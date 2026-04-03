package dto

import (
	"github.com/google/uuid"
	"time"
)

type (
	TodoCreate struct {
		Title string `json:"title"`
	}

	TodoUpdate struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	TodoResponse struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}

	CreateTodoResponse struct {
		ID uuid.UUID `json:"id"`
	}
)
