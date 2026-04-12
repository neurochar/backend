package mapper

import (
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	"github.com/samber/lo"
)

var roomStatusToPb = map[testingEntity.RoomStatusType]typesv1.RoomStatus{
	testingEntity.RoomStatusTypeUnspecified: typesv1.RoomStatus_ROOM_STATUS_UNSPECIFIED,
	testingEntity.RoomStatusTypeNotStarted:  typesv1.RoomStatus_ROOM_STATUS_NOT_STARTED,
	testingEntity.RoomStatusTypeFinished:    typesv1.RoomStatus_ROOM_STATUS_FINISHED,
}

var techniqueItemTypeToPb = map[testingEntity.TechniqueItemType]typesv1.TechniqueItemType{
	testingEntity.TechniqueItemTypeUnspecified:                      typesv1.TechniqueItemType_TECHNIQUE_ITEM_TYPE_UNSPECIFIED,
	testingEntity.TechniqueItemTypeQuestionWithVariantsSignleAnswer: typesv1.TechniqueItemType_TECHNIQUE_ITEM_TYPE_QUESTION_WITH_VARIANTS_SINGLE_ANSWER,
}

var personalityTraitTypeToPb = map[testingEntity.PersonalityTraitType]typesv1.PersonalityTraitType{
	testingEntity.PersonalityTraitTypeUnspecified: typesv1.PersonalityTraitType_PERSONALITY_TRAIT_TYPE_UNSPECIFIED,
	testingEntity.PersonalityTraitTypeBipolar:     typesv1.PersonalityTraitType_PERSONALITY_TRAIT_TYPE_BIPOLAR,
}

var traitPriorityToPb = map[testingEntity.TraitPriority]typesv1.PersonalityTraitPriority{
	testingEntity.TraitPriorityNone:   typesv1.PersonalityTraitPriority_PRESONALITY_TRAIT_PRIORITY_NONE,
	testingEntity.TraitPriorityLow:    typesv1.PersonalityTraitPriority_PRESONALITY_TRAIT_PRIORITY_LOW,
	testingEntity.TraitPriorityMedium: typesv1.PersonalityTraitPriority_PRESONALITY_TRAIT_PRIORITY_MEDIUM,
	testingEntity.TraitPriorityHigh:   typesv1.PersonalityTraitPriority_PRESONALITY_TRAIT_PRIORITY_HIGH,
}

var traitPriorityPbToEntity = make(
	map[typesv1.PersonalityTraitPriority]testingEntity.TraitPriority,
	len(traitPriorityToPb),
)

func init() {
	for k, v := range traitPriorityToPb {
		traitPriorityPbToEntity[v] = k
	}
}

func RoomStatusToPb(item testingEntity.RoomStatusType) typesv1.RoomStatus {
	val, ok := roomStatusToPb[item]
	if !ok {
		return typesv1.RoomStatus_ROOM_STATUS_UNSPECIFIED
	}

	return val
}

func TechniqueItemTypeToPb(item testingEntity.TechniqueItemType) typesv1.TechniqueItemType {
	val, ok := techniqueItemTypeToPb[item]
	if !ok {
		return typesv1.TechniqueItemType_TECHNIQUE_ITEM_TYPE_UNSPECIFIED
	}

	return val
}

func PersonalityTraitTypeToPb(item testingEntity.PersonalityTraitType) typesv1.PersonalityTraitType {
	val, ok := personalityTraitTypeToPb[item]
	if !ok {
		return typesv1.PersonalityTraitType_PERSONALITY_TRAIT_TYPE_UNSPECIFIED
	}

	return val
}

func TraitPriorityToPb(item testingEntity.TraitPriority) typesv1.PersonalityTraitPriority {
	val, ok := traitPriorityToPb[item]
	if !ok {
		return typesv1.PersonalityTraitPriority_PRESONALITY_TRAIT_PRIORITY_NONE
	}

	return val
}

func TraitPriorityPbToEntity(item typesv1.PersonalityTraitPriority) testingEntity.TraitPriority {
	val, ok := traitPriorityPbToEntity[item]
	if !ok {
		return testingEntity.TraitPriorityNone
	}

	return val
}

func ProfilePersonalityTraitsMapItemToPb(
	item testingEntity.ProfilePersonalityTraitsMapItem,
) typesv1.ProfilePersonalityTraitsMapItem {
	return typesv1.ProfilePersonalityTraitsMapItem{
		Priority: TraitPriorityToPb(item.Priority),
		Target:   int32(item.Target),
	}
}

func ProfilePersonalityTraitsMapToPb(
	item testingEntity.ProfilePersonalityTraitsMap,
) *typesv1.ProfilePersonalityTraitsMap {
	return &typesv1.ProfilePersonalityTraitsMap{
		Map: lo.MapValues(
			item,
			func(v testingEntity.ProfilePersonalityTraitsMapItem, _ uint64) *typesv1.ProfilePersonalityTraitsMapItem {
				return lo.ToPtr(ProfilePersonalityTraitsMapItemToPb(v))
			},
		),
	}
}

