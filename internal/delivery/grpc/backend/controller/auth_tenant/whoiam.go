package auth_tenant

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/backend/mapper"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *controller) WhoIAm(
	ctx context.Context,
	req *desc.WhoIAmRequest,
) (*desc.WhoIAmResponse, error) {
	const op = "WhoIAm"

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	accountDTO, err := ctrl.tenantFacade.Account.FindOneByID(
		ctx,
		authData.TenantUserClaims().AccountID,
		nil,
		&tenantUC.AccountDTOOptions{
			FetchTenant:     true,
			FetchPhotoFiles: true,
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.WhoIAmResponse{
		Account: mapper.AuthTenantAccountToPb(accountDTO, ctrl.fileUC, false),
		Tenant:  mapper.TenantToPb(accountDTO.Tenant),
	}, nil
}
