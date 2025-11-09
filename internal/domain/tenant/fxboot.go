package tenant

import (
	"go.uber.org/fx"

	tenantRepo "github.com/neurochar/backend/internal/domain/tenant/repository/pg/tenant"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
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

	// usecases
	fx.Provide(
		fx.Annotate(tenantUC.NewUsecaseImpl, fx.As(new(usecase.TenantUsecase))),
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
