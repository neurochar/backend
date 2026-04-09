package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/server/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	gRPCServer *grpc.Server
	logger     *slog.Logger
	cfg        Config
}

func New(logger *slog.Logger, addr string, opts ...Option) *Server {
	cfg := Config{
		Addr: addr,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) error {
			var err error
			switch errData := p.(type) {
			case error:
				err = errData
			case string:
				err = appErrors.ErrInternal.Extend(fmt.Sprintf("panic: %s", errData))
			default:
				err = appErrors.ErrInternal.Extend("panic: unknown error happend")
			}

			logger.Error("recovered from panic", slog.Any("panic", p))
			return err
		}),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptor.InterceptorRoot(cfg.TrustedProxies),
		interceptor.InterceptorLogger(logger),
		interceptor.InterceptorErrors(),
		recovery.UnaryServerInterceptor(recoveryOpts...),
		interceptor.InterceptorValidate(),
	}

	unaryInterceptors = append(unaryInterceptors, cfg.ExtraInterceptors...)

	streamInterceptors := []grpc.StreamServerInterceptor{}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	reflection.Register(gRPCServer)

	return &Server{
		cfg:        cfg,
		gRPCServer: gRPCServer,
		logger:     logger,
	}
}

func (gw *Server) Listen() (func() error, error) {
	listener, err := net.Listen("tcp", gw.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen gRPC: %w", err)
	}

	return func() error {
		err = gw.gRPCServer.Serve(listener)
		if err != nil {
			return fmt.Errorf("failed to serve gRPC: %w", err)
		}

		return nil
	}, nil
}

func (gw *Server) GRPCServer() *grpc.Server {
	return gw.gRPCServer
}

func (gw *Server) Shutdown() {
	gw.gRPCServer.GracefulStop()
}
