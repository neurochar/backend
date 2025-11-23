package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/dto"

	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
)

// Output

type OutAccount struct {
	Version int64 `json:"_version,omitempty"`

	ID              uuid.UUID `json:"id"`
	RoleID          uint64    `json:"roleID"`
	Email           string    `json:"email"`
	IsConfirmed     bool      `json:"isConfirmed"`
	IsEmailVerified bool      `json:"isEmailVerified"`

	ProfileName              string       `json:"profileName"`
	ProfileSurname           string       `json:"profileSurname"`
	ProfilePhoto100x100File  *dto.FileDTO `json:"profilePhoto100x100File"`
	ProfilePhotoOriginalFile *dto.FileDTO `json:"profilePhotoOriginalFile"`
}

type OutTenant struct {
	Version int64 `json:"_version,omitempty"`

	ID     uuid.UUID `json:"id"`
	TextID string    `json:"textID"`
	Name   string    `json:"name"`
}

type OutTokens struct {
	RefreshJWT     string `json:"refreshJWT"`
	RefreshLifeSec uint64 `json:"refreshLifeSec"`
	AccessJWT      string `json:"accessJWT"`
}

type OutLogin struct {
	Account *OutAccount `json:"account"`
	Tenant  *OutTenant  `json:"tenant"`
	OutTokens
}

type OutWhoIAm struct {
	Account *OutAccount `json:"account"`
	Tenant  *OutTenant  `json:"tenant"`
}

// Helpers

func OutAccountDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	accountDTO *tenantUC.AccountDTO,
) (*OutAccount, error) {
	out := &OutAccount{
		Version: accountDTO.Account.Version(),

		ID:              accountDTO.Account.ID,
		RoleID:          accountDTO.Account.RoleID,
		Email:           accountDTO.Account.Email,
		IsConfirmed:     accountDTO.Account.IsConfirmed,
		IsEmailVerified: accountDTO.Account.IsEmailVerified,

		ProfileName:              accountDTO.Account.ProfileName,
		ProfileSurname:           accountDTO.Account.ProfileSurname,
		ProfilePhoto100x100File:  dto.NewFileDTO(accountDTO.ProfilePhoto100x100File, fc, true),
		ProfilePhotoOriginalFile: dto.NewFileDTO(accountDTO.ProfilePhotoOriginalFile, fc, true),
	}

	return out, nil
}

func OutTenantDTOByAccountDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	accountDTO *tenantUC.AccountDTO,
) (*OutTenant, error) {
	out := &OutTenant{
		Version: accountDTO.Account.Version(),

		ID:     accountDTO.Tenant.ID,
		TextID: accountDTO.Tenant.TextID,
		Name:   accountDTO.Tenant.Name,
	}

	return out, nil
}

func OutTokensDTO(
	c *fiber.Ctx,
	refreshJWT string,
	refreshLifeSec uint64,
	accessJWT string,
) (*OutTokens, error) {
	out := &OutTokens{
		AccessJWT:      accessJWT,
		RefreshLifeSec: refreshLifeSec,
		RefreshJWT:     refreshJWT,
	}

	return out, nil
}

func OutLoginDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	accountDTO *tenantUC.AccountDTO,
	refreshJWT string,
	refreshLifeSec uint64,
	accessJWT string,
) (*OutLogin, error) {
	outAccount, err := OutAccountDTO(c, fc, accountDTO)
	if err != nil {
		return nil, err
	}

	outTenant, err := OutTenantDTOByAccountDTO(c, fc, accountDTO)
	if err != nil {
		return nil, err
	}

	out := &OutLogin{
		Account: outAccount,
		Tenant:  outTenant,
		OutTokens: OutTokens{
			AccessJWT:      accessJWT,
			RefreshLifeSec: refreshLifeSec,
			RefreshJWT:     refreshJWT,
		},
	}

	return out, nil
}

func OutWhoIAmDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	accountDTO *tenantUC.AccountDTO,
) (*OutWhoIAm, error) {
	outAccount, err := OutAccountDTO(c, fc, accountDTO)
	if err != nil {
		return nil, err
	}

	outTenant, err := OutTenantDTOByAccountDTO(c, fc, accountDTO)
	if err != nil {
		return nil, err
	}

	out := &OutWhoIAm{
		Account: outAccount,
		Tenant:  outTenant,
	}

	return out, nil
}
