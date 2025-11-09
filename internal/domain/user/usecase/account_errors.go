package usecase

import appErrors "github.com/neurochar/backend/internal/app/errors"

var ErrRoleNotFound = appErrors.ErrBadRequest.Extend("role not found").WithTextCode("ROLE_NOT_FOUND")

var ErrCodeInvalid = appErrors.ErrBadRequest.Extend("code invalid").WithTextCode("CODE_INVALID")

var ErrCodeExpired = appErrors.ErrBadRequest.Extend("code expired").WithTextCode("CODE_EXPIRED")
