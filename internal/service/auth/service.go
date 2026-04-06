package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"todo-app/internal/apperror"
	"todo-app/internal/auth"
	"todo-app/internal/entity"
	"todo-app/internal/repository"
	"todo-app/internal/service/auth/model"
)

const bcryptCost = 12

type (
	Service interface {
		Register(ctx context.Context, user *entity.User, password string) (uuid.UUID, error)
		Login(ctx context.Context, email, password string) (*model.TokenPair, error)
		RefreshAccessToken(ctx context.Context, refreshToken string) (*model.RefreshAccessTokenResult, error)
	}

	service struct {
		repo       repository.UserRepository
		jwtManager *auth.JWTManager
	}
)

func NewService(repo repository.UserRepository, jwtManager *auth.JWTManager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *service) Register(ctx context.Context, user *entity.User, password string) (uuid.UUID, error) {
	if user == nil {
		return uuid.Nil, apperror.ErrUserRequired
	}

	email := strings.TrimSpace(user.Email)
	if email == "" {
		return uuid.Nil, apperror.ErrEmailRequired
	}
	if password == "" {
		return uuid.Nil, apperror.ErrPasswordRequired
	}
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, apperror.ErrUserNotFound) {
		return uuid.Nil, fmt.Errorf("check user by email: %w", err)
	}
	if existingUser != nil {
		return uuid.Nil, apperror.ErrUserAlreadyExists
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	newUser := &entity.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	return s.repo.Create(ctx, newUser)
}

func (s *service) Login(ctx context.Context, email, password string) (*model.TokenPair, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, apperror.ErrEmailRequired
	}

	if password == "" {
		return nil, apperror.ErrPasswordRequired
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			return nil, apperror.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if user == nil {
		return nil, apperror.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	tokens, err := s.generateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token pair: %w", err)
	}

	return tokens, nil
}

func (s *service) RefreshAccessToken(ctx context.Context, refreshToken string) (*model.RefreshAccessTokenResult, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, apperror.ErrRefreshTokenRequired
	}

	claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperror.ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("generate new access token: %w", err)
	}

	return &model.RefreshAccessTokenResult{
		AccessToken: newAccessToken,
	}, nil
}

func (s *service) generateTokenPair(userID uuid.UUID) (*model.TokenPair, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
