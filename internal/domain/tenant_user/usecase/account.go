package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	tenantEntity "github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
)

type AccountListOptions struct {
	FilterTenantID *uuid.UUID
}

type AccountDTOOptions struct {
	FetchTenant     bool
	FetchPhotoFiles bool
}

type AccountDTO struct {
	Account                  *entity.Account
	Tenant                   *tenantEntity.Tenant
	Role                     *entity.Role
	ProfilePhoto100x100File  *fileEntity.File
	ProfilePhotoOriginalFile *fileEntity.File
}

type AccountDataInputProfilePhotos struct {
	Photo100x100FileID  *uuid.UUID
	PhotoOriginalFileID *uuid.UUID
}

type CreateAccountDataInput struct {
	Email           string
	Password        string
	RoleID          uint64
	IsConfirmed     bool
	IsEmailVerified bool
	IsBlocked       bool
	ProfileName     string
	ProfileSurname  string
	ProfilePhotos   *AccountDataInputProfilePhotos
}

type PatchAccountDataInput struct {
	Version int64

	Email           *string
	Password        *string
	RoleID          *uint64
	IsConfirmed     *bool
	IsEmailVerified *bool
	IsBlocked       *bool
	ProfileName     *string
	ProfileSurname  *string
	ProfilePhotos   *AccountDataInputProfilePhotos
}

type AccountCodeListOptions struct {
	FilterAccountID *uuid.UUID
	FilterType      *entity.AccountCodeType
	FilterIsActive  *bool
}

type ProfileFileTarget string

const (
	FileTargetProfilePhoto100x100  ProfileFileTarget = "profile:photo:100x100"
	FileTargetProfilePhotoOriginal ProfileFileTarget = "profile:photo:original"
)

var ErrRoleNotFound = appErrors.ErrBadRequest.Extend("role not found").WithTextCode("ROLE_NOT_FOUND")

var ErrCodeInvalid = appErrors.ErrBadRequest.Extend("code invalid").WithTextCode("CODE_INVALID")

var ErrCodeExpired = appErrors.ErrBadRequest.Extend("code expired").WithTextCode("CODE_EXPIRED")

type AccountUsecase interface {
	FindOneByEmail(
		ctx context.Context,
		tenantID uuid.UUID,
		email string,
		queryParams *uctypes.QueryGetOneParams,
		dtoOpts *AccountDTOOptions,
	) (resAccount *AccountDTO, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
		dtoOpts *AccountDTOOptions,
	) (resAccount *AccountDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *AccountDTOOptions,
	) (resItems []*AccountDTO, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *AccountDTOOptions,
	) (resItems []*AccountDTO, total uint64, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *AccountDTOOptions,
	) (resItems map[uuid.UUID]*AccountDTO, resErr error)

	CreateAccountByDTO(
		ctx context.Context,
		tenantID uuid.UUID,
		in CreateAccountDataInput,
		author *AccountDTO,
		requestIP net.IP,
	) (resAccountDTO *AccountDTO, activationCode *entity.AccountCode, resErr error)

	PatchAccountByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchAccountDataInput,
		skipVersionCheck bool,
	) (resErr error)

	UpdateAccount(ctx context.Context, item *entity.Account) (resErr error)

	VerifyAccountEmailByCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
	) (resErr error)

	RequestPasswordRecoveryByEmail(
		ctx context.Context,
		tenantID uuid.UUID,
		email string,
		requestIP net.IP,
	) (resCode *entity.AccountCode, resErr error)

	UpdatePasswordByRecoveryCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
		newPassword string,
	) (resCode *entity.AccountCode, resErr error)

	CheckCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
	) (resCode *entity.AccountCode, resErr error)

	UpdateCode(
		ctx context.Context,
		code *entity.AccountCode,
	) (resErr error)

	UploadProfileImageFile(
		ctx context.Context,
		fileName string,
		fileData []byte,
	) (resMap []*fileEntity.File, resErr error)
}

type AccountRepository interface {
	FindOneByEmail(
		ctx context.Context,
		tenantID uuid.UUID,
		email string,
		queryParams *uctypes.QueryGetOneParams,
	) (account *entity.Account, err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (account *entity.Account, err error)

	FindList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Account, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Account, total uint64, err error)

	Create(ctx context.Context, item *entity.Account) (err error)

	Update(ctx context.Context, item *entity.Account) (err error)
}

type AccountCodeRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (item *entity.AccountCode, err error)

	FindList(
		ctx context.Context,
		listOptions *AccountCodeListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.AccountCode, err error)

	Create(ctx context.Context, item *entity.AccountCode) (err error)

	Update(ctx context.Context, item *entity.AccountCode) (err error)
}
