package auth

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
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/pkg/emailnormalize"
)

func (uc *UsecaseImpl) LoginByEmail(
	ctx context.Context,
	tenantID uuid.UUID,
	email string,
	password string,
	ip net.IP,
) (*usecase.AuthSessionDTO, error) {
	const op = "LoginByEmail"

	res, err := emailnormalize.Normalize(strings.TrimSpace(email))
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", uc.pkg, op)
	}

	email = res.NormalizedAddress

	out := &usecase.AuthSessionDTO{}

	timeNow := time.Now()

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		out.AccountDTO, err = uc.accountUC.FindOneByEmail(ctx, tenantID, email, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		}, &usecase.AccountDTOOptions{
			FetchTenant:     true,
			FetchPhotoFiles: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if !out.AccountDTO.Account.VerifyPassword(password) {
			return usecase.ErrPasswordIncorrect
		}

		if !out.AccountDTO.Account.IsConfirmed {
			return usecase.ErrAccountNotConfirmed
		}

		if out.AccountDTO.Account.IsBlocked {
			return usecase.ErrAccountBlocked
		}

		out.Session = entity.NewSession(
			out.AccountDTO.Account.ID,
			ip,
			timeNow,
			time.Duration(uc.cfg.Auth.RefreshTokenLifetimeHrs)*time.Hour,
		)
		if err := uc.repo.Create(ctx, out.Session); err != nil {
			return err
		}

		out.AccountDTO.Account.SetLastRequestAt(&timeNow)
		out.AccountDTO.Account.SetLastLoginAt(&timeNow)
		out.AccountDTO.Account.SetLastRequestIP(&ip)

		err = uc.accountUC.UpdateAccount(ctx, out.AccountDTO.Account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out.RefreshClaims = uc.makeRefreshClaims(out.Session)
	out.AccessClaims = uc.makeAccessClaims(
		out.Session.ID,
		out.AccountDTO.Account.TenantID,
		out.AccountDTO.Account.ID,
		out.AccountDTO.Account.RoleID,
		nil,
		timeNow,
		time.Duration(uc.cfg.Auth.AccessTokenLifetimeSec)*time.Second,
	)

	return out, nil
}

func (uc *UsecaseImpl) IssueAccessJWT(access *entity.SessionAccessClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, access)

	return token.SignedString([]byte(uc.cfg.Auth.JwtAccessSecret))
}

func (uc *UsecaseImpl) IssueRefreshJWT(refresh *entity.SessionRefreshClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh)

	return token.SignedString([]byte(uc.cfg.Auth.JwtRefreshSecret))
}

func (uc *UsecaseImpl) ParseAccessToken(tokenStr string, validate bool) (*entity.SessionAccessClaims, error) {
	const op = "ParseAccessToken"

	ops := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	}

	if !validate {
		ops = append(ops, jwt.WithoutClaimsValidation())
	}

	token, err := jwt.ParseWithClaims(tokenStr, &entity.SessionAccessClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(uc.cfg.Auth.JwtAccessSecret), nil
	}, ops...)
	if err != nil || !token.Valid {
		return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	claims, ok := token.Claims.(*entity.SessionAccessClaims)
	if !ok {
		return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	return claims, nil
}

func (uc *UsecaseImpl) ParseRefreshToken(tokenStr string, validate bool) (*entity.SessionRefreshClaims, error) {
	const op = "ParseRefreshToken"

	ops := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	}

	if !validate {
		ops = append(ops, jwt.WithoutClaimsValidation())
	}

	token, err := jwt.ParseWithClaims(tokenStr, &entity.SessionRefreshClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(uc.cfg.Auth.JwtRefreshSecret), nil
	}, ops...)
	if err != nil || !token.Valid {
		return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	claims, ok := token.Claims.(*entity.SessionRefreshClaims)
	if !ok {
		return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	return claims, nil
}

func (uc *UsecaseImpl) GenerateNewClaims(
	ctx context.Context,
	refreshClaims *entity.SessionRefreshClaims,
	ip net.IP,
) (*usecase.AuthSessionDTO, error) {
	const op = "GenerateNewClaims"

	out := &usecase.AuthSessionDTO{}

	timeNow := time.Now()

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		out.Session, err = uc.repo.FindOneByID(ctx, refreshClaims.SessionID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if out.Session.RefreshToken != refreshClaims.RefreshKey {
			return appErrors.ErrUnauthorized
		}

		if out.Session.RefreshVersion != refreshClaims.RefreshVersion {
			return appErrors.ErrUnauthorized
		}

		if out.Session.RefreshTokenExpiresAt.Before(timeNow) {
			return appErrors.ErrUnauthorized
		}

		out.AccountDTO, err = uc.accountUC.FindOneByID(ctx, out.Session.AccountID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		}, &usecase.AccountDTOOptions{
			FetchTenant: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if !out.AccountDTO.Account.IsConfirmed {
			return usecase.ErrAccountNotConfirmed
		}

		if out.AccountDTO.Account.IsBlocked {
			return usecase.ErrAccountBlocked
		}

		out.Session.GenerateNewRefresh(timeNow, time.Duration(uc.cfg.Auth.RefreshTokenLifetimeHrs)*time.Hour, ip)

		err = uc.repo.Update(ctx, out.Session)
		if err != nil {
			return err
		}

		out.AccountDTO.Account.SetLastRequestAt(&timeNow)
		out.AccountDTO.Account.SetLastLoginAt(&timeNow)
		out.AccountDTO.Account.SetLastRequestIP(&ip)

		err = uc.accountUC.UpdateAccount(ctx, out.AccountDTO.Account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out.RefreshClaims = uc.makeRefreshClaims(out.Session)
	out.AccessClaims = uc.makeAccessClaims(
		out.Session.ID,
		out.AccountDTO.Account.TenantID,
		out.AccountDTO.Account.ID,
		out.AccountDTO.Account.RoleID,
		nil,
		timeNow,
		time.Duration(uc.cfg.Auth.AccessTokenLifetimeSec)*time.Second,
	)

	return out, nil
}

func (uc *UsecaseImpl) RevokeSessionsByAccountID(ctx context.Context, accountID uuid.UUID) error {
	const op = "RevokeSessionsByAccountID"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		sessions, err := uc.repo.FindList(ctx, &usecase.SessionListOptions{
			FilterAccountID: &accountID,
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

func (uc *UsecaseImpl) RevokeSessionByID(ctx context.Context, ID uuid.UUID) error {
	const op = "RevokeSessionByID"

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
