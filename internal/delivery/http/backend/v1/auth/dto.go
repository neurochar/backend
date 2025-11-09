package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/dto"

	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

// Output

type OutAccount struct {
	Version int64 `json:"_version,omitempty"`

	ID              uuid.UUID `json:"id"`
	RoleID          uint64    `json:"roleID"`
	Email           string    `json:"email"`
	IsConfirmed     bool      `json:"isConfirmed"`
	IsEmailVerified bool      `json:"isEmailVerified"`

	ProfileName             string       `json:"profileName"`
	ProfileSurname          string       `json:"profileSurname"`
	ProfilePhoto100x100File *dto.FileDTO `json:"profilePhoto100x100File"`
}

type OutLogin struct {
	Account    *OutAccount `json:"account"`
	RefreshJWT string      `json:"refreshJWT"`
	AccessJWT  string      `json:"accessJWT"`
}

// Helpers

func OutAccountDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	accountDTO *tenantUserUC.AccountDTO,
) (*OutAccount, error) {
	out := &OutAccount{
		Version: accountDTO.Account.Version(),

		ID:              accountDTO.Account.ID,
		RoleID:          accountDTO.Account.RoleID,
		Email:           accountDTO.Account.Email,
		IsConfirmed:     accountDTO.Account.IsConfirmed,
		IsEmailVerified: accountDTO.Account.IsEmailVerified,

		ProfileName:             accountDTO.Account.ProfileName,
		ProfileSurname:          accountDTO.Account.ProfileSurname,
		ProfilePhoto100x100File: dto.NewFileDTO(accountDTO.ProfilePhoto100x100File, fc, true),
	}

	return out, nil
}

func OutLoginDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	accountDTO *tenantUserUC.AccountDTO,
	refreshJWT string,
	accessJWT string,
) (*OutLogin, error) {
	outAccout, err := OutAccountDTO(c, fc, accountDTO)
	if err != nil {
		return nil, err
	}

	out := &OutLogin{
		Account:    outAccout,
		AccessJWT:  accessJWT,
		RefreshJWT: refreshJWT,
	}

	return out, nil
}
