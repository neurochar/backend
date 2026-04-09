package auth_tenant

import (
	"context"
	"errors"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *Controller) Login(
	ctx context.Context,
	req *desc.LoginRequest,
) (*desc.LoginResponse, error) {
	const op = "Login"

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
		backoffSession, ok := ctrl.backoff.GetIfExists(ip.String(), backoffConfigAuthGroupID)

		if ok && !backoffSession.IsAllowed() {
			tryAfter := backoffSession.NextAllowedUntilSeconds()

			return nil, appErrors.Chainf(
				appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
				"%s.%s", ctrl.pkg, op,
			)
		}
	}

	authDTO, err := ctrl.tenantFacade.Auth.LoginByEmail(ctx, tenant.ID, req.Email, req.Password, ip)
	if err != nil {
		if errors.Is(err, appErrors.ErrUnauthorized) {
			backoffSession := ctrl.backoff.GetOrCreate(ip.String(), backoffConfigAuthGroupID)
			backoffSession.AddCounter()
			if backoffSession.Counter() > 1 {
				backoffSession.AddBackoff()
			}
		}

		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	accessJWT, err := ctrl.tenantFacade.Auth.IssueAccessJWT(authDTO.AccessClaims)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	refreshJWT, err := ctrl.tenantFacade.Auth.IssueRefreshJWT(authDTO.RefreshClaims)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	accountDTO := mapper.TenantAccountToPb(authDTO.AccountDTO, ctrl.fileUC, true)

	tenantDTO := mapper.TenantToPb(authDTO.AccountDTO.Tenant)

	tokensDTO := mapper.AuthTenantTokensToPb(refreshJWT, int32(ctrl.cfg.Auth.RefreshTokenLifetime.Seconds()), accessJWT)

	return &desc.LoginResponse{
		Account: accountDTO,
		Tenant:  tenantDTO,
		Tokens:  tokensDTO,
	}, nil
}
