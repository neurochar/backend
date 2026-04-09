package interceptor

import (
	"context"

	"buf.build/go/protovalidate"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func InterceptorValidate() grpc.UnaryServerInterceptor {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		reqProto, ok := req.(proto.Message)
		if !ok {
			return handler(ctx, req)
		}

		if err := v.Validate(reqProto); err != nil {
			return nil, appErrors.ErrBadRequest.WithHints(err.Error())
		}

		return handler(ctx, req)
	}
}
