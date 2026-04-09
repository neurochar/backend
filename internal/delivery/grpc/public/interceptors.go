package public

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	authDelivery "github.com/neurochar/backend/internal/delivery/common/auth"
	"google.golang.org/grpc"
)

func InterceptorPublic(
	authDeliveryCtrl *authDelivery.Controller,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := authDeliveryCtrl.EnrichAuth()(ctx)
		if err != nil {
			return nil, appErrors.ErrInternal.WithWrap(err)
		} else {
			ctx = authCtx
		}

		return handler(ctx, req)
	}
}
