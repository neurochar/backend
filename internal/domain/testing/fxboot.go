package testing

import (
	"go.uber.org/fx"

	profileRepo "github.com/neurochar/backend/internal/domain/testing/repository/pg/profile"
	roomRepo "github.com/neurochar/backend/internal/domain/testing/repository/pg/room"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	crossUC "github.com/neurochar/backend/internal/domain/testing/usecase/cross"
	personalityTraitUC "github.com/neurochar/backend/internal/domain/testing/usecase/personality_trait"
	profileUC "github.com/neurochar/backend/internal/domain/testing/usecase/profile"
	roomUC "github.com/neurochar/backend/internal/domain/testing/usecase/room"
)

// FxModule - fx module
var FxModule = fx.Module(
	"testing_module",

	// repositories
	fx.Provide(
		fx.Private,
		fx.Annotate(profileRepo.NewRepository, fx.As(new(usecase.ProfileRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(roomRepo.NewRepository, fx.As(new(usecase.RoomRepository))),
	),

	// usecases
	fx.Provide(
		fx.Annotate(profileUC.NewUsecaseImpl, fx.As(new(usecase.ProfileUsecase))),
	),
	fx.Provide(
		fx.Annotate(personalityTraitUC.NewUsecaseImpl, fx.As(new(usecase.PersonalityTraitUsecase))),
	),
	fx.Provide(
		fx.Annotate(roomUC.NewUsecaseImpl, fx.As(new(usecase.RoomUsecase))),
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
