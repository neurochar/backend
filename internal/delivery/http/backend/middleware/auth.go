package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

type AuthData struct {
	TenantID  uuid.UUID
	SessionID uuid.UUID
	AccountID uuid.UUID
	RoleID    uint64
}

const AuthDataKey = "auth_data"

func (ctrl *Controller) MiddlewareAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	if accessToken == "" {
		return c.Next()
	}

	claims, err := ctrl.authUC.ParseAccessToken(accessToken, true)
	if err != nil {
		if errors.Is(err, tenantUserUC.ErrInvalidToken) {
			return appErrors.ErrUnauthorized
		}
		return err
	}

	tenantID, err := uuid.Parse(claims.TenantId)
	if err != nil {
		return err
	}

	sessionID, err := uuid.Parse(claims.SessionId)
	if err != nil {
		return err
	}

	accountID, err := uuid.Parse(claims.AccountId)
	if err != nil {
		return err
	}

	c.Locals(AuthDataKey, &AuthData{
		TenantID:  tenantID,
		SessionID: sessionID,
		AccountID: accountID,
		RoleID:    uint64(claims.RoleId),
	})

	ctxData, ctxKey := loghandler.SetData(c.Context(), "request.account.id", claims.AccountId)
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
