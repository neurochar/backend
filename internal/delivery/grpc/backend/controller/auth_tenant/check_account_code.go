package auth_tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *controller) CheckAccountCode(
	ctx context.Context,
	req *desc.CheckAccountCodeRequest,
) (*desc.CheckAccountCodeResponse, error) {
	const op = "CheckAccountCode"

	codeID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	_, err = ctrl.tenantFacade.Account.CheckCode(ctx, codeID, req.Code)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.CheckAccountCodeResponse{}, nil
}
