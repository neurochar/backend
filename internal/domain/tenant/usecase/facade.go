package usecase

type Facade struct {
	Tenant TenantUsecase
}

func NewFacade(
	tenantUC TenantUsecase,
) *Facade {
	return &Facade{
		Tenant: tenantUC,
	}
}
