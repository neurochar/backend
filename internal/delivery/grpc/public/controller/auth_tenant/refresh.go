package auth_tenant

import (
	"context"
	"errors"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func (ctrl *Controller) Refresh(
	ctx context.Context,
	req *desc.RefreshRequest,
) (*desc.RefreshResponse, error) {
	const op = "Refresh"

	ip := tools.GetRealIP(ctx)

	claims, err := ctrl.tenantFacade.Auth.ParseRefreshToken(req.RefreshToken, true)
	if err != nil {
		if errors.Is(err, tenantUC.ErrInvalidToken) {
			return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	authDTO, err := ctrl.tenantFacade.Auth.GenerateNewClaims(ctx, claims, ip)
	if err != nil {
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

	tokensDTO := mapper.AuthTenantTokensToPb(refreshJWT, int32(ctrl.cfg.Auth.RefreshTokenLifetime.Seconds()), accessJWT)

	return &desc.RefreshResponse{
		Tokens: tokensDTO,
	}, nil
}
