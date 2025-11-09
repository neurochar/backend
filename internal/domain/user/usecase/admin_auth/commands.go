package adminauth

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	userConstants "github.com/neurochar/backend/internal/domain/user/constants"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/pkg/emailnormalize"
)

func (uc *UsecaseImpl) LoginByEmail(
	ctx context.Context,
	email string,
	password string,
	ip net.IP,
) (*userEntity.AdminSession, *userEntity.Account, *usecase.RoleDTO, error) {
	const op = "LoginByEmail"

	res, err := emailnormalize.Normalize(strings.TrimSpace(email))
	if err != nil {
		return nil, nil, nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", uc.pkg, op)
	}

	email = res.NormalizedAddress

	var account *userEntity.Account
	var role *usecase.RoleDTO
	var session *userEntity.AdminSession

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		account, err = uc.accountUC.FindOneByEmail(ctx, email, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if !account.VerifyPassword(password) {
			return usecase.ErrPasswordIncorrect
		}

		if !account.IsConfirmed() {
			return usecase.ErrAccountNotConfirmed
		}

		if account.IsBlocked {
			return usecase.ErrAccountBlocked
		}

		role, err = uc.roleUC.GetRoleByID(ctx, account.RoleID)
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		accessRight, err := role.GetRightByKey(userConstants.RightKeyAccessToAdminPanel)
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		if accessRight.Value != 1 {
			return usecase.ErrAccessDenied
		}

		session = userEntity.NewSession(account.ID, ip)
		if err := uc.repo.Create(ctx, session); err != nil {
			return err
		}

		timeNow := time.Now()

		account.SetLastRequestAt(&timeNow)
		account.SetLastLoginAt(&timeNow)
		account.SetLastRequestIP(&ip)

		err = uc.accountUC.UpdateAccount(ctx, account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return session, account, role, nil
}

func (uc *UsecaseImpl) SessionToJWT(session *userEntity.AdminSession) (string, error) {
	const op = "SessionToJWT"

	claims := AuthSessionClaims{
		ID:        session.ID,
		AccountID: session.AccountID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := uc.generateJWTToken(claims)
	if err != nil {
		return "", appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return token, nil
}

func (uc *UsecaseImpl) AuthByJWT(
	ctx context.Context,
	tokenStr string,
	ip net.IP,
) (*userEntity.AdminSession, *userEntity.Account, *usecase.RoleDTO, error) {
	const op = "AuthByJWT"

	token, err := jwt.ParseWithClaims(tokenStr, &AuthSessionClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(uc.cfg.CPanelApp.Auth.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, nil, nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	claims, ok := token.Claims.(*AuthSessionClaims)
	if !ok {
		return nil, nil, nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	var account *userEntity.Account
	var role *usecase.RoleDTO
	var session *userEntity.AdminSession

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		session, err = uc.repo.FindOneByID(ctx, claims.ID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return usecase.ErrInvalidToken
			}
			return err
		}

		if session.AccountID != claims.AccountID {
			return usecase.ErrInvalidToken
		}

		if !session.IsAlive(time.Duration(uc.cfg.CPanelApp.Auth.LifeTimeWithoutActivitySec) * time.Second) {
			err := uc.delete(ctx, session)
			if err != nil {
				return err
			}

			return usecase.ErrExpiredToken
		}

		account, err = uc.accountUC.FindOneByID(ctx, session.AccountID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return usecase.ErrInvalidToken
			}
			return err
		}

		if !account.IsConfirmed() {
			return usecase.ErrAccountNotConfirmed
		}

		if account.IsBlocked {
			return usecase.ErrAccountBlocked
		}

		role, err = uc.roleUC.GetRoleByID(ctx, account.RoleID)
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		accessRight, err := role.GetRightByKey(userConstants.RightKeyAccessToAdminPanel)
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		if accessRight.Value != 1 {
			return usecase.ErrAccessDenied
		}

		timeNow := time.Now()

		session.SetLastRequestAt(timeNow)
		session.SetLastRequestIP(ip)

		err = uc.repo.Update(ctx, session)
		if err != nil {
			return err
		}

		account.SetLastRequestAt(&timeNow)
		account.SetLastRequestIP(&ip)

		err = uc.accountUC.UpdateAccount(ctx, account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return session, account, role, nil
}

func (uc *UsecaseImpl) DeleteActiveSessionsByAccountID(ctx context.Context, accountID uuid.UUID) error {
	const op = "DeleteActiveSessionsByAccountID"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		sessions, err := uc.repo.FindList(ctx, &usecase.AdminAuthListOptions{
			AccountID: &accountID,
		}, &uctypes.QueryGetListParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		for _, session := range sessions {
			err := uc.delete(ctx, session)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) DeleteActiveSessionByID(ctx context.Context, ID uuid.UUID) error {
	const op = "DeleteActiveSessionByID"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		session, err := uc.repo.FindOneByID(ctx, ID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		err = uc.delete(ctx, session)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
