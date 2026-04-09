package public

import (
	"log/slog"

	authDelivery "github.com/neurochar/backend/internal/delivery/common/auth"
	"github.com/neurochar/backend/internal/delivery/grpc/server"
)

type PublicServer struct {
	logger *slog.Logger
	server *server.Server
}

func New(
	logger *slog.Logger,
	addr string,
	authDeliveryCtrl *authDelivery.Controller,
	opts ...server.Option,
) *PublicServer {
	opts = append(opts, server.WithExtraInterceptors(
		InterceptorPublic(authDeliveryCtrl),
	))

	server := server.New(logger, addr, opts...)

	return &PublicServer{
		logger: logger,
		server: server,
	}
}

func (s *PublicServer) Server() *server.Server {
	return s.server
}
