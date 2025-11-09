package user

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"

	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/user/usecase"
)

func (uc *UsecaseImpl) FindOneByProfileID(
	ctx context.Context,
	profileID uint64,
) (*usecase.UserDTO, error) {
	const op = "FindOneByProfileID"

	profile, err := uc.profileUC.FindFullOneByID(ctx, profileID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	account, err := uc.accountUC.FindOneByID(ctx, profile.Profile.AccountID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	role, err := uc.roleUC.GetRoleByID(ctx, account.RoleID)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return &usecase.UserDTO{
		Account:    account,
		Role:       role,
		ProfileDTO: profile,
	}, nil
}

func (uc *UsecaseImpl) FindOneByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
) (*usecase.UserDTO, error) {
	const op = "FindOneByAccountID"

	account, err := uc.accountUC.FindOneByID(ctx, accountID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	profiles, err := uc.profileUC.FindFullList(ctx, &usecase.ProfileListOptions{
		AccountID: &account.ID,
	}, &uctypes.QueryGetListParams{
		Limit: 1,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(profiles) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	role, err := uc.roleUC.GetRoleByID(ctx, account.RoleID)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return &usecase.UserDTO{
		Account:    account,
		Role:       role,
		ProfileDTO: profiles[0],
	}, nil
}

func (uc *UsecaseImpl) FindPagedList(
	ctx context.Context,
	listOptions *usecase.UserListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.UserDTO, uint64, error) {
	const op = "FindList"

	items, total, err := uc.repoProfileAccount.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	filesIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		if item.Profile.Photo100x100FileID != nil {
			filesIDs = append(filesIDs, *item.Profile.Photo100x100FileID)
		}
	}

	filesMap, err := uc.fileUC.FindListInMap(ctx, &fileUC.ListOptions{
		IDs: &filesIDs,
	}, nil)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	rolesIDs := make([]uint64, 0, len(items))
	for _, item := range items {
		rolesIDs = append(rolesIDs, item.Account.RoleID)
	}

	rolesMap := uc.roleUC.GetRolesByIDs(ctx, rolesIDs)

	result := make([]*usecase.UserDTO, 0, len(items))
	for _, item := range items {
		resItem := &usecase.UserDTO{
			Account: item.Account,
			ProfileDTO: &usecase.FullProfileDTO{
				Profile: item.Profile,
			},
		}

		if role, ok := rolesMap[item.Account.RoleID]; ok {
			resItem.Role = role
		}

		if item.Profile.Photo100x100FileID != nil {
			if file, ok := filesMap[*item.Profile.Photo100x100FileID]; ok {
				resItem.ProfileDTO.Photo100x100File = file
			}
		}

		result = append(result, resItem)
	}

	return result, total, nil
}
