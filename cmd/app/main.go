package main

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	_ "todo-app/docs"
	"todo-app/internal/config"
	"todo-app/internal/db"
	"todo-app/internal/http"
	"todo-app/internal/http/handler"
	"todo-app/internal/http/middleware"
	"todo-app/internal/logger"
	"todo-app/internal/repository"
	"todo-app/internal/service/auth"
	"todo-app/internal/service/todo"
	"todo-app/internal/service/user"
)

// @title Todo API
// @version 1.0
// @description Simple Todo API with JWT authentication
// @host localhost:8888
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	database := db.NewPostgres(cfg)
	defer database.Close()

	e := echo.New()

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// middleware
	e.Use(middleware.LoggingMiddleware())

	userRepo := repository.NewUserRepository(database)
	todoRepo := repository.NewTodoRepository(database)

	authSvc := auth.NewService(userRepo, cfg.JWTSecret)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	todoHandler := handler.NewTodoHandler(todoSvc)

	// handlers / routes
	http.Routes(e, authHandler, userHandler, todoHandler, cfg.JWTSecret)

	logger.Log.Info("server started")

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
