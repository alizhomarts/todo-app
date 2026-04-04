package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"log"
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
	runMigrations()

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

func runMigrations() {
	m, err := migrate.New(
		"file://database/migrations",
		"postgres://postgres:postgres@db:5432/todo_db?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
