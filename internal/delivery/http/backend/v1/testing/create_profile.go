package testing

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type CreateProfileHandlerIn struct {
	Name                 string                                    `json:"name" validate:"required,min=1,max=150"`
	PersonalityTraitsMap testingEntity.ProfilePersonalityTraitsMap `json:"personalityTraitsMap" validate:"required"`
}

func (ctrl *Controller) CreateProfileHandler(c *fiber.Ctx) error {
	const op = "CreateProfileHandler"

	in := &CreateProfileHandlerIn{}

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

	profileDTO, err := ctrl.testingFacade.Profile.CreateByDTO(
		c.Context(),
		authData.TenantID,
		testingUC.CreateProfileDataInput{
			Name:                 in.Name,
			PersonalityTraitsMap: in.PersonalityTraitsMap,
			CreatedBy:            &authData.AccountID,
		},
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutProfileDTO(
		c,
		profileDTO,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
