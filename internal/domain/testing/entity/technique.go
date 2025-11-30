package entity

import (
	"time"

	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
)

type Technique interface {
	GetID() uint64

	GetTitle() string

	GetItemsByPersonalityTraits(traitsMap map[uint64]ProfilePersonalityTraitsMapItem) []TechniqueItemData

	CountResult(
		traitsMap map[uint64]ProfilePersonalityTraitsMapItem,
		answers map[uint64]any,
		candidateGender crmEntity.CandidateGender,
		candidateBirthday *time.Time,
	) (RoomResultTechnique, error)
}
