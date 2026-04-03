package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"todo-app/internal/auth"
)

const ContextUserIDKey = "user_id"

func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header",
				})
			}

			claims, err := auth.ParseToken(parts[1], jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token",
				})
			}

			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid user id in token",
				})
			}

			c.Set(ContextUserIDKey, userID)
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uuid.UUID, bool) {
	value := c.Get(ContextUserIDKey)
	userID, ok := value.(uuid.UUID)
	return userID, ok
}
