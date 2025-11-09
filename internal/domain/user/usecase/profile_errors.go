package usecase

import appErrors "github.com/neurochar/backend/internal/app/errors"

var ErrFilePhoto100x100BadRequest = appErrors.ErrBadRequest.Extend("photo_100x100_file_id bad request").
	WithHints("photo_100x100_file_id incorrect")
