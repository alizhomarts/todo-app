package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
	jwt "todo-app/internal/auth"
	"todo-app/internal/http/handler"
	"todo-app/internal/http/middleware"
)

func Routes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	todoHandler *handler.TodoHandler,
	jwtManager *jwt.JWTManager,
) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	api := e.Group("/api/v1")

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
	api.POST("/refresh", authHandler.Refresh)

	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtManager))
	users.GET("/by-email", userHandler.GetByEmail)
	users.GET("/:id", userHandler.GetByID)

	todos := api.Group("/todos")
	todos.Use(middleware.AuthMiddleware(jwtManager))
	todos.GET("", todoHandler.GetAllByUser)
	todos.GET("/:id", todoHandler.Get)
	todos.POST("", todoHandler.Create)
	todos.PUT("/:id", todoHandler.Update)
	todos.DELETE("/:id", todoHandler.Delete)
}
