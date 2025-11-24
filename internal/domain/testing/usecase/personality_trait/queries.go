package personalitytrait

import (
	"context"

	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/traits"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/samber/lo"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uint64,
) (entity.PersonalityTrait, error) {
	const op = "FindOneByID"

	item, ok := lo.Find(traits.Traits, func(item entity.PersonalityTrait) bool {
		return item.GetID() == id
	})
	if !ok {
		return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.PersonalityTraitListOptions,
) (resItems []entity.PersonalityTrait, resErr error) {
	items := lo.Filter(traits.Traits, func(item entity.PersonalityTrait, _ int) bool {
		if listOptions == nil {
			return true
		}

		if listOptions.FilterType != nil {
			if *listOptions.FilterType != item.GetType() {
				return false
			}
		}

		return true
	})

	return items, nil
}

func (uc *UsecaseImpl) FindListInMap(
	ctx context.Context,
	listOptions *usecase.PersonalityTraitListOptions,
) (resItems map[uint64]entity.PersonalityTrait, resErr error) {
	const op = "FindListInMap"

	items, err := uc.FindList(ctx, listOptions)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := lo.SliceToMap(items, func(item entity.PersonalityTrait) (uint64, entity.PersonalityTrait) {
		return item.GetID(), item
	})

	return result, nil
}
