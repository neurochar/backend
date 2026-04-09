package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/auth_tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/crm"
)

type HTTPGatewayControllers struct {
	AuthTenant *auth_tenant.Controller
	CRM        *crm.Controller
}

type Server struct {
	httpServer *http.Server
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

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		cfg:        cfg,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (gw *Server) RegisterHandlers(handler http.Handler) {
	gw.httpServer.Handler = ChainMiddleware(
		handler,
		RootMiddleware(gw.cfg.TrustedProxies),
		LoggerMiddleware(gw.cfg.UseLogger, gw.logger),
		ErrorHandler(),
		RecoveryMiddleware(),
		Cors(gw.cfg.CorsAllowOrigins),
	)
}

func (gw *Server) Listen() (func() error, error) {
	listener, err := net.Listen("tcp", gw.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen HTTP: %w", err)
	}

	return func() error {
		err = gw.httpServer.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve HTTP: %w", err)
		}

		return nil
	}, nil
}

func (gw *Server) Shutdown(ctx context.Context) error {
	return gw.httpServer.Shutdown(ctx)
}
