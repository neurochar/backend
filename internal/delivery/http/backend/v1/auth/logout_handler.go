package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
)

func (ctrl *Controller) LogoutHandler(c *fiber.Ctx) error {
	const op = "LogoutHandler"

	auth := middleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	err := ctrl.tenantUserFacade.Auth.RevokeSessionByID(
		c.Context(),
		auth.SessionID,
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
		}
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
