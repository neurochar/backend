package auth_tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *Controller) UpdatePasswordByCode(
	ctx context.Context,
	req *desc.UpdatePasswordByCodeRequest,
) (*desc.UpdatePasswordByCodeResponse, error) {
	const op = "UpdatePasswordByCode"

	ip := tools.GetRealIP(ctx)

	err := ctrl.limiter.Get(limiter.DefaultName).Register(ctx, &limiter.RegisterKey{
		IP: ip,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	codeID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	if req.Password != req.Password2 {
		return nil, appErrors.Chainf(ErrPasswordsMismatch, "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.tenantFacade.Account.UpdatePasswordByRecoveryCode(ctx, codeID, req.Code, req.Password, true)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.UpdatePasswordByCodeResponse{}, nil
}
