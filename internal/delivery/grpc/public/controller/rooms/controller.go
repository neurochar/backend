package rooms

import (
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/grpc/public"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	testiongUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/rooms/v1"
)

type Controller struct {
	desc.UnimplementedRoomsPublicServiceServer
	pkg           string
	cfg           config.Config
	backoff       *backoff.Controller
	limiter       *limiter.Controller
	server        *public.PublicServer
	tenantFacade  *tenantUC.Facade
	testingFacade *testiongUC.Facade
}

func New(
	cfg config.Config,
	backoff *backoff.Controller,
	limiter *limiter.Controller,
	server *public.PublicServer,
	tenantFacade *tenantUC.Facade,
	testingFacade *testiongUC.Facade,
) *Controller {
	ctrl := &Controller{
		pkg:           "grpc.Controller.Rooms",
		cfg:           cfg,
		backoff:       backoff,
		limiter:       limiter,
		server:        server,
		tenantFacade:  tenantFacade,
		testingFacade: testingFacade,
	}

	return ctrl
}

func (ctrl *Controller) Register() {
	desc.RegisterRoomsPublicServiceServer(ctrl.server.Server().GRPCServer(), ctrl)
}
