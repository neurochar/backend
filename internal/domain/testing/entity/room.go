package entity

import (
	"encoding/json"
	"net"
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
	RowAnswer            json.RawMessage
	CandidateAnswerData  map[uint64]any
	Result               *RoomResult
	CreatedBy            *uuid.UUID
	FinishedIP           *net.IP
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

func (item *Room) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Room) SetCandidateID(value *uuid.UUID) error {
	item.CandidateID = value

	return nil
}

func (item *Room) SetProfileID(value *uuid.UUID) error {
	item.ProfileID = value

	return nil
}

func (item *Room) SetPersonalityTraitsMap(value ProfilePersonalityTraitsMap) error {
	if value == nil {
		value = make(ProfilePersonalityTraitsMap)
	}

	item.PersonalityTraitsMap = value

	return nil
}

func (item *Room) SetTechniqueData(value []RoomTechniqueDataItem) error {
	if value == nil {
		value = make([]RoomTechniqueDataItem, 0)
	}

	item.TechniqueData = value

	return nil
}

func NewRoom(
	tenantID uuid.UUID,
	createdBy *uuid.UUID,
	candidateID uuid.UUID,
	profileID uuid.UUID,
) (*Room, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	room := &Room{
		ID:        uuid.New(),
		TenantID:  tenantID,
		CreatedBy: createdBy,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := room.SetCandidateID(&candidateID)
	if err != nil {
		return nil, err
	}

	err = room.SetProfileID(&profileID)
	if err != nil {
		return nil, err
	}

	return room, nil
}
