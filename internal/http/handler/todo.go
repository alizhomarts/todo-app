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
	"todo-app/internal/http/middleware"
	"todo-app/internal/logger"
	"todo-app/internal/service/todo"
)

type TodoHandler struct {
	service todo.Service
}

func NewTodoHandler(service todo.Service) *TodoHandler {
	return &TodoHandler{service: service}
}

// Create godoc
// @Summary Create todo
// @Description Create new todo for user
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TodoCreate true "Todo data"
// @Success 201 {object} dto.CreateTodoResponse
// @Failure 400 {object} map[string]string
// @Unauthorized 401 {object} map[string]string
// @Router /todos [post]
func (h *TodoHandler) Create(c echo.Context) error {
	logger.Log.Info("Create todo called")

	var req dto.TodoCreate

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.ErrTitleRequired.Error(),
		})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": apperror.Unauthorized.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"user_id": userID,
		"title":   req.Title,
	}).Info("creating todo")

	res := &entity.Todo{
		UserID: userID,
		Title:  req.Title,
	}

	id, err := h.service.Create(c.Request().Context(), res)
	if err != nil {
		logger.Log.WithError(err).Error("failed to create todo")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"id":      id,
		"user_id": userID,
	}).Info("todo created")

	return c.JSON(http.StatusCreated, dto.CreateTodoResponse{ID: id})
}

// Delete godoc
// @Summary Delete todo
// @Description Delete todo by ID for user
// @Tags Todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c echo.Context) error {
	logger.Log.Info("Delete todo called")

	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidID.Error(),
		})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": apperror.Unauthorized.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("deleting todo")

	err = h.service.Delete(c.Request().Context(), id, userID)
	if err != nil {
		logger.Log.WithError(err).Error("failed to delete todo")
		switch {
		case errors.Is(err, apperror.ErrIDRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrTodoNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// Update godoc
// @Summary Update todo
// @Description Update existing todo for user
// @Tags Todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param request body dto.TodoUpdate true "Todo data"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c echo.Context) error {
	logger.Log.Info("Update todo called")

	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidID.Error(),
		})
	}

	var req dto.TodoUpdate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidRequestBody.Error(),
		})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": apperror.Unauthorized.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"id":        id,
		"user_id":   userID,
		"title":     req.Title,
		"completed": req.Completed,
	}).Info("updating todo")

	td := &entity.Todo{
		ID:        id,
		UserID:    userID,
		Title:     req.Title,
		Completed: req.Completed,
	}

	err = h.service.Update(c.Request().Context(), td)
	if err != nil {
		logger.Log.WithError(err).Error("failed to update todo")
		switch {
		case errors.Is(err, apperror.ErrIDRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrTodoNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrTitleRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// Get godoc
// @Summary Get todo
// @Description Get todo for user by id
// @Tags Todos
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.TodoResponse
// @Unauthorized 401 {object} map[string]string
// @StatusNotFound 400 {object} map[string]string
// @StatusInternalServerError 500 {object} map[string]string
// @Router /todos [get]
func (h *TodoHandler) Get(c echo.Context) error {
	logger.Log.Info("Get todo called")

	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": apperror.InvalidID.Error(),
		})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": apperror.Unauthorized.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"id":      id,
		"user_id": userID,
	}).Info("getting todo")

	res, err := h.service.Get(c.Request().Context(), id, userID)
	if err != nil {
		logger.Log.WithError(err).Error("failed to get todo")
		switch {
		case errors.Is(err, apperror.ErrIDRequired):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, apperror.ErrTodoNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": apperror.InternalServer.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, toTodoResponses([]*entity.Todo{res}))
}

// GetAllByUser godoc
// @Summary Get all todos
// @Description Get all todos for user
// @Tags Todos
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.TodoResponse
// @Unauthorized 401 {object} map[string]string
// @StatusInternalServerError 500 {object} map[string]string
// @Router /todos [get]
func (h *TodoHandler) GetAllByUser(c echo.Context) error {
	logger.Log.Info("Get todos called")
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": apperror.Unauthorized.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("getting todos")

	todos, err := h.service.GetAllByUser(c.Request().Context(), userID)
	if err != nil {
		logger.Log.WithError(err).Error("failed to get todos")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": apperror.InternalServer.Error(),
		})
	}

	return c.JSON(http.StatusOK, toTodoResponses(todos))
}

func toTodoResponses(todos []*entity.Todo) []*dto.TodoResponse {
	res := make([]*dto.TodoResponse, 0, len(todos))

	for _, t := range todos {
		res = append(res, &dto.TodoResponse{
			ID:        t.ID,
			UserID:    t.UserID,
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	return res
}
