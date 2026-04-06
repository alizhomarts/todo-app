package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"todo-app/internal/apperror"
	"todo-app/internal/auth"
)

type contextKey string

const ContextUserIDKey contextKey = "user_id"

func AuthMiddleware(jwtManager *auth.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": apperror.ErrMissingAuthHeader.Error(),
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": apperror.ErrInvalidAuthHeader.Error(),
				})
			}

			token := parts[1]

			claims, err := jwtManager.ParseAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": apperror.ErrInvalidToken.Error(),
				})
			}

			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": apperror.ErrInvalidToken.Error(),
				})
			}

			c.Set(string(ContextUserIDKey), userID)
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uuid.UUID, bool) {
	value := c.Get(string(ContextUserIDKey))
	userID, ok := value.(uuid.UUID)
	return userID, ok
}
