package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDKey = "requestId"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.NewString()
		c.Locals(RequestIDKey, requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}
