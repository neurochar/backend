package crm

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
)

type OutCandidate struct {
	Version int64 `json:"_version,omitempty"`

	ID                uuid.UUID                 `json:"id"`
	TenantID          uuid.UUID                 `json:"tenantID"`
	CandidateName     string                    `json:"candidateName"`
	CandidateSurname  string                    `json:"candidateSurname"`
	CandidateGender   crmEntity.CandidateGender `json:"candidateGender"`
	CandidateBirthday *time.Time                `json:"candidateBirthday"`
}

func OutCandidateDTO(
	c *fiber.Ctx,
	candidateDTO *crmUC.CandidateDTO,
) (*OutCandidate, error) {
	out := &OutCandidate{
		Version:  candidateDTO.Candidate.Version(),
		ID:       candidateDTO.Candidate.ID,
		TenantID: candidateDTO.Candidate.TenantID,

		CandidateName:     candidateDTO.Candidate.CandidateName,
		CandidateSurname:  candidateDTO.Candidate.CandidateSurname,
		CandidateGender:   candidateDTO.Candidate.CandidateGender,
		CandidateBirthday: candidateDTO.Candidate.CandidateBirthday,
	}

	return out, nil
}
