package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"
)

type UpdateMyEmailHandlerIn struct {
	Email string `json:"email" validate:"required,email"`
}

func (ctrl *Controller) UpdateMyEmailHandler(c *fiber.Ctx) error {
	const op = "UpdateMyEmailHandler"

	in := &UpdateMyEmailHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	authData := middleware.GetAuthData(c)

	if authData == nil {
		return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", ctrl.pkg, op)
	}

	err := authData.Account.SetEmail(in.Email)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.userFacade.Account.UpdateAccount(c.Context(), authData.Account)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
