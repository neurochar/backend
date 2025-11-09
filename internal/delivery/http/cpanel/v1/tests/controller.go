package tests

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/cpanel/v1"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg   string
	vldtr *validator.Validate
	cfg   config.Config
}

func NewController(
	cfg config.Config,
) *Controller {
	controller := &Controller{
		pkg:   "httpController.Tests",
		vldtr: validation.New(),
		cfg:   cfg,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	const url = "tests"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/tests/internal-error", groups.Prefix), fiber.MethodPost)
	routeGroup.Post("/internal-error", ctrl.InternalErrorHandler)

	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/tests/panic-error", groups.Prefix), fiber.MethodPost)
	routeGroup.Post("/panic-error", ctrl.PanicErrorHandler)
}
