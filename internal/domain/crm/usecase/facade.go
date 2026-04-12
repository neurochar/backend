package usecase

type Facade struct {
	Candidate       CandidateUsecase
	CandidateResume CandidateResumeUsecase
	Cross           CrossUsecase
}

func NewFacade(
	candidateUC CandidateUsecase,
	candidateResumeUC CandidateResumeUsecase,
	crossUC CrossUsecase,
) *Facade {
	return &Facade{
		Candidate:       candidateUC,
		CandidateResume: candidateResumeUC,
		Cross:           crossUC,
	}
}
