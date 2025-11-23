package tenant

import (
	"go.uber.org/fx"

	accountRepo "github.com/neurochar/backend/internal/domain/tenant/repository/pg/account"
	accountCodeRepo "github.com/neurochar/backend/internal/domain/tenant/repository/pg/account_code"
	registrationRepo "github.com/neurochar/backend/internal/domain/tenant/repository/pg/registration"
	sessionRepo "github.com/neurochar/backend/internal/domain/tenant/repository/pg/session"
	tenantRepo "github.com/neurochar/backend/internal/domain/tenant/repository/pg/tenant"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	accountUC "github.com/neurochar/backend/internal/domain/tenant/usecase/account"
	authUC "github.com/neurochar/backend/internal/domain/tenant/usecase/auth"
	crossUC "github.com/neurochar/backend/internal/domain/tenant/usecase/cross"
	registrationUC "github.com/neurochar/backend/internal/domain/tenant/usecase/registration"
	sessionUC "github.com/neurochar/backend/internal/domain/tenant/usecase/session"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase/tenant"
)

// FxModule - fx module
var FxModule = fx.Module(
	"tenant_module",

	// repositories
	fx.Provide(
		fx.Private,
		fx.Annotate(tenantRepo.NewRepository, fx.As(new(usecase.TenantRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(registrationRepo.NewRepository, fx.As(new(usecase.RegistrationRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(accountRepo.NewRepository, fx.As(new(usecase.AccountRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(accountCodeRepo.NewRepository, fx.As(new(usecase.AccountCodeRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(sessionRepo.NewRepository, fx.As(new(usecase.SessionRepository))),
	),

	// usecases
	fx.Provide(
		fx.Annotate(tenantUC.NewUsecaseImpl, fx.As(new(usecase.TenantUsecase))),
	),
	fx.Provide(
		fx.Annotate(registrationUC.NewUsecaseImpl, fx.As(new(usecase.RegistrationUsecase))),
	),
	fx.Provide(
		fx.Annotate(accountUC.NewUsecaseImpl, fx.As(new(usecase.AccountUsecase))),
	),
	fx.Provide(
		fx.Annotate(authUC.NewUsecaseImpl, fx.As(new(usecase.AuthUsecase))),
	),
	fx.Provide(
		fx.Annotate(sessionUC.NewUsecaseImpl, fx.As(new(usecase.SessionUsecase))),
	),
	fx.Provide(
		fx.Annotate(crossUC.NewUsecaseImpl, fx.As(new(usecase.CrossUsecase))),
	),

	// facade
	fx.Provide(
		usecase.NewFacade,
	),

	// init
	fx.Provide(
		fx.Annotate(Init, fx.ResultTags(`group:"InvokeInit"`)),
	),
)
