package providing

import (
	"fmt"
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	privateHTTP "github.com/neurochar/backend/internal/delivery/private_http"
	privateHTTPServer "github.com/neurochar/backend/internal/delivery/private_http/server"
	"go.uber.org/fx"
)

var PrivateHTTPGateway = fx.Options(
	fx.Provide(func(
		logger *slog.Logger,
		cfg config.Config,
	) (*privateHTTPServer.Server, error) {
		privateServer := privateHTTPServer.New(
			logger,
			fmt.Sprintf(":%d", cfg.BackendApp.PrivateHTTP.Port),
		)

		srv := privateHTTP.New(logger, privateServer)

		err := srv.RegisterHandlers()
		if err != nil {
			return nil, err
		}

		return privateServer, nil
	}),
)
