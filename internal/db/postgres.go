package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"todo-app/internal/config"
)

func NewPostgres(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	log.Println("connected to postgres")
	return pool
}
