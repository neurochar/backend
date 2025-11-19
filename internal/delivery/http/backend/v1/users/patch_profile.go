package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

type PatchProfileHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	IsBlocked      *bool                               `json:"isBlocked" validate:"omitempty"`
	Password       *string                             `json:"password" validate:"omitempty"`
	RoleID         *uint64                             `json:"roleID" validate:"omitempty"`
	ProfileName    *string                             `json:"profileName" validate:"omitempty,max=150"`
	ProfileSurname *string                             `json:"profileSurname" validate:"omitempty,max=150"`
	ProfilePhotos  *PatchProfileHandlerInProfilePhotos `json:"profilePhotos" validate:"omitempty"`
}

type PatchProfileHandlerInProfilePhotos struct {
	PhotoOriginalFileID string `json:"photoOriginalFileID" validate:"omitempty,uuid"`
	Photo100x100FileID  string `json:"photo100x100FileID" validate:"omitempty,uuid"`
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

	auth := backendMiddleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	isConfirmed, err := ctrl.tenantUserFacade.Auth.IsSessionConfirmed(c.Context(), auth.SessionID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if !isConfirmed {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	var profilePhotos *tenantUserUC.AccountDataInputProfilePhotos

	if in.ProfilePhotos != nil {
		profilePhotos = &tenantUserUC.AccountDataInputProfilePhotos{}

		if in.ProfilePhotos.PhotoOriginalFileID != "" {
			parseID, err := uuid.Parse(in.ProfilePhotos.PhotoOriginalFileID)
			if err != nil {
				return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
			}

			profilePhotos.PhotoOriginalFileID = &parseID
		}

		if in.ProfilePhotos.Photo100x100FileID != "" {
			parseID, err := uuid.Parse(in.ProfilePhotos.Photo100x100FileID)
			if err != nil {
				return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
			}

			profilePhotos.Photo100x100FileID = &parseID
		}
	}

	usecaseInput := tenantUserUC.PatchAccountDataInput{
		Version: in.Version,

		IsBlocked:      in.IsBlocked,
		Password:       in.Password,
		RoleID:         in.RoleID,
		ProfileName:    in.ProfileName,
		ProfileSurname: in.ProfileSurname,
		ProfilePhotos:  profilePhotos,
	}

	err = ctrl.tenantUserFacade.Common.PatchAccountByDTO(
		c.Context(),
		id,
		usecaseInput,
		auth.AccountID,
		in.SkipVersionCheck,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
