package middleware

import userUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"

type Controller struct {
	authUC userUC.AuthUsecase
}

func New(authUC userUC.AuthUsecase) *Controller {
	return &Controller{
		authUC: authUC,
	}
}
