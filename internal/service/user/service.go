package user

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
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	}

	service struct {
		repo repository.UserRepository
	}
)

func NewService(repo repository.UserRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, apperror.ErrEmailRequired
	}

	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if user == nil {
		return nil, apperror.ErrUserNotFound
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, apperror.ErrIDRequired
	}

	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user == nil {
		return nil, apperror.ErrUserNotFound
	}

	return user, nil
}
