package httperrs

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrCantParseBody = appErrors.ErrBadRequest.Extend("cannot parse request body").WithHints("cannot parse request body")

var ErrValidation = appErrors.ErrBadRequest.Extend("validation")
