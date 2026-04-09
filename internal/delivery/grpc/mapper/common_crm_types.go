package mapper

import (
	"github.com/neurochar/backend/internal/delivery/grpc/mapper/helpers"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
)

func CandidateGenderToPb(item crmEntity.CandidateGender) typesv1.Gender {
	switch item {
	case crmEntity.CandidateGenderMale:
		return typesv1.Gender_GENDER_MALE
	case crmEntity.CandidateGenderFemale:
		return typesv1.Gender_GENDER_FEMALE
	default:
		return typesv1.Gender_GENDER_UNSPECIFIED
	}
}

func CandidateDTOToPb(item *crmUC.CandidateDTO) *typesv1.Candidate {
	if item == nil {
		return nil
	}

	resp := &typesv1.Candidate{
		Id:       item.Candidate.ID.String(),
		Version:  item.Candidate.Version(),
		TenantId: item.Candidate.TenantID.String(),
		Name:     item.Candidate.CandidateName,
		Surname:  item.Candidate.CandidateSurname,
		Gender:   CandidateGenderToPb(item.Candidate.CandidateGender),
		Birthday: helpers.TimePtrToPbDate(item.Candidate.CandidateBirthday),
	}

	return resp
}

func CandidateDTOToListPb(item *crmUC.CandidateDTO) *typesv1.ListCandidate {
	if item == nil {
		return nil
	}

	resp := &typesv1.ListCandidate{
		Id:       item.Candidate.ID.String(),
		Version:  item.Candidate.Version(),
		TenantId: item.Candidate.TenantID.String(),
		Name:     item.Candidate.CandidateName,
		Surname:  item.Candidate.CandidateSurname,
		Gender:   CandidateGenderToPb(item.Candidate.CandidateGender),
		Birthday: helpers.TimePtrToPbDate(item.Candidate.CandidateBirthday),
	}

	return resp
}
