package backend

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/grpc/interceptor"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerOptions struct {
	Port               int
	LogResponseSent    bool
	LogPayloadReceived bool
	PrivateIPs         []string
}

func New(logger *slog.Logger, ops ServerOptions, authUC tenantUC.AuthUsecase) *grpc.Server {
	loggingEvents := []logging.LoggableEvent{
		logging.FinishCall,
	}
	if ops.LogPayloadReceived {
		loggingEvents = append(loggingEvents, logging.PayloadReceived)
	}
	if ops.LogResponseSent {
		loggingEvents = append(loggingEvents, logging.PayloadSent)
	}

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(loggingEvents...),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			logger.Error("recovered from panic", slog.Any("panic", p))
			return appErrors.ErrInternal
		}),
	}

	middlewareCtrl := middleware.New(authUC)

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptor.InterceptorErrors(),
		recovery.UnaryServerInterceptor(recoveryOpts...),
		interceptor.InterceptorRequestIP(ops.PrivateIPs),
		interceptor.InterceptorRequestID(),
		logging.UnaryServerInterceptor(interceptor.InterceptorLogger(logger), loggingOpts...),
		interceptor.InterceptorValidate(),
		middlewareCtrl.Auth(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	reflection.Register(gRPCServer)

	return gRPCServer
}

func Start(grpcServer *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}
