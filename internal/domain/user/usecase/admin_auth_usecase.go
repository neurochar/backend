package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type AdminAuthUsecase interface {
	LoginByEmail(
		ctx context.Context,
		email string,
		password string,
		ip net.IP,
	) (resSession *userEntity.AdminSession, resAccount *userEntity.Account, resRole *RoleDTO, resErr error)

	SessionToJWT(session *userEntity.AdminSession) (resToken string, resErr error)

	AuthByJWT(
		ctx context.Context,
		token string,
		ip net.IP,
	) (resSession *userEntity.AdminSession, resAccount *userEntity.Account, resRole *RoleDTO, resErr error)

	DeleteActiveSessionsByAccountID(ctx context.Context, accountID uuid.UUID) (resErr error)

	DeleteActiveSessionByID(ctx context.Context, ID uuid.UUID) (resErr error)
}
