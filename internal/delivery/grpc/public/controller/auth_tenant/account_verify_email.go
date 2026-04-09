package auth_tenant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *Controller) AccountVerifyEmail(
	ctx context.Context,
	req *desc.AccountVerifyEmailRequest,
) (*desc.AccountVerifyEmailResponse, error) {
	const op = "AccountVerifyEmail"

	ip := tools.GetRealIP(ctx)

	err := ctrl.limiter.Get(limiter.DefaultName).Register(ctx, &limiter.RegisterKey{
		IP: ip,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	fmt.Println("!!!", req)

	codeID, err := uuid.Parse(req.CodeId)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.tenantFacade.Account.VerifyAccountEmailByCode(ctx, codeID, req.Code)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.AccountVerifyEmailResponse{}, nil
}
