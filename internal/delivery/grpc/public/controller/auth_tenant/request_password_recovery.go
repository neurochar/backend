package auth_tenant

import (
	"context"
	"errors"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *Controller) RequestPasswordRecovery(
	ctx context.Context,
	req *desc.RequestPasswordRecoveryRequest,
) (*desc.RequestPasswordRecoveryResponse, error) {
	const op = "RequestPasswordRecovery"

	ip := tools.GetRealIP(ctx)

	err := ctrl.limiter.Get(limiter.DefaultName).Register(ctx, &limiter.RegisterKey{
		IP: ip,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	tenant, err := ctrl.tenantFacade.Tenant.FindOneByTextID(ctx, req.TenantTextId, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithHints("tenant not found"), "%s.%s", ctrl.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if ip != nil {
		backoffSession := ctrl.backoff.GetOrCreate(ip.String(), backoffConfigPasswordRecoveryGroupID)

		if !backoffSession.IsAllowed() {
			tryAfter := backoffSession.NextAllowedUntilSeconds()

			return nil, appErrors.Chainf(
				appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
				"%s.%s", ctrl.pkg, op,
			)
		}

		_ = backoffSession.AddBackoff()
	}

	code, err := ctrl.tenantFacade.Account.RequestPasswordRecoveryByEmail(ctx, tenant.ID, req.Email, ip)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.RequestPasswordRecoveryResponse{
		CodeId: code.ID.String(),
	}, nil
}
