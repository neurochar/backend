package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthData struct {
	TenantID  uuid.UUID
	SessionID uuid.UUID
	AccountID uuid.UUID
	RoleID    uint64
}

const AuthDataKey = "auth_data"

func (ctrl *Controller) MiddlewareAuth(c *fiber.Ctx) error {
	// skipErr := false

	// cookieValue := c.Cookies("auth_admin_session")

	// requestIP := net.ParseIP(middleware.GetRealIP(c))

	// session, account, role, err := ctrl.authAdminUC.AuthByJWT(c.Context(), cookieValue, requestIP)
	// if err != nil {
	// 	if skipErr {
	// 		return c.Next()
	// 	}
	// 	return err
	// }

	// c.Locals(AuthDataKey, &AuthData{
	// 	Account: account,
	// 	Session: session,
	// 	Role:    role,
	// })

	// ctxData, ctxKey := loghandler.SetData(c.Context(), "request.account.id", account.ID)
	// c.Locals(ctxKey, ctxData)

	return c.Next()
}

func GetAuthData(c *fiber.Ctx) *AuthData {
	data, ok := c.Locals(AuthDataKey).(*AuthData)
	if !ok {
		return nil
	}

	return data
}
