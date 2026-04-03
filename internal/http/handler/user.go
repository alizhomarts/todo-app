package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"todo-app/internal/apperror"
	"todo-app/internal/dto"
	"todo-app/internal/entity"
	"todo-app/internal/logger"
	"todo-app/internal/service/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetByEmail(c echo.Context) error {
	logger.Log.Info("Get user by email called")

	var req *dto.EmailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("getting todo")

	u, err := h.service.GetByEmail(c.Request().Context(), req.Email)
	if err != nil {
		logger.Log.WithError(err).Error("failed to get user by email")
		switch {
		case errors.Is(err, apperror.ErrEmailRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, toUserResponse(u))
}

func (h *UserHandler) GetByID(c echo.Context) error {
	logger.Log.Info("Get user by id called")

	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidID.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"id": id,
	}).Info("getting todo")

	u, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		logger.Log.WithError(err).Error("failed to get user by ID")
		switch {
		case errors.Is(err, apperror.ErrIDRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, toUserResponse(u))
}

func toUserResponse(u *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
	}
}
