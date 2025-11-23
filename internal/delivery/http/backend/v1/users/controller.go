package users

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware/limiter"
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
		pkg:          "httpController.Users",
		vldtr:        validation.New(),
		cfg:          cfg,
		backoff:      backoff,
		fileUC:       fileUC,
		tenantFacade: tenantFacade,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	accountLimiterMiddleware := groups.RateLimiter.Get(limiter.DefaultName).Create(false, true, "")

	const url = "users"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url), cpanelMdwr.AuthRequired)

	routeGroup.Post("/photo_file", accountLimiterMiddleware, cpanelMdwr.AuthFullCheck, ctrl.UploadPhotoFileHandler)

	routeGroup.Put("/my_profile", cpanelMdwr.AuthFullCheck, ctrl.UpdateMyProfileHandler)

	routeGroup.Put("/my_password", cpanelMdwr.AuthFullCheck, ctrl.UpdateMyPasswordHandler)

	routeGroup.Get("/:id<guid>", ctrl.GetAccountHandler)

	routeGroup.Post("", cpanelMdwr.AuthFullCheck, ctrl.CreateAccountHandler)

	routeGroup.Patch("/:id<guid>", cpanelMdwr.AuthFullCheck, ctrl.PatchProfileHandler)

	routeGroup.Get("", ctrl.ListAccountsHandler)
}
