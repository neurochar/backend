package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/pkg/auth"
	authpb "github.com/neurochar/backend/pkg/proto_pb/auth"
)

func (uc *UsecaseImpl) makeRefreshClaims(session *entity.Session) *entity.SessionRefreshClaims {
	return &entity.SessionRefreshClaims{
		SessionID:      session.ID,
		RefreshKey:     session.RefreshToken,
		RefreshVersion: session.RefreshVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(session.RefreshTokenIssuedAt),
			ExpiresAt: jwt.NewNumericDate(session.RefreshTokenExpiresAt),
		},
	}
}

func (uc *UsecaseImpl) makeAccessClaims(
	sessionID uuid.UUID,
	tenantID uuid.UUID,
	accountID uuid.UUID,
	roleID uint64,
	meta map[string]string,
	issuedAt time.Time,
	duration time.Duration,
) *auth.SessionAccessClaims {
	return &auth.SessionAccessClaims{
		AccessClaims: authpb.AccessClaims{
			AccountId: accountID.String(),
			SessionId: sessionID.String(),
			TenantId:  tenantID.String(),
			RoleId:    int64(roleID),
			Meta:      meta,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(duration)),
		},
	}
}
