package kettel

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

type KettelItemAnswerImpl struct {
	ID             uint64
	TraitID        uint64
	MaleQuestion   string
	FemaleQuestion string
	Variants       []string
	RawVariantKeys []int
}

var _ entity.TechniqueItemQuestionWithVariants = (*KettelItemAnswerImpl)(nil)

func (i *KettelItemAnswerImpl) GetID() uint64 {
	return i.ID
}

func (i *KettelItemAnswerImpl) GetTeqniqueID() uint64 {
	return 1
}

func (i *KettelItemAnswerImpl) GetTitle() string {
	return i.MaleQuestion
}

func (i *KettelItemAnswerImpl) GetType() entity.TechniqueItemType {
	return entity.TechniqueItemTypeQuestionWithVariantsSignleAnswer
}

func (i *KettelItemAnswerImpl) GetQuestion(candidate *crmEntity.Candidate) string {
	if candidate.CandidateGender == crmEntity.CandidateGenderFemale && i.FemaleQuestion != "" {
		return i.FemaleQuestion
	}

	return i.MaleQuestion
}

func (i *KettelItemAnswerImpl) GetVariants(_ *crmEntity.Candidate) []string {
	return i.Variants
}

func (i *KettelItemAnswerImpl) ValidateAnswer(answer any) error {
	variant, ok := answer.(int)
	if !ok {
		return appErrors.ErrBadRequest
	}

	if variant < 0 || variant > len(i.Variants)-1 {
		return appErrors.ErrBadRequest
	}

	return nil
}
