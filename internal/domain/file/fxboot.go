package file

import (
	"go.uber.org/fx"

	fileRepo "github.com/neurochar/backend/internal/domain/file/repository/file"
	"github.com/neurochar/backend/internal/domain/file/usecase"
)

// FxModule - fx module
var FxModule = fx.Module(
	"file_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(fileRepo.NewRepository, fx.As(new(usecase.FileRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewUsecaseImpl, fx.As(new(usecase.Usecase))),
	),
	fx.Invoke(Init),
)
