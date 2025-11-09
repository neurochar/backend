package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
)

// MeHandler - my profile handler
func (ctrl *Controller) MeHandler(c *fiber.Ctx) error {
	const op = "MyProfileHandler"

	authData := middleware.GetAuthData(c)

	if authData == nil {
		return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", ctrl.pkg, op)
	}

	user, err := ctrl.userFacade.Common.FindOneByAccountID(c.Context(), authData.Account.ID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutUserDTO(c, ctrl.fileUC, user.Account, user.ProfileDTO, authData.Role, true)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
