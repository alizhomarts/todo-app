package apperror

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordRequired   = errors.New("password required")
	ErrEmailRequired      = errors.New("email required")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrIDRequired         = errors.New("id is required")
	ErrUserRequired       = errors.New("user is required")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrTodoNotFound   = errors.New("todo not found")
	ErrUserIDRequired = errors.New("user id is required")
	ErrTitleRequired  = errors.New("title is required")
	ErrTodoRequired   = errors.New("todo is nil")

	Unauthorized       = errors.New("unauthorized")
	InvalidRequestBody = errors.New("invalid request body")
	InternalServer     = errors.New("internal server error")
	InvalidID          = errors.New("invalid id")
)
