// Package users contains users http controller
package users

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware/limiter"
	v1 "github.com/neurochar/backend/internal/delivery/http/cpanel/v1"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/validation"

	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	userConstants "github.com/neurochar/backend/internal/domain/user/constants"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

// Controller - auth controller
type Controller struct {
	pkg        string
	vldtr      *validator.Validate
	cfg        config.Config
	backoff    *backoff.Controller
	userFacade *userUC.Facade
	fileUC     fileUC.Usecase
}

// NewController - create new auth controller
func NewController(
	cfg config.Config,
	backoff *backoff.Controller,
	userFacade *userUC.Facade,
	fileUC fileUC.Usecase,
) *Controller {
	controller := &Controller{
		pkg:        "httpController.Users",
		vldtr:      validation.New(),
		cfg:        cfg,
		backoff:    backoff,
		userFacade: userFacade,
		fileUC:     fileUC,
	}
	return controller
}

const backoffConfigAuthGroupID = "http.auth"

const backoffConfigPasswordRecoveryGroupID = "http.password_recovery"

// RegisterRoutes - register auth routes
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

	checkGlobalSettingsAccessMiddleware := cpanelMdwr.MiddlewareCheckRight(userConstants.RightKeyAccessToGlobalSettings, 1)

	ipLimiterMiddleware := groups.RateLimiter.Get(limiter.DefaultName).Create(true, false, "")
	accountLimiterMiddleware := groups.RateLimiter.Get(limiter.DefaultName).Create(false, true, "")

	routeGroup := groups.Default.Group("/users")
	accountsRouteGroup := routeGroup.Group("/accounts")
	profilesRouteGroup := routeGroup.Group("/profiles")
	rolesRouteGroup := routeGroup.Group("/roles", checkGlobalSettingsAccessMiddleware)

	// Создать пользователя
	routeGroup.Post("", checkGlobalSettingsAccessMiddleware, ctrl.CreateUserHandler)

	// Список пользователей
	routeGroup.Get("", checkGlobalSettingsAccessMiddleware, ctrl.ListUsersHandler)

	// Получить пользователя
	routeGroup.Get("/:profile_id<int>", checkGlobalSettingsAccessMiddleware, ctrl.GetUserHandler)

	// Получить текущего пользователя
	routeGroup.Get("/auth", ctrl.MeHandler)

	// Авторизация
	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/users/auth", groups.Prefix), fiber.MethodPost)
	routeGroup.Post("/auth", ipLimiterMiddleware, ctrl.LoginHandler)

	// Logout
	routeGroup.Post("/logout", ctrl.LogoutHandler)

	// Запрос на восстановление пароля
	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/users/password-recovery", groups.Prefix), fiber.MethodPost)
	routeGroup.Post("/password-recovery", ipLimiterMiddleware, ctrl.RequestPasswordRecoveryHandler)

	// Обновить пароль текущего пользователя
	accountsRouteGroup.Put("/me/password", ctrl.UpdateMyPasswordHandler)

	// Обновить email текущего пользователя
	accountsRouteGroup.Put("/me/email", ctrl.UpdateMyEmailHandler)

	// Проверить код аккаунта
	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/users/accounts/check-code", groups.Prefix), fiber.MethodPost)
	accountsRouteGroup.Post("/check-code", ctrl.CheckAccountCodeHandler)

	// Проверить код аккаунта
	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/users/accounts/password-by-code", groups.Prefix), fiber.MethodPost)
	accountsRouteGroup.Post("/password-by-code", ctrl.UpdatePasswordByCodeHandler)

	// Подтвердить email аккаунта по коду
	cpanelMdwr.AddAuthErrSkiping(fmt.Sprintf("%s/users/accounts/verify-email", groups.Prefix), fiber.MethodGet)
	accountsRouteGroup.Get("/verify-email", ipLimiterMiddleware, ctrl.AccountVerifyEmailHandler)

	// Пропатчить аккаунт
	accountsRouteGroup.Patch("/:id<guid>", checkGlobalSettingsAccessMiddleware, ctrl.PatchAccountHandler)

	// Обновить мой профиль
	profilesRouteGroup.Put("/me", ctrl.UpdateMyProfileHandler)

	// Обновить профиль
	profilesRouteGroup.Put("/:id<int>", checkGlobalSettingsAccessMiddleware, ctrl.UpdateProfileHandler)

	// Загрузить фото профиля
	profilesRouteGroup.Post("/photo_file", accountLimiterMiddleware, ctrl.UploadPhotoFileHandler)

	// Список ролей
	rolesRouteGroup.Get("", ctrl.ListRolesHandler)

	// Получить роль
	rolesRouteGroup.Get("/:id<int>", ctrl.GetRoleHandler)

	// Создать роль
	rolesRouteGroup.Post("", ctrl.CreateRoleHandler)

	// Обновить роль
	rolesRouteGroup.Put("/:id<int>", ctrl.UpdateRoleHandler)

	// Удалить роль
	rolesRouteGroup.Delete("/:id<int>", ctrl.DeleteRoleHandler)
}
