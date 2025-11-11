package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
)

type SessionListOptions struct {
	FilterAccountID *uuid.UUID
}

type AuthSessionDTO struct {
	Session       *entity.Session
	RefreshClaims *entity.SessionRefreshClaims
	AccessClaims  *entity.SessionAccessClaims
	AccountDTO    *AccountDTO
}

var ErrPasswordIncorrect = appErrors.ErrUnauthorized.Extend("password incorrect")

var ErrAccountNotConfirmed = appErrors.ErrUnauthorized.Extend("account not confirmed")

var ErrAccountBlocked = appErrors.ErrUnauthorized.Extend("account blocked")

var ErrAccessDenied = appErrors.ErrUnauthorized.Extend("access denied")

var ErrInvalidToken = appErrors.ErrUnauthorized.Extend("invalid token")

var ErrExpiredToken = appErrors.ErrUnauthorized.Extend("expired token").WithTextCode("EXPIRED_TOKEN")

type AuthUsecase interface {
	LoginByEmail(
		ctx context.Context,
		tenantID uuid.UUID,
		email string,
		password string,
		ip net.IP,
	) (res *AuthSessionDTO, resErr error)

	FindSessionByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (res *entity.Session, resErr error)

	IsSessionRevoked(
		ctx context.Context,
		id uuid.UUID,
	) (res bool, resErr error)

	IssueAccessJWT(access *entity.SessionAccessClaims) (resToken string, resErr error)

	IssueRefreshJWT(refresh *entity.SessionRefreshClaims) (resToken string, resErr error)

	ParseAccessToken(token string, validate bool) (res *entity.SessionAccessClaims, resErr error)

	ParseRefreshToken(token string, validate bool) (res *entity.SessionRefreshClaims, resErr error)

	GenerateNewClaims(ctx context.Context, refresh *entity.SessionRefreshClaims, ip net.IP) (res *AuthSessionDTO, resErr error)

	RevokeSessionsByAccountID(ctx context.Context, accountID uuid.UUID) (resErr error)

	RevokeSessionByID(ctx context.Context, ID uuid.UUID) (resErr error)
}

type SessionRepository interface {
	Create(
		ctx context.Context,
		item *entity.Session,
	) (err error)

	Update(
		ctx context.Context,
		item *entity.Session,
	) (err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (session *entity.Session, err error)

	FindList(
		ctx context.Context,
		listOptions *SessionListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Session, err error)
}
