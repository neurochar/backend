package providing

import (
	"fmt"
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	authDelivery "github.com/neurochar/backend/internal/delivery/common/auth"
	grpcPublic "github.com/neurochar/backend/internal/delivery/grpc/public"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller"
	"go.uber.org/fx"
)

var PublicGRPCServer = fx.Options(
	fx.Provide(func(
		logger *slog.Logger,
		cfg config.Config,
		authDeliveryCtrl *authDelivery.Controller,
	) *grpcPublic.PublicServer {
		return grpcPublic.New(
			logger,
			fmt.Sprintf(":%d", cfg.BackendApp.GRPC.Port),
			authDeliveryCtrl,
		)
	}),
	controller.FxModule,
)
