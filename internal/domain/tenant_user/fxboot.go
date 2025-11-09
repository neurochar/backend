package tenant_user

import (
	"go.uber.org/fx"

	accountRepo "github.com/neurochar/backend/internal/domain/tenant_user/repository/pg/account"
	accountCodeRepo "github.com/neurochar/backend/internal/domain/tenant_user/repository/pg/account_code"
	sessionRepo "github.com/neurochar/backend/internal/domain/tenant_user/repository/pg/session"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	accountUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase/account"
	authUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase/auth"
)

// FxModule - fx module
var FxModule = fx.Module(
	"tenant_user_module",

	// repositories
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
		fx.Annotate(accountUC.NewUsecaseImpl, fx.As(new(usecase.AccountUsecase))),
	),
	fx.Provide(
		fx.Annotate(authUC.NewUsecaseImpl, fx.As(new(usecase.AuthUsecase))),
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
