package users_tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) CreateAccount(
	ctx context.Context,
	req *desc.CreateAccountRequest,
) (*desc.CreateAccountResponse, error) {
	const op = "CreateAccount"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	var photoOriginalFileID *uuid.UUID
	var photo100x100FileID *uuid.UUID

	if req.Payload.ProfilePhotos != nil {
		parseID, err := uuid.Parse(req.Payload.ProfilePhotos.OriginalFileId)
		if err != nil {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photoOriginalFileID = lo.ToPtr(parseID)

		parseID, err = uuid.Parse(req.Payload.ProfilePhotos.S100X100FileId)
		if err != nil {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photo100x100FileID = lo.ToPtr(parseID)
	}

	accountDTO, _, err := ctrl.tenantFacade.Account.CreateAccountByDTO(
		ctx,
		authData.TenantUserClaims().TenantID,
		tenantUC.CreateAccountDataInput{
			Email:          req.Payload.Email,
			Password:       req.Payload.Password,
			RoleID:         req.Payload.RoleId,
			ProfileName:    req.Payload.ProfileName,
			ProfileSurname: req.Payload.ProfileSurname,
			ProfilePhotos: &tenantUC.AccountDataInputProfilePhotos{
				PhotoOriginalFileID: photoOriginalFileID,
				Photo100x100FileID:  photo100x100FileID,
			},
		},
		true,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.CreateAccountResponse{
		Item: mapper.TenantAccountToPb(accountDTO, ctrl.fileUC, true),
	}, nil
}
