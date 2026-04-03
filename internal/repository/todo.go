package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"todo-app/internal/apperror"
	"todo-app/internal/entity"
)

type (
	TodoRepository interface {
		Create(ctx context.Context, todo *entity.Todo) (uuid.UUID, error)
		Update(ctx context.Context, todo *entity.Todo) error
		Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
		Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Todo, error)
		GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error)
	}

	todoRepository struct {
		db *pgxpool.Pool
	}
)

func NewTodoRepository(db *pgxpool.Pool) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) (uuid.UUID, error) {
	if todo == nil {
		return uuid.Nil, apperror.ErrTodoRequired
	}

	if todo.UserID == uuid.Nil {
		return uuid.Nil, apperror.ErrTodoRequired
	}

	if todo.Title == "" {
		return uuid.Nil, apperror.ErrTitleRequired
	}

	query := `
		insert into todos (user_id, title)
		values ($1, $2)
		returning id
	`

	var id uuid.UUID

	err := r.db.QueryRow(ctx, query,
		todo.UserID,
		todo.Title,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create todo: %w", err)
	}

	return id, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	if todo == nil {
		return fmt.Errorf("todo is nil")
	}

	if todo.ID == uuid.Nil {
		return apperror.ErrIDRequired
	}

	if todo.UserID == uuid.Nil {
		return apperror.ErrUserIDRequired
	}

	if todo.Title == "" {
		return apperror.ErrTitleRequired
	}

	query := `
		update todos
		set title = $1, completed = $2
		where id = $3 and user_id = $4
	`

	cmdTag, err := r.db.Exec(ctx, query,
		todo.Title,
		todo.Completed,
		todo.ID,
		todo.UserID,
	)
	if err != nil {
		return fmt.Errorf("update todo: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrTodoNotFound
	}

	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil {
		return apperror.ErrIDRequired
	}

	if userID == uuid.Nil {
		return apperror.ErrUserIDRequired
	}

	query := `
		delete from todos
		where id = $1 and user_id = $2
	`

	cmdTag, err := r.db.Exec(ctx, query,
		id,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperror.ErrTodoNotFound
	}

	return nil
}

func (r *todoRepository) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Todo, error) {
	if id == uuid.Nil {
		return nil, apperror.ErrIDRequired
	}

	if userID == uuid.Nil {
		return nil, apperror.ErrUserIDRequired
	}

	query := `
		select id, user_id, title, completed, created_at
		from todos
		where id = $1 and user_id = $2
	`

	todo := &entity.Todo{}

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrTodoNotFound
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	return todo, nil
}

func (r *todoRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error) {
	if userID == uuid.Nil {
		return nil, apperror.ErrUserIDRequired
	}

	query := `
		select id, user_id, title, completed, created_at
		from todos
		where user_id = $1
		order by created_at desc
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get todos by user: %w", err)
	}
	defer rows.Close()

	var todos []*entity.Todo

	for rows.Next() {
		todo := &entity.Todo{}

		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todos: %w", err)
	}

	return todos, nil
}
