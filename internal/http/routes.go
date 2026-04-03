package http

import (
	"github.com/labstack/echo/v4"
	"todo-app/internal/http/handler"
	"todo-app/internal/http/middleware"
)

func Routes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	todoHandler *handler.TodoHandler,
	jwtSecret string,
) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	users := e.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtSecret))
	users.GET("/email", userHandler.GetByEmail)
	users.GET("/:id", userHandler.GetByID)

	todos := e.Group("/todos")
	todos.Use(middleware.AuthMiddleware(jwtSecret))
	todos.GET("", todoHandler.GetAllByUser)
	todos.GET("/:id", todoHandler.Get)
	todos.POST("", todoHandler.Create)
	todos.PUT("/:id", todoHandler.Update)
	todos.DELETE("/:id", todoHandler.Delete)
}
