package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authpb "github.com/neurochar/backend/pkg/proto_pb/common/auth_tenant"
)

type UserTenantSessionAccessClaims struct {
	authpb.AccessClaims
	jwt.RegisteredClaims
}

func (c *UserTenantSessionAccessClaims) GetTenantID() (uuid.UUID, error) {
	return uuid.Parse(c.TenantId)
}

func (c *UserTenantSessionAccessClaims) GetSessionID() (uuid.UUID, error) {
	return uuid.Parse(c.SessionId)
}

func (c *UserTenantSessionAccessClaims) GetAccountID() (uuid.UUID, error) {
	return uuid.Parse(c.AccountId)
}

type S2SClaims struct {
	ServiceID string
}
