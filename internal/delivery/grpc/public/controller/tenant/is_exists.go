package tenant

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/tenant/v1"
)

func (ctrl *Controller) IsExists(ctx context.Context, req *desc.IsExistsRequest) (*desc.IsExistsResponse, error) {
	const op = "IsExists"

	_, err := ctrl.tenantFacade.Tenant.FindOneByTextID(ctx, req.TextId, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.IsExistsResponse{}, nil
}
