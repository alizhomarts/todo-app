package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"todo-app/internal/apperror"
	"todo-app/internal/auth"
	"todo-app/internal/entity"
	"todo-app/internal/repository"
)

const bcryptCost = 12

type (
	Service interface {
		Register(ctx context.Context, user *entity.User, password string) (uuid.UUID, error)
		Login(ctx context.Context, email, password string) (string, error)
	}

	service struct {
		repo      repository.UserRepository
		jwtSecret string
	}
)

func NewService(repo repository.UserRepository, jwtSecret string) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *service) Register(ctx context.Context, user *entity.User, password string) (uuid.UUID, error) {
	if user == nil {
		return uuid.Nil, apperror.ErrUserRequired
	}
	if user.Email == "" {
		return uuid.Nil, apperror.ErrEmailRequired
	}
	if password == "" {
		return uuid.Nil, apperror.ErrPasswordRequired
	}
	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("check user by email: %w", err)
	}
	if existingUser != nil {
		return uuid.Nil, apperror.ErrUserAlreadyExists
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = string(passwordHash)

	return s.repo.Create(ctx, user)
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", apperror.ErrEmailRequired
	}

	if password == "" {
		return "", apperror.ErrPasswordRequired
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if user == nil {
		return "", apperror.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", apperror.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
