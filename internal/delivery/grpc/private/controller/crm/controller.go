package crm

import (
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/grpc/private"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/private/crm/v1"
)

type Controller struct {
	desc.UnimplementedCrmPrivateServiceServer
	pkg       string
	cfg       config.Config
	server    *private.PrivateServer
	crmFacade *crmUC.Facade
}

func New(
	cfg config.Config,
	server *private.PrivateServer,
	crmFacade *crmUC.Facade,
) *Controller {
	ctrl := &Controller{
		pkg:       "grpc.Private.Controller.Crm",
		cfg:       cfg,
		server:    server,
		crmFacade: crmFacade,
	}

	return ctrl
}

func (ctrl *Controller) Register() {
	desc.RegisterCrmPrivateServiceServer(ctrl.server.Server().GRPCServer(), ctrl)
}
