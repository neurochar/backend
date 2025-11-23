package crm

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg       string
	vldtr     *validator.Validate
	cfg       config.Config
	backoff   *backoff.Controller
	fileUC    fileUC.Usecase
	crmFacade *crmUC.Facade
}

func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	fileUC fileUC.Usecase,
	crmFacade *crmUC.Facade,
) *Controller {
	controller := &Controller{
		pkg:       "httpController.CRM",
		vldtr:     validation.New(),
		cfg:       cfg,
		backoff:   backoff,
		fileUC:    fileUC,
		crmFacade: crmFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "crm"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url), cpanelMdwr.AuthRequired)

	routeGroup.Get("/candidates", ctrl.ListCandidatesHandler)
	routeGroup.Post("/candidates", cpanelMdwr.AuthFullCheck, ctrl.CreateCandidateHandler)
	routeGroup.Get("/candidates/:id<guid>", ctrl.GetCandidateHandler)
	routeGroup.Patch("/candidates/:id<guid>", ctrl.PatchCandidateHandler)
	routeGroup.Delete("/candidates/:id<guid>", ctrl.DeleteCandidateHandler)
}
