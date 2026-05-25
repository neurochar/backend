package providing

import (
	"fmt"
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	grpcPrivate "github.com/neurochar/backend/internal/delivery/grpc/private"
	"github.com/neurochar/backend/internal/delivery/grpc/private/controller"
	"go.uber.org/fx"
)

var PrivateGRPCServer = fx.Options(
	fx.Provide(func(
		logger *slog.Logger,
		cfg config.Config,
	) *grpcPrivate.PrivateServer {
		return grpcPrivate.New(
			logger,
			fmt.Sprintf(":%d", cfg.BackendApp.PrivateGRPC.Port),
		)
	}),
	controller.FxModule,
)
