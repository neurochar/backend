package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"
	"github.com/samber/lo"

	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

type UpdateMyProfileHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	InProfile
}

func (ctrl *Controller) UpdateMyProfileHandler(c *fiber.Ctx) error {
	const op = "UpdateMyProfileHandler"

	in := &UpdateMyProfileHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	authData := middleware.GetAuthData(c)

	if authData == nil {
		return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", ctrl.pkg, op)
	}

	profiles, err := ctrl.userFacade.Profile.FindList(c.Context(), &userUC.ProfileListOptions{
		AccountID: lo.ToPtr(authData.Account.ID),
	}, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if len(profiles) == 0 {
		return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", ctrl.pkg, op)
	}

	var photo100x100FileID *uuid.UUID

	if in.Photo100x100FileID != "" {
		parseID, err := uuid.Parse(in.Photo100x100FileID)
		if err != nil {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photo100x100FileID = &parseID
	}

	err = ctrl.userFacade.Profile.UpdateByDTO(c.Context(), profiles[0].ID, userUC.ProfileDataInput{
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
