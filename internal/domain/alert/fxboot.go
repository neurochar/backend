package alert

import (
	"go.uber.org/fx"

	itemRepo "github.com/neurochar/backend/internal/domain/alert/repository/tg"
	"github.com/neurochar/backend/internal/domain/alert/usecase"
)

// FxModule - fx module
var FxModule = fx.Module(
	"alert_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(itemRepo.NewRepository, fx.As(new(usecase.TelegramRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewUsecaseImpl, fx.As(new(usecase.Usecase))),
	),
	fx.Invoke(Init),
)
