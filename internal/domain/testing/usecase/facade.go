package usecase

type Facade struct {
	PersonalityTrait PersonalityTraitUsecase
	Profile          ProfileUsecase
	Room             RoomUsecase
	Cross            CrossUsecase
}

func NewFacade(
	personalityUC PersonalityTraitUsecase,
	candidateUC ProfileUsecase,
	roomUC RoomUsecase,
	crossUC CrossUsecase,
) *Facade {
	return &Facade{
		PersonalityTrait: personalityUC,
		Profile:          candidateUC,
		Room:             roomUC,
		Cross:            crossUC,
	}
}
