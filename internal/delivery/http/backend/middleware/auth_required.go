package middleware

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) MiddlewareAuthRequired(c *fiber.Ctx) error {
	auth := GetAuthData(c)
	if auth == nil {
		return appErrors.ErrUnauthorized
	}

	return c.Next()
}
