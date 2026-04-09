package providing

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	backendGRPC "github.com/neurochar/backend/internal/delivery/grpc/backend"
	"github.com/neurochar/backend/internal/delivery/grpc/backend/controller"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var BackendGRPCServer = fx.Options(
	fx.Provide(func(logger *slog.Logger, cfg config.Config, authUC tenantUC.AuthUsecase) *grpc.Server {
		return backendGRPC.New(
			logger,
			backendGRPC.ServerOptions{
				Port:               cfg.BackendApp.GRPC.Port,
				LogResponseSent:    cfg.BackendApp.GRPC.LogResponseSent,
				LogPayloadReceived: cfg.BackendApp.GRPC.LogPayloadReceived,
				PrivateIPs:         []string{cfg.Global.ServerIP},
			},
			authUC,
		)
	}),
	controller.FxModule,
)
