package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrProfileInvalidName = appErrors.ErrBadRequest.Extend("invalid name").WithTextCode("INVALID_NAME")

type ProfilePersonalityTraitsMapItem struct {
	Priority int `json:"priority"`
	Target   int `json:"target"`
}

type ProfilePersonalityTraitsMap map[uint64]ProfilePersonalityTraitsMapItem

type Profile struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	Name                 string
	PersonalityTraitsMap ProfilePersonalityTraitsMap
	CreatedBy            *uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
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

func (item *Profile) SetPersonalityTraitsMap(value ProfilePersonalityTraitsMap) error {
	if value == nil {
		value = make(ProfilePersonalityTraitsMap)
	}

	item.PersonalityTraitsMap = value

	return nil
}

func NewProfile(
	tenantID uuid.UUID,
	createdBy *uuid.UUID,
	name string,
	personalityTraitsMap ProfilePersonalityTraitsMap,
) (*Profile, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	profile := &Profile{
		ID:        uuid.New(),
		TenantID:  tenantID,
		CreatedBy: createdBy,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := profile.SetName(name)
	if err != nil {
		return nil, err
	}

	err = profile.SetPersonalityTraitsMap(personalityTraitsMap)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
