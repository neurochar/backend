package usecase

import (
	"context"
)

type RoleUsecase interface {
	BuildRolesInMemory(ctx context.Context) error

	GetRoleByID(ctx context.Context, roleID uint64) (*RoleDTO, error)

	GetRolesByIDs(ctx context.Context, IDs []uint64) map[uint64]*RoleDTO

	GetRolesList(ctx context.Context) []*RoleDTO

	CreateRole(ctx context.Context, in CreateRoleInput) (*RoleDTO, error)

	UpdateRole(ctx context.Context, roleID uint64, in UpdateRoleInput, skipVersionCheck bool) (*RoleDTO, error)

	DeleteRole(ctx context.Context, roleID uint64) error
}
