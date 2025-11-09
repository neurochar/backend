package middleware

import (
	"fmt"
	"net"

	"github.com/gofiber/fiber/v2"

	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/neurochar/backend/internal/infra/loghandler"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

type AuthData struct {
	Session *userEntity.AdminSession
	Account *userEntity.Account
	Role    *userUC.RoleDTO
}

const AuthDataKey = "auth_data"

func (ctrl *Controller) AddAuthErrSkiping(url string, method string) {
	ctrl.skipAuth[fmt.Sprintf("%s_%s", url, method)] = struct{}{}
}

func (ctrl *Controller) MiddlewareAuth(c *fiber.Ctx) error {
	skipErr := false
	checkKey := fmt.Sprintf("%s_%s", c.Path(), c.Method())
	if _, ok := ctrl.skipAuth[checkKey]; ok {
		skipErr = true
	}

	cookieValue := c.Cookies("auth_admin_session")

	requestIP := net.ParseIP(middleware.GetRealIP(c))

	session, account, role, err := ctrl.authAdminUC.AuthByJWT(c.Context(), cookieValue, requestIP)
	if err != nil {
		if skipErr {
			return c.Next()
		}
		return err
	}

	c.Locals(AuthDataKey, &AuthData{
		Account: account,
		Session: session,
		Role:    role,
	})

	ctxData, ctxKey := loghandler.SetData(c.Context(), "request.account.id", account.ID)
	c.Locals(ctxKey, ctxData)

	return c.Next()
}

func GetAuthData(c *fiber.Ctx) *AuthData {
	data, ok := c.Locals(AuthDataKey).(*AuthData)
	if !ok {
		return nil
	}

	return data
}
