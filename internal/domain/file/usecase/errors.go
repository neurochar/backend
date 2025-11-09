package usecase

import appErrors "github.com/neurochar/backend/internal/app/errors"

var ErrProcessFileIncorrect = appErrors.ErrBadRequest.Extend("process file incorrect")
