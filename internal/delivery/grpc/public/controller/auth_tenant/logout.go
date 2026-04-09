package auth_tenant

import (
	"context"
	"errors"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *Controller) Logout(
	ctx context.Context,
	req *desc.LogoutRequest,
) (*desc.LogoutResponse, error) {
	const op = "Logout"

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	err := ctrl.tenantFacade.Session.RevokeSessionByID(
		ctx,
		authData.TenantUserClaims().SessionID,
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.LogoutResponse{}, nil
}
