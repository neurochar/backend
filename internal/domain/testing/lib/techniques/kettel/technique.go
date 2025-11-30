package kettel

import (
	"encoding/json"
	"time"

	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/samber/lo/mutable"
)

type KettelImpl struct{}

var Kettel = KettelImpl{}

func (t *KettelImpl) GetID() uint64 {
	return 1
}

func (t *KettelImpl) GetTitle() string {
	return "16-факторный опросник Кеттела"
}

func (t *KettelImpl) MakeDataItemFromRaw(raw json.RawMessage) (*KettelItemDataImpl, error) {
	dataItem := &KettelItemDataImpl{}

	err := json.Unmarshal(raw, dataItem)
	if err != nil {
		return nil, err
	}

	return dataItem, nil
}

func (t *KettelImpl) GetItemsByPersonalityTraits(
	traitsMap map[uint64]entity.ProfilePersonalityTraitsMapItem,
) []entity.TechniqueItemData {
	result := make([]entity.TechniqueItemData, 0, len(traitsMap))

	for traitID := range traitsMap {
		for _, item := range ItemsLib {
			if item.TraitID == traitID {
				result = append(result, &KettelItemDataImpl{
					ID: item.ID,
				})
			}
		}
	}

	mutable.Shuffle(result)

	return result
}

func (t *KettelImpl) CountResult(
	traitsMap map[uint64]entity.ProfilePersonalityTraitsMapItem,
	answers map[uint64]any,
	candidateGender crmEntity.CandidateGender,
	candidateBirthday *time.Time,
) (entity.RoomResultTechnique, error) {
	return nil, nil
}
