package role

import (
	"context"
	"sort"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/samber/lo"

	"github.com/neurochar/backend/internal/domain/user/usecase"
)

func (uc *UsecaseImpl) GetRoleByID(_ context.Context, roleID uint64) (*usecase.RoleDTO, error) {
	const op = "GetRoleByID"

	uc.mu.RLock()
	defer uc.mu.RUnlock()

	role, ok := uc.rolesMap[roleID]
	if !ok {
		return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", uc.pkg, op)
	}

	return role, nil
}

func (uc *UsecaseImpl) GetRolesByIDs(_ context.Context, IDs []uint64) map[uint64]*usecase.RoleDTO {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	IDs = lo.Uniq(IDs)

	result := make(map[uint64]*usecase.RoleDTO, len(IDs))
	for _, roleID := range IDs {
		role, ok := uc.rolesMap[roleID]
		if ok {
			result[roleID] = role
		}
	}

	return result
}

func (uc *UsecaseImpl) GetRolesList(_ context.Context) []*usecase.RoleDTO {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	result := make([]*usecase.RoleDTO, 0, len(uc.rolesMap))
	for _, role := range uc.rolesMap {
		result = append(result, role)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Role.ID < result[j].Role.ID
	})

	return result
}
