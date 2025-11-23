package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type UpdatePasswordByCodeHandlerIn struct {
	ID        string `json:"id" validate:"omitempty,uuid"`
	Code      string `json:"code" validate:"required"`
	Password  string `json:"password" validate:""`
	Password2 string `json:"password2" validate:""`
}

func (ctrl *Controller) UpdatePasswordByCodeHandler(c *fiber.Ctx) error {
	const op = "UpdatePasswordByCodeHandler"

	in := &UpdatePasswordByCodeHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	codeID, err := uuid.Parse(in.ID)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	if in.Password != in.Password2 {
		return appErrors.Chainf(ErrPasswordsMismatch, "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.tenantFacade.Account.UpdatePasswordByRecoveryCode(c.Context(), codeID, in.Code, in.Password, true)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
