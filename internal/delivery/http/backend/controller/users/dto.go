package users

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/dto"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
)

type OutAccount struct {
	Version int64 `json:"_version,omitempty"`

	ID                       uuid.UUID    `json:"id"`
	TenantID                 uuid.UUID    `json:"tenantID"`
	RoleID                   uint64       `json:"roleID"`
	Email                    string       `json:"email"`
	IsConfirmed              bool         `json:"isConfirmed"`
	IsEmailVerified          bool         `json:"isEmailVerified"`
	IsBlocked                bool         `json:"isBlocked"`
	LastLoginAt              *time.Time   `json:"lastLoginAt"`
	LastRequestAt            *time.Time   `json:"lastRequestAt"`
	ProfileName              string       `json:"profileName"`
	ProfileSurname           string       `json:"profileSurname"`
	ProfilePhotoOriginalFile *dto.FileDTO `json:"profilePhotoOriginalFile"`
	ProfilePhoto100x100File  *dto.FileDTO `json:"profilePhoto100x100File"`
}

func OutAccountDTO(
	c *fiber.Ctx,
	fullDTO bool,
	fc fileUC.Usecase,
	accountDTO *tenantUC.AccountDTO,
) (*OutAccount, error) {
	_ = fullDTO

	out := &OutAccount{
		Version:         accountDTO.Account.Version(),
		ID:              accountDTO.Account.ID,
		TenantID:        accountDTO.Account.TenantID,
		RoleID:          accountDTO.Account.RoleID,
		Email:           accountDTO.Account.Email,
		IsConfirmed:     accountDTO.Account.IsConfirmed,
		IsEmailVerified: accountDTO.Account.IsEmailVerified,
		IsBlocked:       accountDTO.Account.IsBlocked,
		LastLoginAt:     accountDTO.Account.LastLoginAt,
		LastRequestAt:   accountDTO.Account.LastRequestAt,

		ProfileName:              accountDTO.Account.ProfileName,
		ProfileSurname:           accountDTO.Account.ProfileSurname,
		ProfilePhotoOriginalFile: dto.NewFileDTO(accountDTO.ProfilePhotoOriginalFile, fc, true),
		ProfilePhoto100x100File:  dto.NewFileDTO(accountDTO.ProfilePhoto100x100File, fc, true),
	}

	return out, nil
}
