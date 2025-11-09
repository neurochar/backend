package role

import (
	"context"
	"fmt"

	"github.com/neurochar/backend/internal/common/uctypes"
	userConstants "github.com/neurochar/backend/internal/domain/user/constants"

	appErrors "github.com/neurochar/backend/internal/app/errors"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"

	"github.com/neurochar/backend/internal/domain/user/usecase"
)

func (uc *UsecaseImpl) BuildRolesInMemory(ctx context.Context) error {
	const op = "BuildRolesInMemory"

	uc.mu.Lock()
	defer uc.mu.Unlock()

	roles, err := uc.repoRole.FindList(ctx, nil, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	rolesToRights, err := uc.repoRoleToRight.FindList(ctx, nil, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	rolesToRightsMap := make(map[uint64]map[uint64]int, len(roles))

	for _, roleToRight := range rolesToRights {
		if _, ok := rolesToRightsMap[roleToRight.RoleID]; !ok {
			rolesToRightsMap[roleToRight.RoleID] = make(map[uint64]int, len(userConstants.Rights))
		}

		rolesToRightsMap[roleToRight.RoleID][roleToRight.RightID] = roleToRight.Value
	}

	uc.rolesMap = make(map[uint64]*usecase.RoleDTO, len(roles))
	for _, role := range roles {
		roleMap := &usecase.RoleDTO{
			Role:   role,
			Rights: make(map[uint64]*usecase.RoleRightDTO, len(userConstants.Rights)),
		}

		for _, right := range userConstants.Rights {
			roleRight := &usecase.RoleRightDTO{
				Right: right,
			}

			if val, ok := rolesToRightsMap[role.ID][right.ID]; ok {
				roleRight.Value = val
			} else if role.IsSuper {
				roleRight.Value = right.DefaultSuperValue
			} else {
				roleRight.Value = right.DefaultValue
			}

			roleMap.Rights[right.ID] = roleRight
		}

		uc.rolesMap[role.ID] = roleMap
	}

	return nil
}

func (uc *UsecaseImpl) CreateRole(ctx context.Context, in usecase.CreateRoleInput) (*usecase.RoleDTO, error) {
	const op = "CreateRole"

	if in.Rights == nil {
		in.Rights = make(map[string]int)
	}

	role, err := userEntity.NewRole(in.Name)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		err := uc.repoRole.Create(ctx, role)
		if err != nil {
			return err
		}

		for _, right := range userConstants.Rights {
			var value int

			if inVal, ok := in.Rights[right.Key]; ok {
				value = inVal
			} else {
				return appErrors.ErrBadRequest.WithHints(fmt.Sprintf("right %s not found", right.Key))
			}

			if right.Type == userEntity.RightTypeBool {
				if value < 0 {
					value = 0
				} else if value > 1 {
					value = 1
				}
			}

			roleToRight := userEntity.NewRoleToRight(role.ID, right.ID, value)

			err := uc.repoRoleToRight.Create(ctx, roleToRight)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.BuildRolesInMemory(ctx)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	roleRes, err := uc.GetRoleByID(ctx, role.ID)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return roleRes, nil
}

func (uc *UsecaseImpl) UpdateRole(
	ctx context.Context,
	roleID uint64,
	in usecase.UpdateRoleInput,
	skipVersionCheck bool,
) (*usecase.RoleDTO, error) {
	const op = "UpdateRole"

	if in.Rights == nil {
		in.Rights = make(map[string]int)
	}

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		role, err := uc.repoRole.FindOneByID(ctx, roleID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if !skipVersionCheck && role.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, role.Version()).
				WithDetail("last_updated_at", false, role.UpdatedAt)
		}

		role.Name = in.Name

		err = uc.repoRole.Update(ctx, role)
		if err != nil {
			return err
		}

		if !role.IsSystem {
			err = uc.repoRoleToRight.DeleteByRoleID(ctx, roleID)
			if err != nil {
				return err
			}

			for _, right := range userConstants.Rights {
				var value int

				if inVal, ok := in.Rights[right.Key]; ok {
					value = inVal
				} else {
					return appErrors.ErrBadRequest.WithHints(fmt.Sprintf("right %s not found", right.Key))
				}

				if right.Type == userEntity.RightTypeBool {
					if value < 0 {
						value = 0
					} else if value > 1 {
						value = 1
					}
				}

				roleToRight := userEntity.NewRoleToRight(role.ID, right.ID, value)

				err := uc.repoRoleToRight.Create(ctx, roleToRight)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.BuildRolesInMemory(ctx)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	roleRes, err := uc.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return roleRes, nil
}

func (uc *UsecaseImpl) DeleteRole(ctx context.Context, roleID uint64) error {
	const op = "DeleteRole"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		role, err := uc.repoRole.FindOneByID(ctx, roleID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if role.IsSystem {
			return usecase.ErrCantDeleteRoleIsSystem
		}

		err = uc.repoRoleToRight.DeleteByRoleID(ctx, roleID)
		if err != nil {
			return err
		}

		err = uc.deleteRole(ctx, role)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.BuildRolesInMemory(ctx)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
