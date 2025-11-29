package kettel

import (
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

type KettelImpl struct{}

var Kettel = KettelImpl{}

func (t *KettelImpl) GetID() uint64 {
	return 1
}

func (t *KettelImpl) GetTitle() string {
	return "16-факторный опросник Кеттела"
}

func (t *KettelImpl) GetItemsByPersonalityTraits(
	traitsMap map[uint64]entity.ProfilePersonalityTraitsMapItem,
) []entity.TechniqueItemData {
	result := make([]entity.TechniqueItemData, len(traitsMap))

	for traitID := range traitsMap {
		for _, item := range ItemsLib {
			if item.TraitID == traitID {
				result = append(result, &KettelItemDataImpl{
					ID: item.ID,
				})
			}
		}
	}

	return result
}

func (t *KettelImpl) CountResult(
	traitsMap map[uint64]entity.ProfilePersonalityTraitsMapItem,
	answers map[uint64]any,
) (entity.RoomResultTechnique, error) {
	return nil, nil
}
