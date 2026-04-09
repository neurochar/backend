package public

import (
	"log/slog"

	"github.com/neurochar/backend/internal/delivery/httpgw/public/controller"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
)

type Gateway struct {
	logger *slog.Logger
	server *server.Server
	ctrl   *controller.Controller
}

func New(logger *slog.Logger, server *server.Server, controls *controller.Controls) *Gateway {
	ctrl := controller.New(controls)

	return &Gateway{
		logger: logger,
		server: server,
		ctrl:   ctrl,
	}
}
