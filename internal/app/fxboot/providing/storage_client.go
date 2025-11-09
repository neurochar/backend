package providing

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/storage"
	"github.com/neurochar/backend/internal/infra/storage/s3d"
)

// NewStorageClient - provide storage client
func NewStorageClient(cfg config.Config) (storage.Client, *s3.Client) {
	usePathStyle := !cfg.Storage.S3URLIsHost

	s3Client := s3d.NewS3Client(
		cfg.Storage.S3Endpoint,
		cfg.Storage.S3Region,
		cfg.Storage.S3AccessKey,
		cfg.Storage.S3SecretKey,
		usePathStyle,
	)

	return s3d.New(s3Client), s3Client
}
