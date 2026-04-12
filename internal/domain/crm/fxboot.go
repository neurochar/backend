package crm

import (
	"go.uber.org/fx"

	candidateRepo "github.com/neurochar/backend/internal/domain/crm/repository/pg/candidate"
	candidateResumeRepo "github.com/neurochar/backend/internal/domain/crm/repository/pg/candidate_resume"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase/candidate"
	candidateResumeUC "github.com/neurochar/backend/internal/domain/crm/usecase/candidate_resume"
	crossUC "github.com/neurochar/backend/internal/domain/crm/usecase/cross"
)

// FxModule - fx module
var FxModule = fx.Module(
	"crm_module",

	// repositories
	fx.Provide(
		fx.Private,
		fx.Annotate(candidateRepo.NewRepository, fx.As(new(usecase.CandidateRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(candidateResumeRepo.NewRepository, fx.As(new(usecase.CandidateResumeRepository))),
	),

	// usecases
	fx.Provide(
		fx.Annotate(candidateUC.NewUsecaseImpl, fx.As(new(usecase.CandidateUsecase))),
	),
	fx.Provide(
		fx.Annotate(candidateResumeUC.NewUsecaseImpl, fx.As(new(usecase.CandidateResumeUsecase))),
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
