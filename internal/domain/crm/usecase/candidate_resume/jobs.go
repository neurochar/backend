package candidate_resume

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/internal/infra/storage"
	workflows_pb "github.com/neurochar/backend/pkg/proto_pb/common/workflows"
	"github.com/samber/lo"
	tclient "go.temporal.io/sdk/client"
)

func (uc *UsecaseImpl) JobProcessCandidatesResumesNew(ctx context.Context) (bool, error) {
	const op = "JobProcessCandidatesResumesNew"

	anyJonDone := false

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(ctx, &usecase.CandidateResumeListOptions{
			FilterStatus: lo.ToPtr(entity.CandidateResumeStatusNew),
		}, &uctypes.QueryGetListParams{
			Limit:               1,
			ForUpdateSkipLocked: true,
		})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		item := items[0]

		checkItems, err := uc.repo.FindList(ctx, &usecase.CandidateResumeListOptions{
			FilterFileHash: lo.ToPtr(item.FileHash),
			FilterStatus:   lo.ToPtr(entity.CandidateResumeStatusProcessed),
			FilterNotIDs:   lo.ToPtr([]uuid.UUID{item.ID}),
			Sort: []uctypes.SortOption[usecase.CandidateResumeListOptionsSortField]{
				{
					Field:  usecase.CandidateResumeListOptionsSortFieldUpdatedAt,
					IsDesc: true,
				},
			},
		}, &uctypes.QueryGetListParams{
			Limit: 1,
		})
		if err != nil {
			return err
		}

		if len(checkItems) > 0 && checkItems[0].UpdatedAt.After(time.Now().Add(time.Hour*24)) {
			item.AnalyzeData = deepcopy.Copy(checkItems[0].AnalyzeData).(*entity.CandidateResumeAnalyzeData)
			item.Status = entity.CandidateResumeStatusProcessed

			err = uc.repo.Update(ctx, item)
			if err != nil {
				return err
			}

			anyJonDone = true
			return nil
		}

		item.Status = entity.CandidateResumeStatusToProcess

		err = uc.repo.Update(ctx, item)
		if err != nil {
			return err
		}

		anyJonDone = true

		return nil
	})
	if err != nil {
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return anyJonDone, nil
}

func (uc *UsecaseImpl) JobProcessCandidatesResumesToProcess(ctx context.Context) (bool, error) {
	const op = "JobProcessCandidatesResumesToProcess"

	anyJonDone := false

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.FindList(
			ctx,
			&usecase.CandidateResumeListOptions{
				FilterStatus: lo.ToPtr(entity.CandidateResumeStatusToProcess),
			},
			&uctypes.QueryGetListParams{
				Limit:               1,
				ForUpdateSkipLocked: true,
			}, &usecase.CandidateResumeDTOOptions{
				FetchFile: true,
			})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		item := items[0]

		if item.File == nil || item.File.StorageFileKey == nil {
			item.Resume.Status = entity.CandidateResumeStatusProcessError
			item.Resume.ErrorText = lo.ToPtr("file is empty")

			err := uc.repo.Update(ctx, item.Resume)
			if err != nil {
				return err
			}

			anyJonDone = true
			return nil
		}

		newTaskID := uuid.NewString()
		payload := workflows_pb.WorkflowProcessResumeFileInput{
			ResumeId: item.Resume.ID.String(),
			Data: &workflows_pb.WorkflowProcessResumeFileInput_FileData{
				Name:          item.File.OriginalFileName,
				StorageBucket: string(storage.BucketCommonFiles),
				StorageKey:    *item.File.StorageFileKey,
			},
		}

		switch item.Resume.FileType {
		case entity.CandidateResumeFileTypePdf:
			payload.Data.Type = workflows_pb.WorkflowProcessResumeFileInput_FILE_TYPE_PDF
		case entity.CandidateResumeFileTypeWord:
			payload.Data.Type = workflows_pb.WorkflowProcessResumeFileInput_FILE_TYPE_PDF
		default:
			payload.Data.Type = workflows_pb.WorkflowProcessResumeFileInput_FILE_TYPE_UNSPECIFIED
		}

		_, err = uc.temporalClient.ExecuteWorkflow(
			ctx,
			tclient.StartWorkflowOptions{
				ID:        "workflow-job-" + newTaskID,
				TaskQueue: workflows_pb.WorkflowsQueue_WORKFLOWS_QUEUE_DEFAULT.String(),
			},
			workflows_pb.Workflow_WORKFLOW_PROCESS_RESUME_FILE.String(),
			&payload,
		)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		item.Resume.Status = entity.CandidateResumeStatusProcessing

		err = uc.repo.Update(ctx, item.Resume)
		if err != nil {
			return err
		}

		anyJonDone = true

		return nil
	})
	if err != nil {
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return anyJonDone, nil
}
