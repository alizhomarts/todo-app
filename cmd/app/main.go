package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"log"
	"time"
	_ "todo-app/docs"
	jwt "todo-app/internal/auth"
	"todo-app/internal/config"
	"todo-app/internal/db"
	"todo-app/internal/http"
	"todo-app/internal/http/handler"
	"todo-app/internal/http/middleware"
	"todo-app/internal/logger"
	"todo-app/internal/repository"
	authsvc "todo-app/internal/service/auth"
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
	runMigrations(cfg)

	e := echo.New()

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// middleware
	e.Use(middleware.LoggingMiddleware())

	// JWT Manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		time.Minute*15,
		time.Hour*24*7,
	)

	userRepo := repository.NewUserRepository(database)
	todoRepo := repository.NewTodoRepository(database)

	authSvc := authsvc.NewService(userRepo, jwtManager)
	userSvc := user.NewService(userRepo)
	todoSvc := todo.NewService(todoRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	todoHandler := handler.NewTodoHandler(todoSvc)

	// handlers / routes
	http.Routes(e, authHandler, userHandler, todoHandler, jwtManager)

	logger.Log.Info("server started")

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}

func runMigrations(cfg *config.Config) {
	dsn := "postgres://" + cfg.DBUser + ":" + cfg.DBPassword +
		"@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"

	m, err := migrate.New(
		"file://database/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
