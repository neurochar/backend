package usecase

import (
	"context"

	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type ProfileUsecase interface {
	FindList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*userEntity.Profile, resErr error)

	FindFullList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*FullProfileDTO, resErr error)

	FindFullOneByID(
		ctx context.Context,
		id uint64,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *FullProfileDTO, resErr error)

	UploadProfileImageFile(
		ctx context.Context,
		fileName string,
		fileData []byte,
	) (resMap UploadFileOut, resErr error)

	CreateByDTO(
		ctx context.Context,
		account *userEntity.Account,
		in ProfileDataInput,
	) (resItem *FullProfileDTO, resErr error)

	UpdateByDTO(
		ctx context.Context,
		id uint64,
		in ProfileDataInput,
		skipVersionCheck bool,
	) (resErr error)
}
