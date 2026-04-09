package users_tenant

import (
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/grpc/public"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
)

type Controller struct {
	desc.UnimplementedUsersTenantPublicServiceServer
	pkg          string
	cfg          config.Config
	backoff      *backoff.Controller
	limiter      *limiter.Controller
	server       *public.PublicServer
	tenantFacade *tenantUC.Facade
	fileUC       fileUC.Usecase
}

func New(
	cfg config.Config,
	backoff *backoff.Controller,
	limiter *limiter.Controller,
	server *public.PublicServer,
	tenantFacade *tenantUC.Facade,
	fileUC fileUC.Usecase,
) *Controller {
	ctrl := &Controller{
		pkg:          "grpc.Controller.UsersTenant",
		cfg:          cfg,
		backoff:      backoff,
		limiter:      limiter,
		server:       server,
		tenantFacade: tenantFacade,
		fileUC:       fileUC,
	}

	return ctrl
}

func (ctrl *Controller) Register() {
	desc.RegisterUsersTenantPublicServiceServer(ctrl.server.Server().GRPCServer(), ctrl)
}
