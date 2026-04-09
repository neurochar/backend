package rooms

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/controller"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg           string
	vldtr         *validator.Validate
	cfg           config.Config
	backoff       *backoff.Controller
	fileUC        fileUC.Usecase
	tenantFacade  *tenantUC.Facade
	testingFacade *testingUC.Facade
}

func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	fileUC fileUC.Usecase,
	tenantFacade *tenantUC.Facade,
	testingFacade *testingUC.Facade,
) *Controller {
	controller := &Controller{
		pkg:           "httpController.Rooms",
		vldtr:         validation.New(),
		cfg:           cfg,
		backoff:       backoff,
		fileUC:        fileUC,
		tenantFacade:  tenantFacade,
		testingFacade: testingFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "rooms"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.Get("/:id<guid>", ctrl.GetRoomHandler)
	routeGroup.Post("/:id<guid>", ctrl.FinishRoomHandler)
}
