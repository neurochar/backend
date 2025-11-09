package middleware

import userUC "github.com/neurochar/backend/internal/domain/user/usecase"

type Controller struct {
	authAdminUC userUC.AdminAuthUsecase
	skipAuth    map[string]struct{}
}

func New(authAdminUC userUC.AdminAuthUsecase) *Controller {
	return &Controller{
		authAdminUC: authAdminUC,
		skipAuth:    make(map[string]struct{}),
	}
}
