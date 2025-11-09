package usecase

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type RoleListOptions struct{}

type RoleRightDTO struct {
	Right *userEntity.Right
	Value int
}

type RoleDTO struct {
	Role   *userEntity.Role
	Rights map[uint64]*RoleRightDTO
}

func (r *RoleDTO) GetRightByID(rightID uint64) (*RoleRightDTO, error) {
	item, ok := r.Rights[rightID]
	if !ok {
		return nil, appErrors.ErrNotFound
	}

	return item, nil
}

func (r *RoleDTO) GetRightByKey(key string) (*RoleRightDTO, error) {
	for _, item := range r.Rights {
		if item.Right.Key == key {
			return item, nil
		}
	}

	return nil, appErrors.ErrNotFound
}

type CreateRoleInput struct {
	Name   string
	Rights map[string]int
}

type UpdateRoleInput struct {
	Version int64

	Name   string
	Rights map[string]int
}

type RoleToRightListOptions struct{}
