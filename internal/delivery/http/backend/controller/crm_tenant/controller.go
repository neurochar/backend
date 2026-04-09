package crm_tenant

import (
	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller"
	"github.com/neurochar/backend/internal/delivery/http/backend/gateway"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg              string
	vldtr            *validator.Validate
	cfg              config.Config
	backoff          *backoff.Controller
	fileUC           fileUC.Usecase
	crmFacade        *crmUC.Facade
	crmTenantService desc.CrmTenantPublicServiceClient
}

func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	fileUC fileUC.Usecase,
	crmFacade *crmUC.Facade,
	gatewayClient *gateway.DeliveryGrpcClient,
) *Controller {
	crmTenantService := desc.NewCrmTenantPublicServiceClient(gatewayClient.Connection)

	controller := &Controller{
		pkg:              "httpController.CRM",
		vldtr:            validation.New(),
		cfg:              cfg,
		backoff:          backoff,
		fileUC:           fileUC,
		crmFacade:        crmFacade,
		crmTenantService: crmTenantService,
	}

	return controller
}

func RegisterRoutes(groups *controller.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	routeGroup := groups.Default.Group("/v1/tenant/crm")

	// лимитеры добавить в грпс потом
	routeGroup.Post("/candidates-resume", ctrl.UploadCandidateResumeFileHandler)
}
