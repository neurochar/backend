package users

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/dto"
	"github.com/samber/lo"

	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

// Output

type OutAccountRoleRight struct {
	ID    uint64 `json:"id"`
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type OutAccountRole struct {
	Version int64 `json:"_version,omitempty"`

	ID       uint64                         `json:"id"`
	Name     string                         `json:"name"`
	IsSystem bool                           `json:"isSystem"`
	IsSuper  bool                           `json:"isSuper"`
	Rights   map[string]OutAccountRoleRight `json:"rights,omitempty"`
}

type OutAccount struct {
	Version int64 `json:"_version,omitempty"`

	ID              uuid.UUID      `json:"id"`
	Email           string         `json:"email"`
	IsEmailVerified bool           `json:"isEmailVerified"`
	IsBlocked       bool           `json:"isBlocked"`
	LastLoginAt     *time.Time     `json:"lastLoginAt"`
	LastRequestAt   *time.Time     `json:"lastRequestAt"`
	LastRequestIP   *string        `json:"lastRequestIP"`
	Role            OutAccountRole `json:"role"`
}

type OutProfile struct {
	Version int64 `json:"_version,omitempty"`

	ID               uint64       `json:"id"`
	Name             string       `json:"name"`
	Surname          string       `json:"surname"`
	Photo100x100File *dto.FileDTO `json:"photo100x100File"`
}

type OutUser struct {
	Account OutAccount `json:"account"`
	Profile OutProfile `json:"profile"`
}

// Input

type InAccountCreate struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"min=0"`
	RoleID          uint64 `json:"roleID" validate:"required"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	IsSendPassword  bool   `json:"isSendPassword"`
}

type InAccountRole struct {
	Name   string         `json:"name" validate:"required,min=1,max=150"`
	Rights map[string]int `json:"rights"`
}

type InProfile struct {
	Name               string `json:"name" validate:"required,min=1,max=150"`
	Surname            string `json:"surname" validate:"required,min=1,max=150"`
	Photo100x100FileID string `json:"photo100x100FileID" validate:"omitempty,uuid"`
}

// Helpers

func OutUserDTO(
	c *fiber.Ctx,
	fc fileUC.Usecase,
	account *userEntity.Account,
	profileDTO *userUC.FullProfileDTO,
	roleDTO *userUC.RoleDTO,
	echoRightsMap bool,
) (*OutUser, error) {
	var LastRequestIP *string
	if account.LastRequestIP != nil {
		LastRequestIP = lo.ToPtr(account.LastRequestIP.String())
	}

	out := &OutUser{
		Account: OutAccount{
			Version: account.Version(),

			ID:              account.ID,
			Email:           account.Email,
			IsEmailVerified: account.IsEmailVerified,
			IsBlocked:       account.IsBlocked,
			LastLoginAt:     account.LastLoginAt,
			LastRequestAt:   account.LastRequestAt,
			LastRequestIP:   LastRequestIP,
			Role: OutAccountRole{
				Version:  roleDTO.Role.Version(),
				ID:       roleDTO.Role.ID,
				Name:     roleDTO.Role.Name,
				IsSuper:  roleDTO.Role.IsSuper,
				IsSystem: roleDTO.Role.IsSystem,
			},
		},
		Profile: OutProfile{
			Version:          profileDTO.Profile.Version(),
			ID:               profileDTO.Profile.ID,
			Name:             profileDTO.Profile.Name,
			Surname:          profileDTO.Profile.Surname,
			Photo100x100File: dto.NewFileDTO(profileDTO.Photo100x100File, fc, true),
		},
	}

	if echoRightsMap {
		out.Account.Role.Rights = make(map[string]OutAccountRoleRight, len(roleDTO.Rights))

		if roleDTO.Rights != nil {
			for _, right := range roleDTO.Rights {
				out.Account.Role.Rights[right.Right.Key] = OutAccountRoleRight{
					ID:    right.Right.ID,
					Key:   right.Right.Key,
					Type:  right.Right.Type.String(),
					Value: right.Value,
				}
			}
		}
	}

	return out, nil
}
