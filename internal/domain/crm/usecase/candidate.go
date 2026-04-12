package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
)

type CandidateListOptions struct {
	FilterIDs      *[]uuid.UUID
	FilterTenantID *uuid.UUID
	SearchQuery    *string
	Sort           []uctypes.SortOption[CandidateListOptionsSortField]
}

type CandidateListOptionsSortField string

const (
	CandidateListOptionsSortFieldCreatedAt CandidateListOptionsSortField = "created_at"
)

type CandidateDTOOptions struct {
	FetchCreatedBy bool
	FetchResume    bool
}

type CandidateDTO struct {
	Candidate *entity.Candidate
	Resume    *CandidateResumeDTO
	CreatedBy *tenantUC.AccountDTO
}

type CreateCandidateDataInput struct {
	CandidateName     string
	CandidateSurname  string
	CandidateGender   entity.CandidateGender
	CandidateBirthday *time.Time
	CreatedBy         *uuid.UUID
	ResumeFileID      *uuid.UUID
}

type PatchCandidateDataInput struct {
	Version int64

	CandidateName     *string
	CandidateSurname  *string
	CandidateGender   *entity.CandidateGender
	CandidateBirthday **time.Time
	ResumeFileID      **uuid.UUID
}

type CandidateUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
		dtoOpts *CandidateDTOOptions,
	) (resCandidate *CandidateDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *CandidateListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *CandidateDTOOptions,
	) (resItems []*CandidateDTO, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *CandidateListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *CandidateDTOOptions,
	) (resItems []*CandidateDTO, total uint64, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *CandidateListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *CandidateDTOOptions,
	) (resItems map[uuid.UUID]*CandidateDTO, resErr error)

	CreateByDTO(
		ctx context.Context,
		tenantID uuid.UUID,
		in CreateCandidateDataInput,
	) (resCandidateDTO *CandidateDTO, resErr error)

	PatchByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchCandidateDataInput,
		skipVersionCheck bool,
	) (resErr error)

	Update(ctx context.Context, item *entity.Candidate) (resErr error)
}

type CandidateRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (item *entity.Candidate, err error)

	FindList(
		ctx context.Context,
		listOptions *CandidateListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Candidate, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *CandidateListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Candidate, total uint64, err error)

	Create(ctx context.Context, item *entity.Candidate) (err error)

	Update(ctx context.Context, item *entity.Candidate) (err error)
}
