package auth_tenant

import (
	"context"

	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *controller) Login(
	ctx context.Context,
	req *desc.LoginRequest,
) (*desc.LoginResponse, error) {
	const op = "Login"

	// tenant, err := ctrl.tenantFacade.Tenant.FindOneByTextID(ctx, req.TenantTextId, nil)
	// if err != nil {
	// 	if errors.Is(err, appErrors.ErrNotFound) {
	// 		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithHints("tenant not found"), "%s.%s", ctrl.pkg, op)
	// 	}
	// 	return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	// }

	return &desc.LoginResponse{}, nil
}
