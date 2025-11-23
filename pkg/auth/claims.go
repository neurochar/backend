package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authpb "github.com/neurochar/backend/pkg/proto_pb/auth"
)

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
