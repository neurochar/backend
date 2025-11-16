package cabinet

import appErrors "github.com/neurochar/backend/internal/app/errors"

var ErrPasswordsMismatch = appErrors.ErrBadRequest.WithTextCode("PASSWORDS_MISMATCH")
