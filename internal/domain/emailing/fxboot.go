package emailing

import (
	"go.uber.org/fx"

	itemRepo "github.com/neurochar/backend/internal/domain/emailing/repository/item"
	"github.com/neurochar/backend/internal/domain/emailing/usecase"
)

// FxModule - fx module
var FxModule = fx.Module(
	"emailing_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(itemRepo.NewRepository, fx.As(new(usecase.ItemRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewUsecaseImpl, fx.As(new(usecase.Usecase))),
	),
	fx.Invoke(Init),
)
