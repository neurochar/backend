package entity

import (
	"time"

	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
)

type TechniqueItemType int

const (
	TechniqueItemTypeQuestionWithVariantsSignleAnswer TechniqueItemType = 1
)

var TechniqueItemTypeMap = map[TechniqueItemType]string{
	TechniqueItemTypeQuestionWithVariantsSignleAnswer: "question_with_variants_single_answer",
}

func (t TechniqueItemType) String() string {
	return TechniqueItemTypeMap[t]
}

type TechniqueItem interface {
	GetTeqniqueID() uint64
	GetID() uint64
	GetType() TechniqueItemType
	GetTitle() string
	ValidateAnswer(answer any) error
}

type TechniqueItemQuestionWithVariants interface {
	TechniqueItem
	GetQuestion(candidateGender crmEntity.CandidateGender, candidateBirthday *time.Time) string
	GetVariants(candidateGender crmEntity.CandidateGender, candidateBirthday *time.Time) []string
}
