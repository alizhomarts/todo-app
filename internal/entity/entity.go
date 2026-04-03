package entity

import (
	"github.com/google/uuid"
	"time"
)

type (
	Todo struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}

	User struct {
		ID           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		PasswordHash string    `json:"-"`
		FirstName    string    `json:"first_name"`
		LastName     string    `json:"last_name"`
		CreatedAt    time.Time `json:"created_at"`
	}
)
