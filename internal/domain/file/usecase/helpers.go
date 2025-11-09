package usecase

import (
	"fmt"
	"strings"

	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	"github.com/neurochar/backend/internal/infra/storage"
)

func (uc *UsecaseImpl) FileUrl(file *fileEntity.File) string {
	var builder strings.Builder
	if uc.cfg.Storage.S3URLIsHost {
		builder.WriteString(
			fmt.Sprintf("%s%s%s/", uc.cfg.Storage.S3URLHostPrefix, storage.BucketCommonFiles, uc.cfg.Storage.S3URLHostPostfix))
	} else {
		builder.WriteString(fmt.Sprintf("%s/%s/", uc.cfg.Storage.S3URL, storage.BucketCommonFiles))
	}
	if file.StorageFileKey != nil {
		builder.WriteString(*file.StorageFileKey)
	}

	return builder.String()
}
