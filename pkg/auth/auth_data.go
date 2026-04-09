package auth

import "github.com/google/uuid"

type AuthDataTokenType uint8

const (
	AuthDataTokenTenantUser AuthDataTokenType = 1
	AuthDataTokenS2S        AuthDataTokenType = 2
)

type AuthData struct {
	tokenType        AuthDataTokenType
	tenantUserClaims *AuthDataTenantUserClaims
	s2sClaims        *AuthDataS2SClaims
}

func (a *AuthData) Type() AuthDataTokenType {
	return a.tokenType
}

func (a *AuthData) IsS2S() bool {
	return a.tokenType == AuthDataTokenS2S
}

func (a *AuthData) IsTenantUser() bool {
	return a.tokenType == AuthDataTokenTenantUser
}

func (a *AuthData) TenantUserClaims() *AuthDataTenantUserClaims {
	return a.tenantUserClaims
}

func (a *AuthData) S2SClaims() *AuthDataS2SClaims {
	return a.s2sClaims
}

type AuthDataTenantUserClaims struct {
	TenantID  uuid.UUID
	SessionID uuid.UUID
	AccountID uuid.UUID
	RoleID    uint64
}

type AuthDataS2SClaims struct {
	ServiceID string
}

func UserTenantClaimsToAuthData(claims *UserTenantSessionAccessClaims) (*AuthData, error) {
	data := &AuthData{
		tokenType: AuthDataTokenTenantUser,
	}

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

	data.tenantUserClaims = &AuthDataTenantUserClaims{
		TenantID:  tenantID,
		SessionID: sessionID,
		AccountID: accountID,
		RoleID:    uint64(claims.RoleId),
	}

	return data, nil
}

func S2SClaimsToAuthData(claims *S2SClaims) (*AuthData, error) {
	data := &AuthData{
		tokenType: AuthDataTokenS2S,
		s2sClaims: &AuthDataS2SClaims{
			ServiceID: claims.ServiceID,
		},
	}

	return data, nil
}
