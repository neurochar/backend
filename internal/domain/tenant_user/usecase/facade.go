package usecase

type Facade struct {
	Common  CommonUsecase
	Account AccountUsecase
	Auth    AuthUsecase
}

func NewFacade(
	commonUC CommonUsecase,
	accountUC AccountUsecase,
	authUC AuthUsecase,
) *Facade {
	return &Facade{
		Common:  commonUC,
		Account: accountUC,
		Auth:    authUC,
	}
}
