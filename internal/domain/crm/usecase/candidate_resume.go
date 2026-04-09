package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

type CandidateResumeListOptions struct {
	FilterResumeID *uuid.UUID
}

type CandidateResumeDTOOptions struct {
	FetchFile bool
}

type CandidateResumeDTO struct {
	Resume *entity.CandidateResume
	File   *fileEntity.File
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
