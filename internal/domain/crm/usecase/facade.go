package usecase

type Facade struct {
	Candidate CandidateUsecase
	Cross     CrossUsecase
}

func NewFacade(
	candidateUC CandidateUsecase,
	crossUC CrossUsecase,
) *Facade {
	return &Facade{
		Candidate: candidateUC,
		Cross:     crossUC,
	}
}
