package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

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

func (s *Server) RegisterHandlers(handler http.Handler) {
	s.httpServer.Handler = RecoveryMiddleware(handler)
}

func (s *Server) Listen() (func() error, error) {
	listener, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen private HTTP: %w", err)
	}

	return func() error {
		err = s.httpServer.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve private HTTP: %w", err)
		}

		return nil
	}, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
