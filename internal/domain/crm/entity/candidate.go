package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrCandidateInvalidName = appErrors.ErrBadRequest.Extend("invalid name").WithTextCode("INVALID_NAME")

var ErrCandidateInvalidSurname = appErrors.ErrBadRequest.Extend("invalid surname").WithTextCode("INVALID_SURNAME")

var ErrCandidateGenderUnknown = appErrors.ErrBadRequest.Extend("gender unknown").WithTextCode("GENDER_UNKNOWN")

type CandidateGender uint8

const (
	CandidateGenderUnknown CandidateGender = 0
	CandidateGenderMale    CandidateGender = 1
	CandidateGenderFemale  CandidateGender = 2
)

func CandidateGenderFromUint8(value uint8) (CandidateGender, error) {
	switch value {
	case 0:
		return CandidateGenderUnknown, nil
	case 1:
		return CandidateGenderMale, nil
	case 2:
		return CandidateGenderFemale, nil
	default:
		return CandidateGenderUnknown, ErrCandidateGenderUnknown
	}
}

type Candidate struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	CandidateName     string
	CandidateSurname  string
	CandidateGender   CandidateGender
	CandidateBirthday *time.Time
	CreatedBy         *uuid.UUID

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

func (item *Candidate) SetCandidateGender(value CandidateGender) error {
	item.CandidateGender = value

	return nil
}

func (item *Candidate) SetCandidateBirthday(value *time.Time) error {
	item.CandidateBirthday = value

	return nil
}

func (item *Candidate) CalcAge(now time.Time) *int {
	if item.CandidateBirthday == nil {
		return nil
	}

	bd := time.Date(
		item.CandidateBirthday.Year(),
		item.CandidateBirthday.Month(),
		item.CandidateBirthday.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	age := n.Year() - bd.Year()
	if n.Month() < bd.Month() || (n.Month() == bd.Month() && n.Day() < bd.Day()) {
		age--
	}

	return &age
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
