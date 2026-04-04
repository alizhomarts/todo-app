package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"todo-app/internal/entity"
)

type (
	UserRepository interface {
		Create(ctx context.Context, user *entity.User) (uuid.UUID, error)
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	}

	userRepository struct {
		db *pgxpool.Pool
	}
)

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	query := `
			  insert into users (email, password_hash, first_name, last_name) 
			  values ($1, $2, $3, $4) 
			  returning id
	`
	var id uuid.UUID

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, fmt.Errorf("create user repository error: %w", err)
	}

	return id, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `
		select id, email, password_hash, first_name, last_name, created_at
		from users
		where email = $1
	`
	if err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	query := `
		select id, email, password_hash, first_name, last_name, created_at
		from users
		where id = $1
	`
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return &user, nil
}
