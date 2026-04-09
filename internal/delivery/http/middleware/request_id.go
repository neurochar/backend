package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

// RequestID - middleware for request id
func RequestID() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestID, err := uuid.Parse(c.Get("X-Request-ID"))
		if err != nil {
			requestID = uuid.New()
		}

		c.Request().Header.Set("X-Request-ID", requestID.String())

		c.Locals("requestID", requestID)

		ctxData, ctxKey := loghandler.SetData(c.Context(), "request.id", requestID)
		c.Locals(ctxKey, ctxData)

		return c.Next()
	}
}
