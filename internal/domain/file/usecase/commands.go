package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	"github.com/neurochar/backend/internal/infra/storage"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) UploadAndCreateFiles(
	ctx context.Context,
	input UploadFilesIn,
) (resMap UploadFilesOut, resCancelFn func(context.Context) error, resErr error) {
	const op = "UploadAndCreateFiles"

	createdIDs := make([]uuid.UUID, 0, len(input))

	resCancelFn = func(ctx context.Context) error {
		err := uc.DeleteByIDs(ctx, createdIDs)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s.%s", uc.pkg, op, "CancelFn")
		}

		return nil
	}

	defer func() {
		if resErr != nil {
			// nolint
			_ = resCancelFn(ctx)
		}
	}()

	result := make(UploadFilesOut, len(input))

	for i, item := range input {
		if item.Process != nil {
			fileData, err := item.Process(item.FileData)
			if err != nil {
				return nil, resCancelFn, appErrors.Chainf(
					err.
						WithDetail("fileTarget", false, item.Target).
						WithDetail("fileName", false, item.FileName),
					"%s.%s", uc.pkg, op)
			}

			input[i].FileData = fileData
		}
	}

	groupID := uuid.New()

	for _, item := range input {
		fileKey, fileHash, fileMimeType, _ := uc.storageClient.FileMetaByBytes(ctx, item.FileName, item.FileData)

		file := fileEntity.NewFile(groupID, item.Target, false, item.FileName)
		file.StorageFileKey = &fileKey

		err := uc.Create(ctx, file)
		if err != nil {
			return nil, resCancelFn, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		createdIDs = append(createdIDs, file.ID)

		_, err = uc.storageClient.UploadBytes(ctx, storage.BucketCommonFiles, storage.UploadBytesInput{
			Key:         fileKey,
			Hash:        fileHash,
			ContentType: fileMimeType,
			Data:        item.FileData,
		})
		if err != nil {
			return nil, resCancelFn, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
		}

		file.UploadedToStorage = true

		err = uc.Update(ctx, file)
		if err != nil {
			return nil, resCancelFn, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		result[file.ID] = file
	}

	return result, resCancelFn, nil
}

func (uc *UsecaseImpl) Create(ctx context.Context, item *fileEntity.File) (resErr error) {
	const op = "Create"

	err := uc.repo.Create(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Update(ctx context.Context, item *fileEntity.File) (resErr error) {
	const op = "Update"

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Delete(ctx context.Context, item *fileEntity.File) (resErr error) {
	const op = "Delete"

	timeNow := time.Now()
	item.ToDeleteFromStorage = true
	item.DeletedAt = &timeNow

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) DeleteByIDs(ctx context.Context, ids []uuid.UUID) (resErr error) {
	const op = "DeleteByIDs"

	if len(ids) == 0 {
		return nil
	}

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(ctx, &ListOptions{IDs: &ids}, nil)
		if err != nil {
			return err
		}

		for _, item := range items {
			err := uc.Delete(ctx, item)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) ProcessFilesToTarget(
	ctx context.Context,
	items []ProcessFileToTargetIn,
) (map[uuid.UUID]*fileEntity.File, error) {
	const op = "ProcessFilesToTarget"

	var filesMap map[uuid.UUID]*fileEntity.File

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		filesIDs := make([]uuid.UUID, 0, len(items)*2)
		for _, item := range items {
			if item.CurrentFileID != nil {
				filesIDs = append(filesIDs, *item.CurrentFileID)
			}

			if item.NewFileID != nil {
				filesIDs = append(filesIDs, *item.NewFileID)
			}
		}

		if len(filesIDs) == 0 {
			return nil
		}

		var err error

		filesMap, err = uc.FindListInMap(ctx, &ListOptions{IDs: lo.ToPtr(filesIDs)}, nil)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		toDelete := make([]uuid.UUID, 0)
		toAssign := make([]uuid.UUID, 0)

		groupsMap := map[string]uuid.UUID{}

		for _, item := range items {
			var currentFile, newFile *fileEntity.File

			if item.CurrentFileID != nil {
				file, ok := filesMap[*item.CurrentFileID]
				if ok {
					currentFile = file
				}
			}

			if item.NewFileID != nil {
				var ok bool
				newFile, ok = filesMap[*item.NewFileID]
				if !ok {
					return ErrProcessFileIncorrect.
						WithDetail("fileID", false, item.NewFileID.String()).
						WithDetail("fileTarget", false, item.Target).
						WithHints(fmt.Sprintf("file %s not found", item.NewFileID.String()))
				}
			}

			if item.Group != "" {
				groupMapItem, ok := groupsMap[item.Group]
				var groupMapValue uuid.UUID
				if newFile != nil {
					groupMapValue = newFile.GroupID
				}

				if !ok {
					groupsMap[item.Group] = groupMapValue
				} else if ok && groupMapItem != groupMapValue {
					return ErrProcessFileIncorrect.
						WithDetail("fileID", false, item.NewFileID.String()).
						WithDetail("fileTarget", false, item.Target).
						WithHints(fmt.Sprintf("file %s incorrect group", item.NewFileID.String()))
				}
			}

			if newFile != nil && newFile.Target != item.Target {
				return ErrProcessFileIncorrect.
					WithDetail("fileID", false, item.NewFileID.String()).
					WithDetail("fileTarget", false, item.Target).
					WithHints(fmt.Sprintf("file %s target is incorrect", item.NewFileID.String()))
			}

			if newFile != nil && (currentFile == nil || currentFile.ID != newFile.ID) {
				if newFile.AssignedToTarget {
					return ErrProcessFileIncorrect.
						WithDetail("fileID", false, item.NewFileID.String()).
						WithDetail("fileTarget", false, item.Target).
						WithHints(fmt.Sprintf("file %s already assigned", item.NewFileID.String()))
				}

				toAssign = append(toAssign, newFile.ID)
			}

			if currentFile != nil && (newFile == nil || newFile.ID != currentFile.ID) {
				toDelete = append(toDelete, currentFile.ID)
			}
		}

		toDelete = lo.Uniq(toDelete)
		toAssign = lo.Uniq(toAssign)

		for _, id := range toDelete {
			err := uc.Delete(ctx, filesMap[id])
			if err != nil {
				return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
			}
		}

		for _, id := range toAssign {
			file := filesMap[id]
			file.SetAssignedToTarget(true)

			err := uc.Update(ctx, file)
			if err != nil {
				return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
			}
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return filesMap, nil
}

func (uc *UsecaseImpl) ProcessFilesSliceToTarget(
	ctx context.Context,
	oldIDs []uuid.UUID,
	newIDs []uuid.UUID,
	target string,
) (map[uuid.UUID]*fileEntity.File, error) {
	const op = "ProcessFilesSliceToTarget"

	onlyInNew, onlyInOld := lo.Difference(newIDs, oldIDs)
	allIDs := append(oldIDs, onlyInNew...)

	var filesMap map[uuid.UUID]*fileEntity.File

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		filesMap, err = uc.FindListInMap(ctx, &ListOptions{IDs: lo.ToPtr(allIDs)}, nil)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		for _, fileID := range onlyInNew {
			file, ok := filesMap[fileID]
			if !ok {
				return ErrProcessFileIncorrect.
					WithDetail("fileID", false, fileID.String()).
					WithDetail("fileTarget", false, target).
					WithHints(fmt.Sprintf("file %s not found", fileID.String()))
			}

			if file.Target != target {
				return ErrProcessFileIncorrect.
					WithDetail("fileID", false, fileID.String()).
					WithDetail("fileTarget", false, target).
					WithHints(fmt.Sprintf("file %s target is incorrect", fileID.String()))
			}

			if file.AssignedToTarget {
				return ErrProcessFileIncorrect.
					WithDetail("fileID", false, fileID.String()).
					WithDetail("fileTarget", false, target).
					WithHints(fmt.Sprintf("file %s already assigned", fileID.String()))
			}
		}

		toDelete := lo.Uniq(onlyInOld)
		toAssign := lo.Uniq(onlyInNew)

		for _, id := range toDelete {
			err := uc.Delete(ctx, filesMap[id])
			if err != nil {
				return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
			}
		}

		for _, id := range toAssign {
			file := filesMap[id]
			file.SetAssignedToTarget(true)

			err := uc.Update(ctx, file)
			if err != nil {
				return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
			}
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return filesMap, nil
}
