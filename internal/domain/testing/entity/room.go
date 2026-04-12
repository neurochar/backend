package entity

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/convert"
)

type RoomStatusType uint8

const (
	RoomStatusTypeUnspecified RoomStatusType = 0
	RoomStatusTypeNotStarted  RoomStatusType = 1
	RoomStatusTypeFinished    RoomStatusType = 10
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
	Result               *RoomResult
	ResultIndex          *int
	CreatedBy            *uuid.UUID
	FinishedIP           *netip.Addr
	FinishedAt           *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
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

func (item *Room) SetResultIndex(value *int) error {
	item.ResultIndex = value

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

func (item *Room) SetCandidateAnswerData(value map[uint64]any) error {
	answerData := make(map[uint64]any, len(value))

	if item.TechniqueData == nil {
		return nil
	}

	for i, v := range value {
		if int(i) >= 0 && int(i) < len(item.TechniqueData) {
			techniqueDataItem := item.TechniqueData[i]

			techniqueItem, err := techniqueDataItem.ItemData.GetItem()
			if err != nil {
				return err
			}

			switch techniqueItem.GetType() {
			case TechniqueItemTypeQuestionWithVariantsSignleAnswer:
				valueInt, ok := convert.ToInt(v)
				if !ok {
					return appErrors.ErrBadRequest
				}

				answerData[i] = valueInt
			}
		}
	}

	if len(answerData) != len(item.TechniqueData) {
		return appErrors.ErrBadRequest
	}

	item.CandidateAnswerData = answerData

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
		Status:    RoomStatusTypeNotStarted,
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
