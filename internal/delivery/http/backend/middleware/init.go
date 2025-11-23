package middleware

import tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"

type Controller struct {
	authUC tenantUC.AuthUsecase
}

func New(authUC tenantUC.AuthUsecase) *Controller {
	return &Controller{
		authUC: authUC,
	}
}
