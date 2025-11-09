package usecase

type Facade struct {
	Account   AccountUsecase
	AdminAuth AdminAuthUsecase
	Common    UserUsecase
	Profile   ProfileUsecase
	Role      RoleUsecase
}

func NewFacade(
	accountUC AccountUsecase,
	adminAuthUC AdminAuthUsecase,
	commonUC UserUsecase,
	profileUC ProfileUsecase,
	roleUC RoleUsecase,
) *Facade {
	return &Facade{
		Account:   accountUC,
		AdminAuth: adminAuthUC,
		Common:    commonUC,
		Profile:   profileUC,
		Role:      roleUC,
	}
}
