package auth

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware/limiter"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/validation"
)

type Controller struct {
	pkg              string
	vldtr            *validator.Validate
	cfg              config.Config
	backoff          *backoff.Controller
	fileUC           fileUC.Usecase
	tenantFacade     *tenantUC.Facade
	tenantUserFacade *tenantUserUC.Facade
}

func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	fileUC fileUC.Usecase,
	tenantFacade *tenantUC.Facade,
	tenantUserFacade *tenantUserUC.Facade,
) *Controller {
	controller := &Controller{
		pkg:              "httpController.Auth",
		vldtr:            validation.New(),
		cfg:              cfg,
		backoff:          backoff,
		fileUC:           fileUC,
		tenantFacade:     tenantFacade,
		tenantUserFacade: tenantUserFacade,
	}
	return controller
}

const backoffConfigAuthGroupID = "http.auth"

const backoffConfigPasswordRecoveryGroupID = "http.password_recovery"

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, cpanelMdwr *middleware.Controller) {
	ctrl.backoff.SetConfigForGroup(
		backoffConfigAuthGroupID,
		backoff.WithTtl(time.Minute*10),
		backoff.WithInitialInterval(time.Second*5),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*1),
	)

	ctrl.backoff.SetConfigForGroup(
		backoffConfigPasswordRecoveryGroupID,
		backoff.WithTtl(time.Minute*30),
		backoff.WithInitialInterval(time.Second*30),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*10),
	)

	ipLimiterMiddleware := groups.RateLimiter.Get(limiter.DefaultName).Create(true, false, "")

	const url = "auth"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.Post("/login", ipLimiterMiddleware, ctrl.LoginHandler)

	routeGroup.Post("/refresh", ctrl.RefreshHandler)

	routeGroup.Get("/whoiam", cpanelMdwr.MiddlewareAuthRequired, ctrl.WhoIAmHandler)

	routeGroup.Post("/logout", cpanelMdwr.MiddlewareAuthRequired, ctrl.LogoutHandler)
}
