package kettel

import (
	"encoding/json"
	"time"

	"github.com/govalues/decimal"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/pkg/convert"
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
	techniqueData []entity.RoomTechniqueDataItem,
	answers map[uint64]any,
	candidateGender crmEntity.CandidateGender,
	candidateBirthday *time.Time,
) (entity.RoomResultTechnique, error) {
	result := make(entity.RoomResultTechnique)

	rawValues := make(map[uint64]int, len(traitsMap))
	for traitID := range traitsMap {
		rawValues[traitID] = 0
	}

	for index, item := range techniqueData {
		if item.TechniqueID != 1 {
			continue
		}

		techniqueItem, err := item.ItemData.GetItem()
		if err != nil {
			return nil, err
		}

		libItem, ok := ItemsLib[techniqueItem.GetID()]
		if !ok {
			return nil, appErrors.ErrInternal.WithHints("lib item not found")
		}

		v, ok := answers[uint64(index)]
		if !ok {
			return nil, appErrors.ErrInternal.WithHints("value not found in answers")
		}

		valueInt, ok := convert.ToInt(v)
		if !ok {
			return nil, appErrors.ErrInternal.WithHints("cant convert answer value to int")
		}

		if _, ok := rawValues[libItem.TraitID]; !ok {
			return nil, appErrors.ErrInternal.WithHints("trait not found in raw map values")
		}

		if valueInt < 0 || valueInt >= len(libItem.RawVariantKeys) {
			return nil, appErrors.ErrInternal.WithHints("value dont exists in RawVariantKeys")
		}

		rawValues[libItem.TraitID] += libItem.RawVariantKeys[valueInt]
	}

	for traitID, value := range rawValues {
		stenValue := convertRawToSten(value, traitID, candidateGender, candidateBirthday)

		if stenValue == -1 {
			return nil, appErrors.ErrInternal.WithHints("cant convert raw to sten")
		}

		res, err := decimal.NewFromInt64(int64(stenValue), 0, 0)
		if err != nil {
			return nil, appErrors.ErrInternal.WithParent(err).WithHints("cant convert to decimal")
		}

		result[traitID] = entity.RoomResultTechniquesItem{
			Result: res,
		}
	}

	return result, nil
}
