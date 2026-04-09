package interceptor

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"google.golang.org/grpc"
)

func InterceptorErrors() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		data, err := handler(ctx, req)

		return data, appErrors.ToGrpcStatus(err)
	}
}
