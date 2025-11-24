package usecase

type Facade struct {
	PersonalityTrait PersonalityTraitUsecase
	Profile          ProfileUsecase
	Cross            CrossUsecase
}

func NewFacade(
	personalityUC PersonalityTraitUsecase,
	candidateUC ProfileUsecase,
	crossUC CrossUsecase,
) *Facade {
	return &Facade{
		PersonalityTrait: personalityUC,
		Profile:          candidateUC,
		Cross:            crossUC,
	}
}
