package middleware

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/auth"
)

func (ctrl *Controller) AuthFullCheck(c *fiber.Ctx) error {
	authData := GetAuthData(c)
	if authData == nil || !authData.IsTenantUser() {
		return appErrors.ErrUnauthorized
	}

	c.Locals(auth.ContextKeyAuthCheckTenantAccess, true)

	isConfirmed, err := ctrl.authUC.IsSessionConfirmed(c.Context(), authData.TenantUserClaims().SessionID)
	if err != nil {
		return err
	}

	if !isConfirmed {
		return appErrors.ErrUnauthorized
	}

	return c.Next()
}
