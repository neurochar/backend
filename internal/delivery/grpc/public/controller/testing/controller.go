package testing

import (
	"time"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/grpc/public"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
)

type Controller struct {
	desc.UnimplementedTestingPublicServiceServer
	pkg           string
	cfg           config.Config
	backoff       *backoff.Controller
	limiter       *limiter.Controller
	server        *public.PublicServer
	crmFacade     *crmUC.Facade
	testingFacade *testingUC.Facade
	fileUC        fileUC.Usecase
}

func New(
	cfg config.Config,
	backoff *backoff.Controller,
	limiter *limiter.Controller,
	server *public.PublicServer,
	crmFacade *crmUC.Facade,
	testingFacade *testingUC.Facade,
	fileUC fileUC.Usecase,
) *Controller {
	ctrl := &Controller{
		pkg:           "grpc.Controller.Testing",
		cfg:           cfg,
		backoff:       backoff,
		limiter:       limiter,
		server:        server,
		crmFacade:     crmFacade,
		testingFacade: testingFacade,
		fileUC:        fileUC,
	}

	return ctrl
}

const backoffConfigLLMGroupID = "controller.llm"

func (ctrl *Controller) Register() {
	desc.RegisterTestingPublicServiceServer(ctrl.server.Server().GRPCServer(), ctrl)

	ctrl.backoff.SetConfigForGroup(
		backoffConfigLLMGroupID,
		backoff.WithTtl(time.Minute*10),
		backoff.WithInitialInterval(time.Second),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*1),
	)
}
