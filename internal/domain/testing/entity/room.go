package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RoomStatusType uint8

const (
	RoomStatusTypeNotStarted RoomStatusType = 0
	RoomStatusTypeFinished   RoomStatusType = 10
)

var RoomStatusTypeMap = map[RoomStatusType]string{
	RoomStatusTypeNotStarted: "not_started",
	RoomStatusTypeFinished:   "finished",
}

func (t RoomStatusType) String() string {
	return RoomStatusTypeMap[t]
}

type RoomTechniqueDataItem struct {
	TechniqueID uint64
	ItemData    TechniqueItemData
}

type Room struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	Status               RoomStatusType
	CandidateID          *uuid.UUID
	ProfileID            *uuid.UUID
	PersonalityTraitsMap ProfilePersonalityTraitsMap
	TechniqueData        []RoomTechniqueDataItem
	CandidateAnswerData  map[uint64]any
	RowResult            json.RawMessage
	Result               *RoomResult
	CreatedBy            *uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

func (item *Room) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func NewRoom(
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
