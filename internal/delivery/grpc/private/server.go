package private

import (
	"log/slog"

	"github.com/neurochar/backend/internal/delivery/grpc/server"
)

type PrivateServer struct {
	logger *slog.Logger
	server *server.Server
}

func New(
	logger *slog.Logger,
	addr string,
	opts ...server.Option,
) *PrivateServer {
	opts = append(opts, server.WithExtraInterceptors())

	server := server.New(logger, addr, opts...)

	return &PrivateServer{
		logger: logger,
		server: server,
	}
}

func (s *PrivateServer) Server() *server.Server {
	return s.server
}
