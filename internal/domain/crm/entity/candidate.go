package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrCandidateInvalidName = appErrors.ErrBadRequest.Extend("invalid name").WithTextCode("INVALID_NAME")

var ErrCandidateInvalidSurname = appErrors.ErrBadRequest.Extend("invalid surname").WithTextCode("INVALID_SURNAME")

type Candidate struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	CandidateName    string
	CandidateSurname string
	CreatedBy        *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Candidate) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Candidate) SetCandidateName(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrCandidateInvalidName
	}

	item.CandidateName = value

	return nil
}

func (item *Candidate) SetCandidateSurname(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrCandidateInvalidSurname
	}

	item.CandidateSurname = value

	return nil
}

func NewCandidate(
	tenantID uuid.UUID,
	createdBy *uuid.UUID,
	candidateName string,
	candidateSurname string,
) (*Candidate, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	account := &Candidate{
		ID:        uuid.New(),
		TenantID:  tenantID,
		CreatedBy: createdBy,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := account.SetCandidateName(candidateName)
	if err != nil {
		return nil, err
	}

	err = account.SetCandidateSurname(candidateSurname)
	if err != nil {
		return nil, err
	}

	return account, nil
}
