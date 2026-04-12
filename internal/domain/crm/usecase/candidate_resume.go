package usecase

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

var ErrInvalidCandidateResumeFileType = appErrors.ErrBadRequest.WithTextCode("INVALID_MIMETYPE")

var FileTargetCandidateResumeFile string = "crm:candidate:resume"

type CandidateResumeListOptions struct {
	FilterIDs           *[]uuid.UUID
	FilterCandidatesIDs *[]uuid.UUID
	FilterStatus        *entity.CandidateResumeStatus
	Sort                []uctypes.SortOption[CandidateResumeListOptionsSortField]
}

type CandidateResumeListOptionsSortField string

const (
	CandidateResumeListOptionsSortFieldCreatedAt CandidateResumeListOptionsSortField = "created_at"
)

type CandidateResumeDTOOptions struct {
	FetchFile bool
}

type CandidateResumeDTO struct {
	Resume *entity.CandidateResume
	File   *fileEntity.File
}

type CandidateResumeUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
		dtoOpts *CandidateResumeDTOOptions,
	) (resCandidate *CandidateResumeDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *CandidateResumeListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *CandidateResumeDTOOptions,
	) (resItems []*CandidateResumeDTO, resErr error)

	Create(ctx context.Context, item *entity.CandidateResume) (resErr error)

	Update(ctx context.Context, item *entity.CandidateResume) (resErr error)

	Delete(ctx context.Context, selectReq *CandidateResumeListOptions) (resErr error)

	UploadResumeFile(
		ctx context.Context,
		fileName string,
		fileData []byte,
	) ([]*fileEntity.File, error)
}

type CandidateResumeRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (item *entity.CandidateResume, err error)

	FindList(
		ctx context.Context,
		listOptions *CandidateResumeListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.CandidateResume, err error)

	Create(ctx context.Context, item *entity.CandidateResume) (err error)

	Update(ctx context.Context, item *entity.CandidateResume) (err error)
}

var MimetypesForFileTypeOnCandidateResume = map[string]entity.CandidateResumeFileType{
	"application/pdf":    entity.CandidateResumeFileTypePdf,
	"application/msword": entity.CandidateResumeFileTypeWord,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": entity.CandidateResumeFileTypeWord,
	"application/vnd.oasis.opendocument.text":                                 entity.CandidateResumeFileTypeWord,
}
