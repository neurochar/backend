package pg

import (
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	CandidateTable = "crm_candidate"
)

var CandidateTableFields = []string{}

func init() {
	CandidateTableFields = dbhelper.ExtractDBFields(&CandidateDBModel{})
}

type CandidateDBModel struct {
	ID                uuid.UUID  `db:"id"`
	TenantID          uuid.UUID  `db:"tenant_id"`
	CandidateName     string     `db:"candidate_name"`
	CandidateSurname  string     `db:"candidate_surname"`
	CandidateGender   uint8      `db:"candidate_gender"`
	CandidateBirthday *time.Time `db:"candidate_birthday"`
	CreatedBy         *uuid.UUID `db:"created_by"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *CandidateDBModel) ToEntity() *entity.Candidate {
	return &entity.Candidate{
		ID:                db.ID,
		TenantID:          db.TenantID,
		CandidateName:     db.CandidateName,
		CandidateSurname:  db.CandidateSurname,
		CandidateGender:   entity.CandidateGender(db.CandidateGender),
		CandidateBirthday: db.CandidateBirthday,
		CreatedBy:         db.CreatedBy,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapCandidateEntityToDBModel(entity *entity.Candidate) *CandidateDBModel {
	return &CandidateDBModel{
		ID:                entity.ID,
		TenantID:          entity.TenantID,
		CandidateName:     entity.CandidateName,
		CandidateSurname:  entity.CandidateSurname,
		CandidateGender:   uint8(entity.CandidateGender),
		CandidateBirthday: entity.CandidateBirthday,
		CreatedBy:         entity.CreatedBy,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
