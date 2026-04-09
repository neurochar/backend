package tests

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) InternalErrorHandler(c *fiber.Ctx) error {
	const op = "InternalErrorHandler"

	return appErrors.Chainf(appErrors.ErrInternal.WithHints("this is test internal error"), "%s.%s", ctrl.pkg, op)
}
