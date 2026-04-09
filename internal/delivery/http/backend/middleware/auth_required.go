package middleware

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/auth"
)

func (ctrl *Controller) AuthRequired(c *fiber.Ctx) error {
	authData := GetAuthData(c)
	if authData == nil || !authData.IsTenantUser() {
		return appErrors.ErrUnauthorized
	}

	c.Locals(auth.ContextKeyAuthCheckTenantAccess, true)

	return c.Next()
}
