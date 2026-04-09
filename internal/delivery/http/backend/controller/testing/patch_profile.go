package testing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/null"
	"github.com/neurochar/backend/pkg/validation"

	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
)

type PatchProfileHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	Name *string `json:"name" validate:"required,min=1,max=150"`
	// nolint
	PersonalityTraitsMap null.Nullable[testingEntity.ProfilePersonalityTraitsMap] `json:"personalityTraitsMap" validate:"required"`
}

func (ctrl *Controller) PatchProfileHandler(c *fiber.Ctx) error {
	const op = "PatchProfileHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	in := &PatchProfileHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	usecaseInput := testingUC.PatchProfileDataInput{
		Version: in.Version,

		Name: in.Name,
	}

	if in.PersonalityTraitsMap.IsSet {
		var value testingEntity.ProfilePersonalityTraitsMap
		if in.PersonalityTraitsMap.IsValid {
			value = in.PersonalityTraitsMap.Value
		}
		usecaseInput.PersonalityTraitsMap = &value
	}

	err = ctrl.testingFacade.Profile.PatchByDTO(
		c.Context(),
		id,
		usecaseInput,
		in.SkipVersionCheck,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
