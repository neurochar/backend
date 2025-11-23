package usecase

type Facade struct {
	Tenant       TenantUsecase
	Registration RegistrationUsecase
	Account      AccountUsecase
	Session      SessionUsecase
	Auth         AuthUsecase
	Cross        CrossUsecase
}

func NewFacade(
	tenantUC TenantUsecase,
	registrationUC RegistrationUsecase,
	accountUC AccountUsecase,
	sessionUC SessionUsecase,
	authUC AuthUsecase,
	crossUC CrossUsecase,
) *Facade {
	return &Facade{
		Tenant:       tenantUC,
		Registration: registrationUC,
		Account:      accountUC,
		Session:      sessionUC,
		Auth:         authUC,
		Cross:        crossUC,
	}
}
