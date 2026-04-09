package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type CheckAccountCodeHandlerIn struct {
	ID   string `json:"id" validate:"omitempty,uuid"`
	Code string `json:"code" validate:"required"`
}

func (ctrl *Controller) CheckAccountCodeHandler(c *fiber.Ctx) error {
	const op = "CheckAccountCodeHandler"

	in := &CheckAccountCodeHandlerIn{}

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

	_, err = ctrl.tenantFacade.Account.CheckCode(c.Context(), codeID, in.Code)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
