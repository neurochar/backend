package auth

import "github.com/google/uuid"

type AuthData struct {
	TenantID  uuid.UUID
	SessionID uuid.UUID
	AccountID uuid.UUID
	RoleID    uint64
}

func ClaimsToAuthData(claims *SessionAccessClaims) (*AuthData, error) {
	tenantID, err := claims.GetTenantID()
	if err != nil {
		return nil, err
	}

	sessionID, err := claims.GetSessionID()
	if err != nil {
		return nil, err
	}

	accountID, err := claims.GetAccountID()
	if err != nil {
		return nil, err
	}

	data := &AuthData{
		TenantID:  tenantID,
		SessionID: sessionID,
		AccountID: accountID,
		RoleID:    uint64(claims.RoleId),
	}

	return data, nil
}
