package rooms

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/neurochar/backend/pkg/validation"
)

type FinishRoomHandlerIn struct {
	AnswerData map[uint64]any `json:"dataAnswers"`
}

func (ctrl *Controller) FinishRoomHandler(c *fiber.Ctx) error {
	const op = "FinishRoomHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	in := &FinishRoomHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	err = ctrl.testingFacade.Room.Finish(auth.WithoutCheckTenantAccess(c.Context()), id, in.AnswerData)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
