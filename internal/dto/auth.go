package dto

import "github.com/google/uuid"

type (
	RegisterRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		Token string `json:"token"`
	}

	CreateUserResponse struct {
		ID uuid.UUID `json:"id"`
	}
)
