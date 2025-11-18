package users

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
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

	user, err := ctrl.tenantUserFacade.Account.FindOneByID(c.Context(), auth.AccountID, nil, nil)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithParent(err), "%s.%s", ctrl.pkg, op)
	}

	ip := middleware.GetRealIP(c)
	requestIP := net.ParseIP(ip)

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

	accountDTO, err := ctrl.tenantUserFacade.Common.CreateUser(c.Context(), auth.TenantID, tenantUserUC.CreateAccountDataInput{
		Email:          in.Email,
		Password:       in.Password,
		RoleID:         in.RoleID,
		ProfileName:    in.ProfileName,
		ProfileSurname: in.ProfileSurname,
		ProfilePhotos: &tenantUserUC.AccountDataInputProfilePhotos{
			PhotoOriginalFileID: photoOriginalFileID,
			Photo100x100FileID:  photo100x100FileID,
		},
	}, user, requestIP)
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
