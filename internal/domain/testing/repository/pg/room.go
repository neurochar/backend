package pg

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques/kettel"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	RoomTable = "testing_room"
)

var RoomTableFields = []string{}

func init() {
	RoomTableFields = dbhelper.ExtractDBFields(&RoomDBModel{})
}

type RoomDBModel struct {
	ID                   uuid.UUID       `db:"id"`
	TenantID             uuid.UUID       `db:"tenant_id"`
	Status               uint8           `db:"status"`
	CandidateID          *uuid.UUID      `db:"candidate_id"`
	ProfileID            *uuid.UUID      `db:"profile_id"`
	PersonalityTraitsMap json.RawMessage `db:"personality_traits_map"`
	TechniqueData        json.RawMessage `db:"technique_data"`
	CandidateAnswerData  json.RawMessage `db:"candidate_answer_data"`
	Result               json.RawMessage `db:"result"`
	ResultIndex          *int            `db:"result_index"`
	CreatedBy            *uuid.UUID      `db:"created_by"`
	FinishedIP           *netip.Addr     `db:"finished_ip"`
	FinishedAt           *time.Time      `db:"finished_at"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type RoomDBModelTechniqueDataItem struct {
	TechniqueID uint64
	ItemData    json.RawMessage
}

func (db *RoomDBModel) ToEntity() *entity.Room {
	traitsMap := make(entity.ProfilePersonalityTraitsMap)
	err := json.Unmarshal(db.PersonalityTraitsMap, &traitsMap)
	if err != nil {
		panic(err)
	}

	techniqueDataRow := make([]RoomDBModelTechniqueDataItem, 0)
	err = json.Unmarshal(db.TechniqueData, &techniqueDataRow)
	if err != nil {
		panic(err)
	}

	techniqueData := make([]entity.RoomTechniqueDataItem, 0, len(techniqueDataRow))
	for _, item := range techniqueDataRow {
		dataItem := entity.RoomTechniqueDataItem{
			TechniqueID: item.TechniqueID,
		}

		switch item.TechniqueID {
		case 1:
			itemData, err := kettel.Kettel.MakeDataItemFromRaw(item.ItemData)
			if err != nil {
				panic(err)
			}

			dataItem.ItemData = itemData
		default:
			panic(fmt.Sprintf("technique id=%d not found", item.TechniqueID))
		}

		techniqueData = append(techniqueData, dataItem)
	}

	var candidateAnswerData map[uint64]any
	if db.CandidateAnswerData != nil {
		candidateAnswerData = make(map[uint64]any)
		err = json.Unmarshal(db.CandidateAnswerData, &candidateAnswerData)
		if err != nil {
			panic(err)
		}
	}

	var result *entity.RoomResult
	if db.Result != nil {
		result = &entity.RoomResult{}
		err = json.Unmarshal(db.Result, result)
		if err != nil {
			panic(err)
		}
	}

	return &entity.Room{
		ID:                   db.ID,
		TenantID:             db.TenantID,
		Status:               entity.RoomStatusType(db.Status),
		CandidateID:          db.CandidateID,
		ProfileID:            db.ProfileID,
		PersonalityTraitsMap: traitsMap,
		TechniqueData:        techniqueData,
		CandidateAnswerData:  candidateAnswerData,
		Result:               result,
		ResultIndex:          db.ResultIndex,
		CreatedBy:            db.CreatedBy,
		FinishedIP:           db.FinishedIP,
		FinishedAt:           db.FinishedAt,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapRoomEntityToDBModel(entity *entity.Room) *RoomDBModel {
	traitsMap, err := json.Marshal(entity.PersonalityTraitsMap)
	if err != nil {
		panic(err)
	}

	techniqueData, err := json.Marshal(entity.TechniqueData)
	if err != nil {
		panic(err)
	}

	var candidateAnswerData json.RawMessage
	if entity.CandidateAnswerData != nil {
		candidateAnswerData, err = json.Marshal(entity.CandidateAnswerData)
		if err != nil {
			panic(err)
		}
	}

	var result json.RawMessage
	if entity.Result != nil {
		result, err = json.Marshal(entity.Result)
		if err != nil {
			panic(err)
		}
	}

	return &RoomDBModel{
		ID:                   entity.ID,
		TenantID:             entity.TenantID,
		Status:               uint8(entity.Status),
		CandidateID:          entity.CandidateID,
		ProfileID:            entity.ProfileID,
		PersonalityTraitsMap: traitsMap,
		TechniqueData:        techniqueData,
		CandidateAnswerData:  candidateAnswerData,
		Result:               result,
		ResultIndex:          entity.ResultIndex,
		CreatedBy:            entity.CreatedBy,
		FinishedAt:           entity.FinishedAt,
		FinishedIP:           entity.FinishedIP,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
