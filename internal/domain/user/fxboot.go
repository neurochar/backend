package user

import (
	"go.uber.org/fx"

	accountRepo "github.com/neurochar/backend/internal/domain/user/repository/account"
	accountCodeRepo "github.com/neurochar/backend/internal/domain/user/repository/account_code"
	adminSessionRepo "github.com/neurochar/backend/internal/domain/user/repository/admin_session"
	profileRepo "github.com/neurochar/backend/internal/domain/user/repository/profile"
	profileAccountRepo "github.com/neurochar/backend/internal/domain/user/repository/profile__account"
	roleRepo "github.com/neurochar/backend/internal/domain/user/repository/role"
	roleToRightRepo "github.com/neurochar/backend/internal/domain/user/repository/role_to_right"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	accountUC "github.com/neurochar/backend/internal/domain/user/usecase/account"
	adminAuthUC "github.com/neurochar/backend/internal/domain/user/usecase/admin_auth"
	profileUC "github.com/neurochar/backend/internal/domain/user/usecase/profile"
	roleUC "github.com/neurochar/backend/internal/domain/user/usecase/role"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase/user"
)

// FxModule - fx module
var FxModule = fx.Module(
	"user_module",

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
		fx.Annotate(adminSessionRepo.NewRepository, fx.As(new(usecase.SessionRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(profileRepo.NewRepository, fx.As(new(usecase.ProfileRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(roleRepo.NewRepository, fx.As(new(usecase.RoleRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(roleToRightRepo.NewRepository, fx.As(new(usecase.RoleToRightRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(profileAccountRepo.NewRepository, fx.As(new(usecase.ProfileAccountRepository))),
	),

	// usecases
	fx.Provide(
		fx.Annotate(accountUC.NewUsecaseImpl, fx.As(new(usecase.AccountUsecase))),
	),
	fx.Provide(
		fx.Annotate(adminAuthUC.NewUsecaseImpl, fx.As(new(usecase.AdminAuthUsecase))),
	),
	fx.Provide(
		fx.Annotate(profileUC.NewUsecaseImpl, fx.As(new(usecase.ProfileUsecase))),
	),
	fx.Provide(
		fx.Annotate(roleUC.NewUsecaseImpl, fx.As(new(usecase.RoleUsecase))),
	),
	fx.Provide(
		fx.Annotate(userUC.NewUsecaseImpl, fx.As(new(usecase.UserUsecase))),
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
