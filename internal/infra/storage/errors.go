package storage

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrBucketAlreadyExists = appErrors.ErrConflict.Extend("bucket already exists")
