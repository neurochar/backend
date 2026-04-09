package tenant

import (
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/grpc/public"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/tenant/v1"
)

type Controller struct {
	desc.UnimplementedTenantPublicServiceServer
	pkg          string
	cfg          config.Config
	backoff      *backoff.Controller
	limiter      *limiter.Controller
	server       *public.PublicServer
	tenantFacade *tenantUC.Facade
}

func New(
	cfg config.Config,
	backoff *backoff.Controller,
	limiter *limiter.Controller,
	server *public.PublicServer,
	tenantFacade *tenantUC.Facade,
) *Controller {
	ctrl := &Controller{
		pkg:          "grpc.Controller.Tenant",
		cfg:          cfg,
		backoff:      backoff,
		limiter:      limiter,
		server:       server,
		tenantFacade: tenantFacade,
	}

	return ctrl
}

func (ctrl *Controller) Register() {
	desc.RegisterTenantPublicServiceServer(ctrl.server.Server().GRPCServer(), ctrl)
}
