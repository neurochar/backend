package usecase

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
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

type GenerateRoomResultsRequest struct {
	Job           GenerateRoomResultsRequestJob
	Candidate     GenerateRoomResultsRequestCandidate
	PsyTestResult GenerateRoomResultsRequestPsyTestResult `json:"psy_test_results"`
}

type GenerateRoomResultsRequestJob struct {
	Role        string
	Description string
}

type GenerateRoomResultsRequestCandidate struct {
	Age *int
	Sex crmEntity.CandidateGender
}

type GenerateRoomResultsRequestPsyTestResult struct {
	Traits    []GenerateRoomResultsRequestTrait `json:"traits"`
	SumResult float64                           `json:"sum_result"`
}

type GenerateRoomResultsRequestTrait struct {
	Name           string
	Description    string
	LeftStateName  string
	RightStateName string
	Priority       entity.TraitPriority
	Target         int
	Result         int
}

type GenerateRoomResultsResponse struct {
	Analyze *entity.RoomResultAnalyze
}

type LLMRepository interface {
	GenerateProfileDescriptionByName(ctx context.Context, name string) (string, error)

	GenerateProfileTraitsMapByDescription(
		ctx context.Context,
		req *GenerateProfileTraitsMapByDescriptionRequest,
	) (*GenerateProfileTraitsMapByDescriptionResponse, error)

	GenerateRoomResults(
		ctx context.Context,
		req *GenerateRoomResultsRequest,
	) (*GenerateRoomResultsResponse, error)
}
