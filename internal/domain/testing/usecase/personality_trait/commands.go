package personalitytrait

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/traits"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) ValidatePersonalityTraitsMap(traitsMap entity.ProfilePersonalityTraitsMap) error {
	const op = "ValidatePersonalityTraitsMap"

	if len(traitsMap) == 0 {
		return appErrors.Chainf(usecase.ErrPersonalityTraitsMapEmpty, "%s.%s", uc.pkg, op)
	}

	for id, trait := range traitsMap {
		item, ok := lo.Find(traits.Traits, func(item entity.PersonalityTrait) bool {
			return item.GetID() == id
		})
		if !ok {
			return appErrors.Chainf(usecase.ErrPersonalityTraitsMapIncorrect, "%s.%s", uc.pkg, op)
		}

		if item.GetType() != entity.PersonalityTraitTypeBipolar {
			return appErrors.Chainf(usecase.ErrPersonalityTraitsMapIncorrect, "%s.%s", uc.pkg, op)
		}

		if trait.Priority < 0 || trait.Priority > 3 {
			return appErrors.Chainf(usecase.ErrPersonalityTraitsMapIncorrect, "%s.%s", uc.pkg, op)
		}

		if trait.Target < 0 || trait.Target > 10 {
			return appErrors.Chainf(usecase.ErrPersonalityTraitsMapIncorrect, "%s.%s", uc.pkg, op)
		}
	}

	return nil
}
