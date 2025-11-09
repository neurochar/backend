package usecase

import appErrors "github.com/neurochar/backend/internal/app/errors"

var ErrCantDeleteRoleIsSystem = appErrors.ErrConflict.Extend("can't delete role, its system").
	WithTextCode("CANT_DELETE_ROLE")
