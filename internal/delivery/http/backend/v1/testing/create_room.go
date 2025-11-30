package testing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type CreateRoomHandlerIn struct {
	CandidateID uuid.UUID `json:"candidateID" validate:"required"`
	ProfileID   uuid.UUID `json:"profileID" validate:"required"`
}

func (ctrl *Controller) CreateRoomHandler(c *fiber.Ctx) error {
	const op = "CreateRoomHandler"

	in := &CreateRoomHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	authData := backendMiddleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	roomDTO, err := ctrl.testingFacade.Room.CreateByDTO(
		c.Context(),
		authData.TenantID,
		testingUC.CreateRoomDataInput{
			CandidateID: in.CandidateID,
			ProfileID:   in.ProfileID,
			CreatedBy:   &authData.AccountID,
		},
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutRoomDTO(
		c,
		roomDTO,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
