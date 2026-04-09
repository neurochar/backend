package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type CreateAccountHandlerIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=0"`
	RoleID   uint64 `json:"roleID" validate:"required"`

	ProfileName                string `json:"profileName" validate:"required,min=1,max=150"`
	ProfileSurname             string `json:"profileSurname" validate:"required,min=1,max=150"`
	ProfilePhotoOriginalFileID string `json:"profilePhotoOriginalFileID" validate:"omitempty,uuid"`
	ProfilePhoto100x100FileID  string `json:"profilePhoto100x100FileID" validate:"omitempty,uuid"`
}

func (ctrl *Controller) CreateAccountHandler(c *fiber.Ctx) error {
	const op = "CreateAccountHandler"

	in := &CreateAccountHandlerIn{}

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
	if auth == nil || !auth.IsTenantUser() {
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

	accountDTO, _, err := ctrl.tenantFacade.Account.CreateAccountByDTO(
		c.Context(),
		auth.TenantUserClaims().TenantID,
		tenantUC.CreateAccountDataInput{
			Email:          in.Email,
			Password:       in.Password,
			RoleID:         in.RoleID,
			ProfileName:    in.ProfileName,
			ProfileSurname: in.ProfileSurname,
			ProfilePhotos: &tenantUC.AccountDataInputProfilePhotos{
				PhotoOriginalFileID: photoOriginalFileID,
				Photo100x100FileID:  photo100x100FileID,
			},
		},
		true,
		nil,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutAccountDTO(
		c,
		true,
		ctrl.fileUC,
		accountDTO,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
