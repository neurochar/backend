package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var (
	ErrProfileInvalidName    = appErrors.ErrBadRequest.Extend("invalid name")
	ErrProfileInvalidSurname = appErrors.ErrBadRequest.Extend("invalid surname")
)

type Profile struct {
	ID                 uint64
	AccountID          uuid.UUID
	Name               string
	Surname            string
	Photo100x100FileID *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Profile) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Profile) SetName(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrProfileInvalidName
	}

	item.Name = value

	return nil
}

func (item *Profile) SetSurname(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrProfileInvalidSurname
	}

	item.Surname = value

	return nil
}

func NewProfile(accountID uuid.UUID, name string, surname string) (*Profile, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Profile{
		AccountID:          accountID,
		Name:               name,
		Surname:            surname,
		Photo100x100FileID: nil,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := item.SetName(name)
	if err != nil {
		return nil, err
	}

	err = item.SetSurname(surname)
	if err != nil {
		return nil, err
	}

	return item, nil
}
