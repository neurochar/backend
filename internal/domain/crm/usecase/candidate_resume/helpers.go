package candidate_resume

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
)

func (uc *UsecaseImpl) entitiesToDTO(
	ctx context.Context,
	items []*entity.CandidateResume,
	dtoOpts *usecase.CandidateResumeDTOOptions,
) ([]*usecase.CandidateResumeDTO, error) {
	const op = "entitiesToDTO"

	filesMap := make(map[uuid.UUID]*fileEntity.File, 0)
	filesIDs := make([]uuid.UUID, 0)

	for _, item := range items {
		filesIDs = append(filesIDs, item.FilesIDs()...)
	}

	if (dtoOpts == nil || dtoOpts.FetchFile) && len(filesIDs) > 0 {
		var err error
		filesMap, err = uc.fileUC.FindListInMap(ctx, &fileUC.ListOptions{
			IDs: &filesIDs,
		}, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}
	}

	out := make([]*usecase.CandidateResumeDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.CandidateResumeDTO{
			Resume: item,
		}

		if dtoOpts == nil || dtoOpts.FetchFile {

			file, ok := filesMap[item.FileID]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("file not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.File = file

		}

		out = append(out, resItem)
	}

	return out, nil
}
