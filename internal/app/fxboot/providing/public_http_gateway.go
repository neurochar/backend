package providing

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	authDelivery "github.com/neurochar/backend/internal/delivery/common/auth"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/auth_tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/crm"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/registration"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/rooms"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/testing"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/users_tenant"
	"github.com/neurochar/backend/internal/delivery/httpgw/public"
	"github.com/neurochar/backend/internal/delivery/httpgw/public/controller"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
	"go.uber.org/fx"
)

var PublicHTTPGateway = fx.Options(
	fx.Provide(func(
		logger *slog.Logger,
		cfg config.Config,
		authDeliveryCtrl *authDelivery.Controller,
		authTenantCtrl *auth_tenant.Controller,
		crmCtrl *crm.Controller,
		tenantCtrl *tenant.Controller,
		registrationCtrl *registration.Controller,
		roomsCtrl *rooms.Controller,
		usersTenantCtrl *users_tenant.Controller,
		testingCtrl *testing.Controller,
	) (*server.Server, *public.Gateway, error) {
		publicServer := server.New(
			logger,
			fmt.Sprintf(":%d", cfg.BackendApp.HTTP.Port),
			server.WithCORS(cfg.BackendApp.HTTP.CorsAllowOrigins),
			server.WithTrustedProxies([]string{cfg.Global.ServerIP}),
			server.WithLogger(cfg.BackendApp.Base.LogHTTP),
		)

		gw := public.New(
			logger,
			publicServer,
			&controller.Controls{
				AuthTenant:   authTenantCtrl,
				CRM:          crmCtrl,
				Tenant:       tenantCtrl,
				Registration: registrationCtrl,
				Rooms:        roomsCtrl,
				UsersTenant:  usersTenantCtrl,
				Testing:      testingCtrl,
			},
		)

		err := gw.RegisterHandlers(context.Background(), authDeliveryCtrl)
		if err != nil {
			return nil, nil, err
		}

		return publicServer, gw, nil
	}),
)
