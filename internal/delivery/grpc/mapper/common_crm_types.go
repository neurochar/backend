package mapper

import (
	"github.com/neurochar/backend/internal/delivery/grpc/mapper/helpers"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
)

var candidateResumeStatusToPb = map[crmEntity.CandidateResumeStatus]typesv1.CandidateResumeStatus{
	crmEntity.CandidateResumeStatusUnspecified:  typesv1.CandidateResumeStatus_CANDIDATE_RESUME_STATUS_UNSPECIFIED,
	crmEntity.CandidateResumeStatusNew:          typesv1.CandidateResumeStatus_CANDIDATE_RESUME_STATUS_NEW,
	crmEntity.CandidateResumeStatusToProcess:    typesv1.CandidateResumeStatus_CANDIDATE_RESUME_STATUS_TO_PROCESS,
	crmEntity.CandidateResumeStatusProcessing:   typesv1.CandidateResumeStatus_CANDIDATE_RESUME_STATUS_PROCESSING,
	crmEntity.CandidateResumeStatusProcessed:    typesv1.CandidateResumeStatus_CANDIDATE_RESUME_STATUS_PROCESSED,
	crmEntity.CandidateResumeStatusProcessError: typesv1.CandidateResumeStatus_CANDIDATE_RESUME_STATUS_PROCESS_ERROR,
}

var candidateResumeStatusPbToEntity = make(
	map[typesv1.CandidateResumeStatus]crmEntity.CandidateResumeStatus,
	len(candidateResumeStatusToPb),
)

func init() {
	for k, v := range candidateResumeStatusToPb {
		candidateResumeStatusPbToEntity[v] = k
	}
}

func CandidateResumeStatusPbToEntity(item typesv1.CandidateResumeStatus) crmEntity.CandidateResumeStatus {
	val, ok := candidateResumeStatusPbToEntity[item]
	if !ok {
		return crmEntity.CandidateResumeStatusUnspecified
	}

	return val
}

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

func CandidateDTOToPb(item *crmUC.CandidateDTO, fileUC fileUC.Usecase, isFull bool) *typesv1.Candidate {
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

	if item.Resume != nil && item.Resume.File != nil {
		resp.ResumeFile = FileToPb(item.Resume.File, fileUC, isFull)
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
