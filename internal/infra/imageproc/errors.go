package imageproc

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrInvalidImageFile = appErrors.ErrBadRequest.Extend("invalid image file").WithTextCode("INVALID_IMAGE_FILE")

var ErrInvalidImageSize = appErrors.ErrBadRequest.Extend("invalid image size").WithTextCode("INVALID_IMAGE_SIZE")

var ErrCantConvertImage = appErrors.ErrInternal.Extend("cannot convert image")
