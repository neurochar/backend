package auth_tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *controller) AccountVerifyEmail(
	ctx context.Context,
	req *desc.AccountVerifyEmailRequest,
) (*desc.AccountVerifyEmailResponse, error) {
	const op = "AccountVerifyEmail"

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
