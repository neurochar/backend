package usecase

import (
	"context"
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/infra/storage"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) JobProcessFilesToDelete(ctx context.Context) (bool, error) {
	const op = "JobProcessFilesToDelete"

	anyJonDone := false

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(ctx, &ListOptions{
			ToDeleteFromStorage: lo.ToPtr(true),
		}, &uctypes.QueryGetListParams{
			WithDeleted:         true,
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

		if item.StorageFileKey != nil {
			check, err := uc.repo.FindList(ctx, &ListOptions{
				StorageFileKey: lo.ToPtr(item.StorageFileKey),
			}, &uctypes.QueryGetListParams{
				Limit: 1,
			})
			if err != nil {
				return err
			}

			if len(check) > 0 || !item.UploadedToStorage {
				goto updateToProcessed
			}

			err = uc.storageClient.Delete(ctx, storage.BucketCommonFiles, *item.StorageFileKey)
			if err != nil {
				return err
			}
		}

	updateToProcessed:
		item.ToDeleteFromStorage = false

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

func (uc *UsecaseImpl) JobProcessUnusedFiles(ctx context.Context, unusedTTL time.Duration) (bool, error) {
	const op = "JobProcessUnusedFiles"

	anyJonDone := false

	maxTime := time.Now().Add(-unusedTTL)

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(ctx, &ListOptions{
			AssignedToTarget: lo.ToPtr(false),
			CreatedAt: &ListOptionsCreatedAt{
				Value:   &maxTime,
				Compare: uctypes.CompareLessOrEqual,
			},
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

		err = uc.Delete(ctx, items[0])
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
