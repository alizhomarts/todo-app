package todo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"todo-app/internal/apperror"
	"todo-app/internal/entity"
	"todo-app/internal/repository"
)

type (
	Service interface {
		Create(ctx context.Context, todo *entity.Todo) (uuid.UUID, error)
		Update(ctx context.Context, todo *entity.Todo) error
		Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
		Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Todo, error)
		GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error)
	}

	service struct {
		repo repository.TodoRepository
	}
)

func NewService(repo repository.TodoRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, todo *entity.Todo) (uuid.UUID, error) {
	if todo == nil {
		return uuid.Nil, apperror.ErrTodoRequired
	}

	if todo.UserID == uuid.Nil {
		return uuid.Nil, apperror.ErrUserIDRequired
	}

	if todo.Title == "" {
		return uuid.Nil, apperror.ErrTitleRequired
	}

	id, err := s.repo.Create(ctx, todo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create todo: %w", err)
	}

	return id, nil
}

func (s *service) Update(ctx context.Context, todo *entity.Todo) error {
	if todo == nil {
		return apperror.ErrTodoRequired
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

	err := s.repo.Update(ctx, todo)
	if err != nil {
		return fmt.Errorf("update todo: %w", err)
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if id == uuid.Nil {
		return apperror.ErrIDRequired
	}

	if userID == uuid.Nil {
		return apperror.ErrUserIDRequired
	}

	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	return nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Todo, error) {
	if id == uuid.Nil {
		return nil, apperror.ErrIDRequired
	}

	if userID == uuid.Nil {
		return nil, apperror.ErrUserIDRequired
	}

	todo, err := s.repo.Get(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get todo by id: %w", err)
	}

	return todo, nil
}

func (s *service) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error) {
	if userID == uuid.Nil {
		return nil, apperror.ErrUserIDRequired
	}

	todos, err := s.repo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get todos by user: %w", err)
	}

	return todos, nil
}