func TestingListProfileDTOToPb(item *testingUC.ProfileDTO) *typesv1.TestingListProfile {
	return &typesv1.TestingListProfile{
		Id:       item.Profile.ID.String(),
		Version:  item.Profile.Version(),
		TenantId: item.Profile.TenantID.String(),
		Name:     item.Profile.Name,
	}
}

func TestingProfileDTOToPb(item *testingUC.ProfileDTO) *typesv1.TestingProfile {
	return &typesv1.TestingProfile{
		Id:                item.Profile.ID.String(),
		Version:           item.Profile.Version(),
		TenantId:          item.Profile.TenantID.String(),
		Name:              item.Profile.Name,
		Description:       item.Profile.Description,
		PersonalityTraits: ProfilePersonalityTraitsMapToPb(item.Profile.PersonalityTraitsMap),
	}
}

func TestingPersonalityTraitsMapItemPbToEntity(
	v *typesv1.ProfilePersonalityTraitsMapItem,
) testingEntity.ProfilePersonalityTraitsMapItem {
	res := testingEntity.ProfilePersonalityTraitsMapItem{}
	if v == nil {
		return res
	}

	res.Priority = TraitPriorityPbToEntity(v.Priority)
	res.Target = int(v.Target)

	return res
}

func TestingPersonalityTraitToPb(item testingEntity.PersonalityTrait) *typesv1.PersonalityTrait {
	return &typesv1.PersonalityTrait{
		Id:             item.GetID(),
		Type:           PersonalityTraitTypeToPb(item.GetType()),
		Name:           item.GetName(),
		Description:    item.GetDescription(),
		LeftStateName:  item.GetLeftStateName(),
		RightStateName: item.GetRightStateName(),
	}
}

func CrmCandidateToTestingRoomPb(item *crmUC.CandidateDTO) *typesv1.TestingRoomCandidate {
	return &typesv1.TestingRoomCandidate{
		Id:      item.Candidate.ID.String(),
		Name:    item.Candidate.CandidateName,
		Surname: item.Candidate.CandidateSurname,
	}
}

func TestingProfileToTestingRoomPb(item *testingUC.ProfileDTO) *typesv1.TestingRoomProfile {
	return &typesv1.TestingRoomProfile{
		Id:   item.Profile.ID.String(),
		Name: item.Profile.Name,
	}
}

func TestingRoomResultToPb(item *testingEntity.RoomResult) *typesv1.TestingRoomResult {
	if item == nil {
		return nil
	}

	res := &typesv1.TestingRoomResult{
		TotalMatchTip: item.TotalMatchTip,
		Traits: lo.MapValues(item.Traits, func(v testingEntity.RoomResultTraitItem, k uint64) *typesv1.TestingRoomResultTrait {
			r := &typesv1.TestingRoomResultTrait{
				Tip: v.Tip,
			}

			if match, ok := v.Match.Float64(); ok {
				r.Match = float32(match)
			}

			return r
		}),
	}

	if totalMatch, ok := item.TotalMatch.Float64(); ok {
		res.TotalMatch = float32(totalMatch)
	}

	return res
}

func TestingRoomDTOToListPb(item *testingUC.RoomDTO) *typesv1.TestingListRoom {
	return &typesv1.TestingListRoom{
		Id:        item.Room.ID.String(),
		Version:   item.Room.Version(),
		TenantId:  item.Room.TenantID.String(),
		Status:    RoomStatusToPb(item.Room.Status),
		Candidate: CrmCandidateToTestingRoomPb(item.CandidateDTO),
		Profile:   TestingProfileToTestingRoomPb(item.ProfileDTO),
		Result:    TestingRoomResultToPb(item.Room.Result),
	}
}

func TestingRoomDTOToPb(item *testingUC.RoomDTO) *typesv1.TestingRoom {
	return &typesv1.TestingRoom{
		Id:        item.Room.ID.String(),
		Version:   item.Room.Version(),
		TenantId:  item.Room.TenantID.String(),
		Status:    RoomStatusToPb(item.Room.Status),
		Candidate: CrmCandidateToTestingRoomPb(item.CandidateDTO),
		Profile:   TestingProfileToTestingRoomPb(item.ProfileDTO),
		Result:    TestingRoomResultToPb(item.Room.Result),
		PersonalityTraits: &typesv1.ProfilePersonalityTraitsMap{
			Map: lo.MapValues(
				item.Room.PersonalityTraitsMap,
				func(v testingEntity.ProfilePersonalityTraitsMapItem, _ uint64) *typesv1.ProfilePersonalityTraitsMapItem {
					return lo.ToPtr(ProfilePersonalityTraitsMapItemToPb(v))
				},
			),
		},
	}
}
