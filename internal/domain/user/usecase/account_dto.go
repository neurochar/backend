package usecase

import (
	"github.com/google/uuid"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type AccountListOptions struct {
	IDs    *[]uuid.UUID
	RoleID *uint64
}

type AccountDataInput struct {
	Email           string
	Password        string
	RoleID          uint64
	IsEmailVerified bool
}

type PatchAccountDataInput struct {
	Version int64

	Email           *string
	Password        *string
	RoleID          *uint64
	IsEmailVerified *bool
	IsBlocked       *bool
}

type AccountCodeListOptions struct {
	AccountID *uuid.UUID
	Type      *userEntity.AccountCodeType
	IsActive  *bool
}
