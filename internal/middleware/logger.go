package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"github.com/yourname/user-api/internal/logger"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start).Milliseconds()
		requestID, _ := c.Locals(RequestIDKey).(string)
		status := c.Response().StatusCode()
		fields := []zap.Field{
			zap.String("method", string(c.Method())),
			zap.String("path", c.OriginalURL()),
			zap.Int("status_code", status),
			zap.Int64("latency_ms", latency),
			zap.String("request_id", requestID),
		}

		if status >= 500 {
			logger.Logger.Error("request completed", fields...)
		} else if status >= 400 {
			logger.Logger.Warn("request completed", fields...)
		} else {
			logger.Logger.Info("request completed", fields...)
		}

		return err
	}
}
