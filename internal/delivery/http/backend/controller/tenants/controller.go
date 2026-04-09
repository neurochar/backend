package tenants

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/controller"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg          string
	vldtr        *validator.Validate
	cfg          config.Config
	backoff      *backoff.Controller
	fileUC       fileUC.Usecase
	tenantFacade *tenantUC.Facade
}

func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	fileUC fileUC.Usecase,
	tenantFacade *tenantUC.Facade,
) *Controller {
	controller := &Controller{
		pkg:          "httpController.Tenants",
		vldtr:        validation.New(),
		cfg:          cfg,
		backoff:      backoff,
		fileUC:       fileUC,
		tenantFacade: tenantFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "tenants"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.Get("/is_exists/:text_id<string>", ctrl.GetIsExistsTenantHandler)

	routeGroup.Patch("", cpanelMdwr.AuthFullCheck, ctrl.PatchTenantHandler)
}
