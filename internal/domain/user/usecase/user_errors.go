package usecase

import appErrors "github.com/neurochar/backend/internal/app/errors"

var ErrCantDeleteRoleAccountsExists = appErrors.ErrConflict.Extend("can't delete role, accounts exists").
	WithTextCode("CANT_DELETE_ROLE_ACCOUNTS_EXISTS")
