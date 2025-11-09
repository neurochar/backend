package account

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/app/config"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

func (uc *UsecaseImpl) createAccountEmailVerificationCode(
	ctx context.Context,
	account *userEntity.Account,
	requestIP net.IP,
) (*userEntity.AccountCode, error) {
	const op = "createAccountEmailVerificationCode"

	code, err := userEntity.NewAccountCode(account.ID, userEntity.AccountCodeTypeEmailVerification, requestIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repoAccountCode.Create(ctx, code)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

func (uc *UsecaseImpl) createAccountPasswordRecoveryCode(
	ctx context.Context,
	account *userEntity.Account,
	requestIP net.IP,
) (*userEntity.AccountCode, error) {
	const op = "createAccountPasswordRecoveryCode"

	code, err := userEntity.NewAccountCode(account.ID, userEntity.AccountCodeTypePasswordRecovery, requestIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repoAccountCode.Create(ctx, code)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

type RecoveryCodeEmailData struct {
	Cfg       config.Config
	Account   *userEntity.Account
	Code      string
	RequestIP string
}

func (uc *UsecaseImpl) sendRecoveryCodeEmailToUser(
	ctx context.Context,
	account *userEntity.Account,
	recoveryCode *userEntity.AccountCode,
	requestIP net.IP,
) error {
	const op = "sendRecoveryCodeEmailToUser"

	data := RecoveryCodeEmailData{
		Cfg:       uc.cfg,
		Account:   account,
		Code:      recoveryCode.Code,
		RequestIP: requestIP.String(),
	}

	msg, err := emailing.NewMessageFromJetTpl("users/email_password_recovery_code.jet", data)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.emailing.SendAsyc(emailing.Message{
		To:            account.Email,
		Subject:       fmt.Sprintf("Восстановление пароля на сайте %s", uc.cfg.Global.ProjectName),
		TextHtml:      msg,
		AutoTextPlain: true,
	}, func(sendErr error) {
		if sendErr != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "email sending", slog.Any("error", sendErr))
		}
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

// Максимальное количество попыток проверки кода
const maxCodeAttempts = 3

// Максимальное время жизни кода
const maxCodeLifetime = time.Hour * 24 * 7

func (uc *UsecaseImpl) checkCodeByID(ctx context.Context, codeID uuid.UUID, codeValue string) (*userEntity.AccountCode, error) {
	const op = "checkCodeByID"

	var codeItem *userEntity.AccountCode
	var resErr error

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		codeItem, err = uc.repoAccountCode.FindOneByID(ctx, codeID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if !codeItem.IsActive {
			resErr = appErrors.ErrBadRequest

			return nil
		}

		if !codeItem.VerifyCode(codeValue) {
			codeItem.AddAttempt()

			if codeItem.Attempts >= maxCodeAttempts {
				codeItem.Deactivate()
			}

			maxAttempts := maxCodeAttempts - codeItem.Attempts
			if maxAttempts < 0 {
				maxAttempts = 0
			}

			err := uc.repoAccountCode.Update(ctx, codeItem)
			if err != nil {
				return err
			}

			resErr = usecase.ErrCodeInvalid.WithDetail("max_attempts", false, maxAttempts)

			return nil
		}

		if !codeItem.IsAlive(maxCodeLifetime) {
			codeItem.Deactivate()

			err := uc.repoAccountCode.Update(ctx, codeItem)
			if err != nil {
				return err
			}

			resErr = usecase.ErrCodeExpired

			return nil
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return codeItem, resErr
}
