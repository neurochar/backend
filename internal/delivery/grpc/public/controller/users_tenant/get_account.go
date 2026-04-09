package users_tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
)

func (ctrl *Controller) GetAccount(ctx context.Context, req *desc.GetAccountRequest) (*desc.GetAccountResponse, error) {
	const op = "GetAccount"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	accountDTO, err := ctrl.tenantFacade.Account.FindOneByID(ctx, id, nil, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrForbidden) {
			return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", ctrl.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.GetAccountResponse{
		Item: mapper.TenantAccountToPb(accountDTO, ctrl.fileUC, true),
	}, nil
}
