package usecase

import (
	"context"
	"net/netip"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
)

var (
	ErrRegistrationDemoAlreadyExists = appErrors.ErrBadRequest.WithTextCode("DEMO_ALREADY_EXISTS")
	ErrRegistrationAlreadyFinished   = appErrors.ErrBadRequest.WithTextCode("ALREADY_FINISHED")
)

type RegistrationListOptions struct {
	FilterEmail      *string
	FilterTariff     *uint64
	FilterIsFinished *bool
}

type CreateRegistrationIn struct {
	Email     string
	Tariff    uint64
	RequestIP *netip.Addr
}

type FinishRegistrationIn struct {
	TenantTextID   string
	ProfileName    string
	ProfileSurname string
}

type RegistrationUsecase interface {
	FindList(
		ctx context.Context,
		listOptions *RegistrationListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Registration, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *RegistrationListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Registration, total uint64, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *entity.Registration, resErr error)

	CreateByDTO(
		ctx context.Context,
		in CreateRegistrationIn,
		requestIP *netip.Addr,
	) (resItem *entity.Registration, resErr error)

	FinishByDTO(
		ctx context.Context,
		id uuid.UUID,
		in FinishRegistrationIn,
		requestIP *netip.Addr,
	) (resTenant *entity.Tenant, resErr error)
}

type RegistrationRepository interface {
	FindList(
		ctx context.Context,
		listOptions *RegistrationListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Registration, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *RegistrationListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Registration, total uint64, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resFile *entity.Registration, resErr error)

	Create(ctx context.Context, item *entity.Registration) (resErr error)

	Update(ctx context.Context, item *entity.Registration) (resErr error)

	Delete(ctx context.Context, listOptions *RegistrationListOptions) (res uint64, resErr error)

	DeleteByID(ctx context.Context, id uuid.UUID) (resErr error)
}
