package users_tenant

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) UpdateMyPassword(
	ctx context.Context,
	req *desc.UpdateMyPasswordRequest,
) (*desc.UpdateMyPasswordResponse, error) {
	const op = "UpdateMyPassword"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithoutCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	account, err := ctrl.tenantFacade.Account.FindOneByID(
		ctx,
		authData.TenantUserClaims().AccountID,
		nil,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if !account.Account.VerifyPassword(req.Payload.CurrentPassword) {
		return nil, appErrors.Chainf(
			appErrors.ErrBadRequest.WithTextCode("CURRENT_PASSWORD_INCORRECT").WithHints("current password is incorrect"),
			"%s.%s",
			ctrl.pkg,
			op,
		)
	}

	if req.Payload.NewPassword != req.Payload.NewPassword2 {
		return nil, appErrors.Chainf(ErrPasswordsMismatch, "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.tenantFacade.Account.PatchAccountByDTO(
		ctx,
		authData.TenantUserClaims().AccountID,
		tenantUC.PatchAccountDataInput{
			Password: lo.ToPtr(req.Payload.NewPassword),
		},
		true,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.UpdateMyPasswordResponse{}, nil
}
