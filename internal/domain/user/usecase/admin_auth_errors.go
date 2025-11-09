package usecase

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrPasswordIncorrect = appErrors.ErrUnauthorized.Extend("password incorrect")

var ErrAccountNotConfirmed = appErrors.ErrUnauthorized.Extend("account not confirmed")

var ErrAccountBlocked = appErrors.ErrUnauthorized.Extend("account blocked")

var ErrAccessDenied = appErrors.ErrUnauthorized.Extend("access denied")

var ErrInvalidToken = appErrors.ErrUnauthorized.Extend("invalid token")

var ErrExpiredToken = appErrors.ErrUnauthorized.Extend("expired token").WithTextCode("EXPIRED_TOKEN")
