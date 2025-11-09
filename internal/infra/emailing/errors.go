package emailing

import appErrors "github.com/neurochar/backend/internal/app/errors"

var (
	// ErrIncorrectTo - error for incorrect to
	ErrIncorrectTo = appErrors.ErrInternal.Extend(`incorrect field: "to"`)

	// ErrIncorrectSubject - error for incorrect subject
	ErrIncorrectSubject = appErrors.ErrInternal.Extend(`incorrect field: "subject"`)
)
