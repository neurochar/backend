package tests

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg          string
	vldtr        *validator.Validate
	cfg          config.Config
	tenantFacade *tenantUC.Facade
}

func NewController(
	cfg config.Config,
	tenantFacade *tenantUC.Facade,
) *Controller {
	controller := &Controller{
		pkg:          "httpController.Tests",
		vldtr:        validation.New(),
		cfg:          cfg,
		tenantFacade: tenantFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "tests"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.All("/internal-error", ctrl.InternalErrorHandler)

	routeGroup.Post("/panic-error", ctrl.PanicErrorHandler)
}
