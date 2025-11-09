package usecase

type Facade struct {
	Account AccountUsecase
	Auth    AuthUsecase
}

func NewFacade(
	accountUC AccountUsecase,
	authUC AuthUsecase,
) *Facade {
	return &Facade{
		Account: accountUC,
		Auth:    authUC,
	}
}
