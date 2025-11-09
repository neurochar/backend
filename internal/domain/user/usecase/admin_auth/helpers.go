package adminauth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type AuthSessionClaims struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	jwt.RegisteredClaims
}

func (uc *UsecaseImpl) delete(ctx context.Context, session *userEntity.AdminSession) error {
	nowTime := time.Now()
	session.DeletedAt = &nowTime
	return uc.repo.Update(ctx, session)
}

func (uc *UsecaseImpl) generateJWTToken(claims AuthSessionClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(uc.cfg.CPanelApp.Auth.JWTSecret))
}
