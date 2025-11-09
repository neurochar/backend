package usecase

import (
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type User struct {
	Account *userEntity.Account
	Profile *userEntity.Profile
}

type UserDTO struct {
	Account    *userEntity.Account
	Role       *RoleDTO
	ProfileDTO *FullProfileDTO
}

type UserListOptions struct {
	Query  *string
	RoleID *uint64
}

type CreateUserInput struct {
	Account AccountDataInput
	Profile ProfileDataInput

	IsSendPassword bool
}
