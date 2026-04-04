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
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	logger.Log.WithField("email", req.Email).Info("login attempt")

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("login user")

	token, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		logger.Log.WithError(err).Error("login failed")
		switch {
		case errors.Is(err, apperror.ErrUserAlreadyExists):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
