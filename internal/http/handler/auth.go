package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"todo-app/internal/apperror"
	"todo-app/internal/dto"
	"todo-app/internal/entity"
	"todo-app/internal/logger"
	"todo-app/internal/service/auth"
)

type AuthHandler struct {
	service auth.Service
}

func NewAuthHandler(service auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register user
// @Description Create new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	logger.Log.Info("user registration started")

	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	user := &entity.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	id, err := h.service.Register(c.Request().Context(), user, req.Password)
	if err != nil {
		logger.Log.WithError(err).Error("failed to register user")
		switch {
		case errors.Is(err, apperror.ErrUserAlreadyExists):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"user_id": id,
		"email":   req.Email,
	}).Info("user registered successfully")

	return c.JSON(http.StatusCreated, dto.CreateUserResponse{
		ID: id,
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		logger.Log.WithError(err).Warn("invalid login request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	logger.Log.WithField("email", req.Email).Info("login attempt")

	tokens, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		logger.Log.WithError(err).WithField("email", req.Email).Error("login failed")

		switch {
		case errors.Is(err, apperror.ErrEmailRequired),
			errors.Is(err, apperror.ErrPasswordRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		case errors.Is(err, apperror.ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	response := &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	logger.Log.WithField("email", req.Email).Info("login successful")

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return nil
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generates a new access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} map[string]string "Invalid request body or refresh token is required"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		logger.Log.WithError(err).Warn("invalid refresh token request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	token, err := h.service.RefreshAccessToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		logger.Log.WithError(err).Warn("failed to refresh access token")

		switch {
		case errors.Is(err, apperror.ErrRefreshTokenRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		case errors.Is(err, apperror.ErrInvalidRefreshToken):
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	logger.Log.Info("token refreshed")

	response := &dto.RefreshTokenResponse{
		AccessToken: token.AccessToken,
	}

	return c.JSON(http.StatusOK, response)
}
