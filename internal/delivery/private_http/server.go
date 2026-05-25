package private_http

import (
	"log/slog"
	"net/http"

	"github.com/neurochar/backend/internal/delivery/private_http/server"
)

type Server struct {
	logger *slog.Logger
	server *server.Server
}

func New(logger *slog.Logger, srv *server.Server) *Server {
	return &Server{
		logger: logger,
		server: srv,
	}
}

func (s *Server) RegisterHandlers() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /metrics/shoud-gpu-nodes", s.MetricsShouldGPUNodes)

	s.server.RegisterHandlers(mux)

	return nil
}

func (s *Server) Server() *server.Server {
	return s.server
}
