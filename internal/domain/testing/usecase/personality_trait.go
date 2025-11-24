package usecase

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

var (
	ErrPersonalityTraitsMapEmpty = appErrors.ErrBadRequest.WithTextCode("PERSONALITY_TRAITS_EMPTY").
					WithHints("personality traits map is empty")
	ErrPersonalityTraitsMapIncorrect = appErrors.ErrBadRequest.WithTextCode("PERSONALITY_TRAITS_INCORRECT").
						WithHints("personality traits map is incorrect")
)

type PersonalityTraitListOptions struct {
	FilterType *entity.PersonalityTraitType
}

type PersonalityTraitUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uint64,
	) (res entity.PersonalityTrait, resErr error)

	FindList(
		ctx context.Context,
		listOptions *PersonalityTraitListOptions,
	) (resItems []entity.PersonalityTrait, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *PersonalityTraitListOptions,
	) (resItems map[uint64]entity.PersonalityTrait, resErr error)

	ValidatePersonalityTraitsMap(traitsMap entity.ProfilePersonalityTraitsMap) error
}
