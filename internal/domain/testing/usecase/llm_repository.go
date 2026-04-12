package usecase

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

var (
	ErrLLMInvalidResponse = appErrors.ErrInternal.WithTextCode("LLM_INVALID_RESPONSE")
	ErrLLMBadRequest      = appErrors.ErrBadRequest
)

type GenerateProfileTraitsMapByDescriptionRequest struct {
	Role        string
	Description string
}

type GenerateProfileTraitsMapByDescriptionResponse struct {
	TraitsMap entity.ProfilePersonalityTraitsMap
}

type LLMRepository interface {
	GenerateProfileDescriptionByName(ctx context.Context, name string) (string, error)

	GenerateProfileTraitsMapByDescription(
		ctx context.Context,
		req *GenerateProfileTraitsMapByDescriptionRequest,
	) (*GenerateProfileTraitsMapByDescriptionResponse, error)
}
