package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
)

func (ctrl *Controller) Auth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	if accessToken == "" {
		return c.Next()
	}

	claims, err := ctrl.authUC.ParseAccessToken(accessToken, true)
	if err != nil {
		if errors.Is(err, tenantUC.ErrInvalidToken) {
			return c.Next()
		}
		return err
	}

	authData, err := auth.ClaimsToAuthData(claims)
	if err != nil {
		return err
	}

	c.Locals(auth.ContextKeyAuthData, authData)

	ctxData, ctxKey := loghandler.SetData(c.Context(), "request.account.id", claims.AccountId)
	c.Locals(ctxKey, ctxData)

	ctxData, ctxKey = loghandler.SetData(c.Context(), "request.tenant.id", claims.TenantId)
	c.Locals(ctxKey, ctxData)

	return c.Next()
}

func GetAuthData(c *fiber.Ctx) *auth.AuthData {
	return auth.GetAuthData(c.Context())
}
