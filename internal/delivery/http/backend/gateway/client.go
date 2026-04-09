package gateway

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/app/config"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/app/fxboot/invoking"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
)

var ErrGRPCServerNotCreated = appErrors.ErrInternal.Extend("cant create connect to grpc server")

type DeliveryGrpcClient struct {
	Connection *grpc.ClientConn
}

func NewGrpcClient(cfg config.Config) (*DeliveryGrpcClient, error) {
	intercepts := []grpc.UnaryClientInterceptor{
		HeaderUnaryClientInterceptor(nil),
	}

	cc, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", cfg.BackendApp.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(intercepts...),
	)
	if err != nil {
		return nil, err
	}

	return &DeliveryGrpcClient{
		Connection: cc,
	}, nil
}

func HeaderUnaryClientInterceptor(headers map[string]string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		requestIP, ok := ctx.Value("requestIP").(string)
		if ok && requestIP != "" {
			md.Set("x-forwarded-for", requestIP)
		}

		requestID, ok := ctx.Value("requestID").(uuid.UUID)
		if ok && requestID != uuid.Nil {
			md.Set("x-request-id", requestID.String())
		}

		authToken, ok := ctx.Value(middleware.AuthCtxKeyToken).(string)
		if ok && authToken != "" {
			md.Set("authorization", "Bearer "+authToken)
		}

		for k, v := range headers {
			md.Set(k, v)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func InitGrpcClient(grpcClient *DeliveryGrpcClient) invoking.InvokeInit {
	return invoking.InvokeInit{
		StartBeforeOpen: func(ctx context.Context) error {
			grpcClient.Connection.Connect()
			if !grpcClient.Connection.WaitForStateChange(ctx, connectivity.Idle) {
				// nolint
				_ = grpcClient.Connection.Close()
				return ErrGRPCServerNotCreated
			}

			return nil
		},
		Stop: func(ctx context.Context) error {
			return grpcClient.Connection.Close()
		},
	}
}
