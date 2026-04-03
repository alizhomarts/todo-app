package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
	"todo-app/internal/logger"
)

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			logger.Log.WithFields(logrus.Fields{
				"method":     c.Request().Method,
				"path":       c.Request().URL.Path,
				"status":     c.Response().Status,
				"latency_ms": time.Since(start).Milliseconds(),
				"ip":         c.RealIP(),
			}).Info("http request")

			return err
		}
	}
}
