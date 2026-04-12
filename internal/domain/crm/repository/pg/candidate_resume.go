package pg

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	CandidateResumeTable = "crm_candidates_resume"
)

var CandidateResumeTableFields = []string{}

func init() {
	CandidateResumeTableFields = dbhelper.ExtractDBFields(&CandidateResumeDBModel{})
}

type CandidateResumeDBModel struct {
	ID          uuid.UUID       `db:"id"`
	TenantID    uuid.UUID       `db:"tenant_id"`
	Status      uint8           `db:"status"`
	CandidateID *uuid.UUID      `db:"candidate_id"`
	FileID      uuid.UUID       `db:"file_id"`
	FileHash    string          `db:"file_hash"`
	FileType    uint8           `db:"file_type"`
	AnalyzeData json.RawMessage `db:"analyze_data"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *CandidateResumeDBModel) ToEntity() *entity.CandidateResume {
	status := entity.CandidateResumeStatusFromUint8(db.Status)

	fileType := entity.CandidateResumeFileTypeFromUint8(db.FileType)

	analyzeData := &entity.CandidateResumeAnalyzeData{}
	err := json.Unmarshal(db.AnalyzeData, &analyzeData)
	if err != nil {
		panic(err)
	}

	return &entity.CandidateResume{
		ID:          db.ID,
		TenantID:    db.TenantID,
		Status:      status,
		CandidateID: db.CandidateID,
		FileID:      db.FileID,
		FileHash:    db.FileHash,
		FileType:    fileType,
		AnalyzeData: analyzeData,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapCandidateResumeEntityToDBModel(entity *entity.CandidateResume) *CandidateResumeDBModel {
	analyzeData, err := json.Marshal(entity.AnalyzeData)
	if err != nil {
		panic(err)
	}

	return &CandidateResumeDBModel{
		ID:          entity.ID,
		TenantID:    entity.TenantID,
		Status:      uint8(entity.Status),
		CandidateID: entity.CandidateID,
		FileID:      entity.FileID,
		FileHash:    entity.FileHash,
		FileType:    uint8(entity.FileType),
		AnalyzeData: analyzeData,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
