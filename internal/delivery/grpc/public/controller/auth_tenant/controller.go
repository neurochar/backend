package auth_tenant

import (
	"time"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/grpc/public"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

type Controller struct {
	desc.UnimplementedAuthTenantPublicServiceServer
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
		pkg:          "grpc.Controller.AuthTenant",
		cfg:          cfg,
		backoff:      backoff,
		limiter:      limiter,
		server:       server,
		tenantFacade: tenantFacade,
		fileUC:       fileUC,
	}

	return ctrl
}

const backoffConfigAuthGroupID = "controller.auth"

const backoffConfigPasswordRecoveryGroupID = "controller.password_recovery"

func (ctrl *Controller) Register() {
	desc.RegisterAuthTenantPublicServiceServer(ctrl.server.Server().GRPCServer(), ctrl)

	ctrl.backoff.SetConfigForGroup(
		backoffConfigAuthGroupID,
		backoff.WithTtl(time.Minute*10),
		backoff.WithInitialInterval(time.Second*5),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*1),
	)

	ctrl.backoff.SetConfigForGroup(
		backoffConfigPasswordRecoveryGroupID,
		backoff.WithTtl(time.Minute*30),
		backoff.WithInitialInterval(time.Second*30),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*10),
	)
}
