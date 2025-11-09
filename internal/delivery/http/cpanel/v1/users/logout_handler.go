package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
)

func (ctrl *Controller) LogoutHandler(c *fiber.Ctx) error {
	const op = "LogoutHandler"

	authData := middleware.GetAuthData(c)

	if authData == nil {
		return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", ctrl.pkg, op)
	}

	err := ctrl.userFacade.AdminAuth.DeleteActiveSessionByID(c.Context(), authData.Session.ID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	c.ClearCookie(authCookie)

	return nil
}
