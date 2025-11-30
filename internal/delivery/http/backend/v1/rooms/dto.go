package rooms

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	tenantEntity "github.com/neurochar/backend/internal/domain/tenant/entity"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
)

type OutRoomCandidate struct {
	CandidateName string `json:"candidateName"`
}

type OutRoomTechniqueItem struct {
	Type     testingEntity.TechniqueItemType `json:"type"`
	Question string                          `json:"question,omitempty"`
	Variants []string                        `json:"variants,omitempty"`
}

type OutRoom struct {
	ID            uuid.UUID                    `json:"id"`
	Status        testingEntity.RoomStatusType `json:"status"`
	TenantName    string                       `json:"tenantName"`
	Candidate     *OutRoomCandidate            `json:"candidate"`
	TechniqueData []OutRoomTechniqueItem       `json:"techniqueData"`
}

func OutRoomDTO(
	c *fiber.Ctx,
	roomDTO *testingUC.RoomDTO,
	tenant *tenantEntity.Tenant,
) (*OutRoom, error) {
	out := &OutRoom{
		ID:         roomDTO.Room.ID,
		TenantName: tenant.Name,
	}

	candidateGender := crmEntity.CandidateGenderUnknown
	var candidateBirthday *time.Time

	if roomDTO.CandidateDTO != nil {
		out.Candidate = &OutRoomCandidate{
			CandidateName: roomDTO.CandidateDTO.Candidate.CandidateName,
		}

		candidateGender = roomDTO.CandidateDTO.Candidate.CandidateGender
		candidateBirthday = roomDTO.CandidateDTO.Candidate.CandidateBirthday
	}

	if roomDTO.Room.Status != testingEntity.RoomStatusTypeFinished {
		for _, techniqueDataItem := range roomDTO.Room.TechniqueData {
			techniqueItem, err := techniqueDataItem.ItemData.GetItem()
			if err != nil {
				return nil, err
			}

			item := OutRoomTechniqueItem{
				Type: techniqueItem.GetType(),
			}

			if techniqueItem.GetType() == testingEntity.TechniqueItemTypeQuestionWithVariantsSignleAnswer {
				itemQuestionWithVariantsSignleAnswer, ok := techniqueItem.(testingEntity.TechniqueItemQuestionWithVariants)
				if !ok {
					return nil, appErrors.ErrInternal
				}

				item.Question = itemQuestionWithVariantsSignleAnswer.GetQuestion(candidateGender, candidateBirthday)
				item.Variants = itemQuestionWithVariantsSignleAnswer.GetVariants(candidateGender, candidateBirthday)
			}

			out.TechniqueData = append(out.TechniqueData, item)
		}
	}

	return out, nil
}
