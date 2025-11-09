package middleware

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func XCheck() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
			check := c.Get("X-Check")

			if check != "true" {
				return appErrors.ErrBadRequest.WithDetail("check", false, false)
			}

			return c.Next()
		default:
			return c.Next()
		}
	}
}
