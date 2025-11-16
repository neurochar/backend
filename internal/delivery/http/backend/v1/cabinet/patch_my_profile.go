package cabinet

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

type PatchMyProfileHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	ProfileName                string `json:"profileName" validate:"required,min=1,max=150"`
	ProfileSurname             string `json:"profileSurname" validate:"required,min=1,max=150"`
	ProfilePhotoOriginalFileID string `json:"profilePhotoOriginalFileID" validate:"omitempty,uuid"`
	ProfilePhoto100x100FileID  string `json:"profilePhoto100x100FileID" validate:"omitempty,uuid"`
}

func (ctrl *Controller) PatchMyProfileHandler(c *fiber.Ctx) error {
	const op = "PatchMyProfileHandler"

	in := &PatchMyProfileHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	auth := middleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	isRevoked, err := ctrl.tenantUserFacade.Auth.IsSessionRevoked(c.Context(), auth.SessionID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if isRevoked {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	var photoOriginalFileID *uuid.UUID
	var photo100x100FileID *uuid.UUID

	if in.ProfilePhotoOriginalFileID != "" {
		parseID, err := uuid.Parse(in.ProfilePhotoOriginalFileID)
		if err != nil {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photoOriginalFileID = &parseID
	}

	if in.ProfilePhoto100x100FileID != "" {
		parseID, err := uuid.Parse(in.ProfilePhoto100x100FileID)
		if err != nil {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photo100x100FileID = &parseID
	}

	err = ctrl.tenantUserFacade.Account.PatchAccountByDTO(c.Context(), auth.AccountID, tenantUserUC.PatchAccountDataInput{
		Version: in.Version,

		ProfileName:    &in.ProfileName,
		ProfileSurname: &in.ProfileSurname,
		ProfilePhotos: &tenantUserUC.PatchAccountDataInputProfilePhotos{
			PhotoOriginalFileID: photoOriginalFileID,
			Photo100x100FileID:  photo100x100FileID,
		},
	}, in.SkipVersionCheck)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
