package auth

import (
	"context"
	"errors"
	"strings"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
)

func (ctrl *Controller) EnrichAuth() func(
	ctx context.Context,
) (context.Context, error) {
	return func(
		ctx context.Context,
	) (context.Context, error) {
		enrich := tools.GetEnrich(ctx)

		if enrich == nil {
			return ctx, nil
		}

		if enrich.S2SToken != "" {
			authData, err := auth.S2SClaimsToAuthData(
				&auth.S2SClaims{
					ServiceID: enrich.S2SToken,
				},
			)
			if err != nil {
				return ctx, err
			}

			ctx = auth.SetAuthData(ctx, authData)
			ctx = loghandler.SetContextData(ctx, "request.s2s", authData.S2SClaims().ServiceID)

			return ctx, nil
		}

		if enrich.AuthorizationToken != "" {
			authToken := strings.TrimPrefix(enrich.AuthorizationToken, "Bearer ")

			claims, err := ctrl.authUC.ParseAccessToken(authToken, true)
			if err != nil {
				if !errors.Is(err, appErrors.ErrUnauthorized) {
					return ctx, err
				}
				return ctx, nil
			}

			authData, err := auth.UserTenantClaimsToAuthData(claims)
			if err != nil {
				if !errors.Is(err, appErrors.ErrUnauthorized) {
					return ctx, err
				}
				return ctx, nil
			}

			ctx = auth.SetAuthData(ctx, authData)
			ctx = loghandler.SetContextData(ctx, "request.account.id", authData.TenantUserClaims().AccountID.String())
			ctx = loghandler.SetContextData(ctx, "request.tenant.id", authData.TenantUserClaims().TenantID.String())

			return ctx, nil
		}

		return ctx, nil
	}
}
