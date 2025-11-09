package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authpb "github.com/neurochar/backend/pkg/proto_pb/auth"
)

type SessionRefreshClaims struct {
	SessionID      uuid.UUID `json:"sid,omitempty"`
	RefreshKey     uuid.UUID `json:"rk,omitempty"`
	RefreshVersion uint64    `json:"rv,omitempty"`
	jwt.RegisteredClaims
}

type SessionAccessClaims struct {
	authpb.AccessClaims
	jwt.RegisteredClaims
}

func (c *SessionAccessClaims) GetTenantID() (uuid.UUID, error) {
	return uuid.Parse(c.TenantId)
}

func (c *SessionAccessClaims) GetSessionID() (uuid.UUID, error) {
	return uuid.Parse(c.SessionId)
}

func (c *SessionAccessClaims) GetAccountID() (uuid.UUID, error) {
	return uuid.Parse(c.AccountId)
}
