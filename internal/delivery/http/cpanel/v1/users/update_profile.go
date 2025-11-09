package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

type UpdateProfileHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	InProfile
}

func (ctrl *Controller) UpdateProfileHandler(c *fiber.Ctx) error {
	const op = "UpdateProfileHandler"

	in := &UpdateProfileHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	var photo100x100FileID *uuid.UUID

	if in.Photo100x100FileID != "" {
		parseID, err := uuid.Parse(in.Photo100x100FileID)
		if err != nil {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photo100x100FileID = &parseID
	}

	err = ctrl.userFacade.Profile.UpdateByDTO(c.Context(), uint64(id), userUC.ProfileDataInput{
		Version: in.Version,

		Name:               in.Name,
		Surname:            in.Surname,
		Photo100x100FileID: photo100x100FileID,
	}, in.SkipVersionCheck)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
