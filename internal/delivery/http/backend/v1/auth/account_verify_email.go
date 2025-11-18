package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) AccountVerifyEmailHandler(c *fiber.Ctx) error {
	const op = "AccountVerifyEmailHandler"

	codeIDstr := c.Query("code_id")
	codeID, err := uuid.Parse(codeIDstr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	codeValue := c.Query("code")

	err = ctrl.tenantUserFacade.Account.VerifyAccountEmailByCode(c.Context(), codeID, codeValue)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
