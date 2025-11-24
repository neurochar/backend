package testing

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
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
	testingFacade *testingUC.Facade
}

func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	fileUC fileUC.Usecase,
	testingFacade *testingUC.Facade,
) *Controller {
	controller := &Controller{
		pkg:           "httpController.Testing",
		vldtr:         validation.New(),
		cfg:           cfg,
		backoff:       backoff,
		fileUC:        fileUC,
		testingFacade: testingFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "testing"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url), cpanelMdwr.AuthRequired)

	routeGroup.Get("/personality_traits", ctrl.ListPersonalityTraitsHandler)

	routeGroup.Get("/profiles", ctrl.ListProfilesHandler)
	routeGroup.Post("/profiles", cpanelMdwr.AuthFullCheck, ctrl.CreateProfileHandler)
	routeGroup.Get("/profiles/:id<guid>", ctrl.GetProfileHandler)
	routeGroup.Patch("/profiles/:id<guid>", cpanelMdwr.AuthFullCheck, ctrl.PatchProfileHandler)
	routeGroup.Delete("/profiles/:id<guid>", cpanelMdwr.AuthFullCheck, ctrl.DeleteProfileHandler)
}
