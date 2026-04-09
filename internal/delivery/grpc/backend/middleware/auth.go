package middleware

import (
	"context"
	"strings"

	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (ctrl *Controller) Auth() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if s2sValues := md.Get("s2s"); len(s2sValues) > 0 {
				s2sToken := s2sValues[0]
				if s2sToken != "" {
					authData, err := auth.S2SClaimsToAuthData(
						&auth.S2SClaims{
							ServiceID: s2sToken,
						},
					)
					if err != nil {
						return nil, err
					}

					ctx = auth.SetAuthData(ctx, authData)
					ctx = loghandler.SetContextData(ctx, "request.s2s", authData.S2SClaims().ServiceID)

					return handler(ctx, req)
				}
			}

			if authValues := md.Get("authorization"); len(authValues) > 0 {
				authToken := authValues[0]
				if authToken != "" {
					authToken = strings.TrimPrefix(authToken, "Bearer ")

					claims, err := ctrl.authUC.ParseAccessToken(authToken, true)
					if err != nil {
						return nil, err
					}

					authData, err := auth.UserTenantClaimsToAuthData(claims)
					if err != nil {
						return nil, err
					}

					ctx = auth.SetAuthData(ctx, authData)
					ctx = loghandler.SetContextData(ctx, "request.account.id", authData.TenantUserClaims().AccountID.String())
					ctx = loghandler.SetContextData(ctx, "request.tenant.id", authData.TenantUserClaims().TenantID.String())

					return handler(ctx, req)
				}
			}

			return handler(ctx, req)
		}

		return handler(ctx, req)
	}
}
