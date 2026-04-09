package auth

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/neurochar/backend/internal/app/config"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/controller"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware/limiter"
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
		pkg:          "httpController.Auth",
		vldtr:        validation.New(),
		cfg:          cfg,
		backoff:      backoff,
		fileUC:       fileUC,
		tenantFacade: tenantFacade,
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

	routeGroup.Get("/whoiam", cpanelMdwr.AuthFullCheck, ctrl.WhoIAmHandler)

	routeGroup.Post("/logout", cpanelMdwr.AuthFullCheck, ctrl.LogoutHandler)

	routeGroup.Post("/password-recovery", ipLimiterMiddleware, ctrl.RequestPasswordRecoveryHandler)

	routeGroup.Post("/check-code", ipLimiterMiddleware, ctrl.CheckAccountCodeHandler)

	routeGroup.Post("/password-by-code", ipLimiterMiddleware, ctrl.UpdatePasswordByCodeHandler)

	routeGroup.Post("/verify-email", ipLimiterMiddleware, ctrl.AccountVerifyEmailHandler)
}
