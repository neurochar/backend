package registration

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
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
		pkg:          "httpController.Registration",
		vldtr:        validation.New(),
		cfg:          cfg,
		backoff:      backoff,
		fileUC:       fileUC,
		tenantFacade: tenantFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "registration"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.Post("/finish", ctrl.FinishRegistrationHandler)

	routeGroup.Post("", ctrl.StartRegistrationHandler)

	routeGroup.Get("/:id<guid>", ctrl.CheckRegistrationHandler)
}
