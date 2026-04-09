package mapper

import (
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	tenantEntity "github.com/neurochar/backend/internal/domain/tenant/entity"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	roomsv1 "github.com/neurochar/backend/pkg/proto_pb/public/rooms/v1"
)

func RoomToPb(
	roomDTO *testingUC.RoomDTO,
	tenant *tenantEntity.Tenant,
) (*roomsv1.Room, error) {
	if roomDTO == nil {
		return nil, nil
	}

	out := &roomsv1.Room{
		Id:     roomDTO.Room.ID.String(),
		Status: RoomStatusToPb(roomDTO.Room.Status),
	}

	if tenant != nil {
		out.TenantName = tenant.Name
	}

	candidateGender := crmEntity.CandidateGenderUnspecified
	var candidateBirthday *time.Time

	if roomDTO.CandidateDTO != nil {
		out.Candidate = &roomsv1.RoomCandidate{
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

			item := &roomsv1.RoomTechniqueItem{
				Type: TechniqueItemTypeToPb(techniqueItem.GetType()),
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

func ParseAnswerValue(v *roomsv1.AnswerValue) any {
	if v == nil {
		return nil
	}

	switch x := v.Value.(type) {
	case *roomsv1.AnswerValue_StringValue:
		return x.StringValue
	case *roomsv1.AnswerValue_IntValue:
		return x.IntValue
	case *roomsv1.AnswerValue_DoubleValue:
		return x.DoubleValue
	case *roomsv1.AnswerValue_BoolValue:
		return x.BoolValue
	default:
		return nil
	}
}
