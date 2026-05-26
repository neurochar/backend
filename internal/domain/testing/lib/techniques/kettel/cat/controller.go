package cat

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/neurochar/backend/embed"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques/kettel"
)

type Controller struct {
	grpParams *GRMParams
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) LoadGRParamsFromFile() error {
	loader := embed.NewLoader()
	reader, err := loader.Open("grm_params.json")
	if err != nil {
		return fmt.Errorf("failed to open grm params file: %w", err)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read grm params file: %w", err)
	}

	var params GRMParams
	if err := json.Unmarshal(data, &params); err != nil {
		return fmt.Errorf("failed to unmarshal grm params: %w", err)
	}

	ctrl.grpParams = &params
	return nil
}

type SessionAnswer struct {
	QuestionID uint64
	VariantID  int
}

type PlaySessionResult struct {
	IsFinished   bool
	IsSure       bool
	ResultSten   *int
	NextAnswerID *uint64
}

func (ctrl *Controller) PlaySession(traitID uint64, answers []SessionAnswer) (*PlaySessionResult, error) {
	catTrait, ok := libTraitToCat[traitID]
	if !ok {
		return nil, fmt.Errorf("unknown trait id: %d", traitID)
	}

	if ctrl.grpParams == nil {
		return nil, fmt.Errorf("grm params not loaded")
	}

	factor, ok := ctrl.grpParams.Factors[catTrait]
	if !ok {
		return nil, fmt.Errorf("unknown factor: %s", catTrait)
	}

	const confidenceScore = 0.35

	session := NewCATSession(factor, 10, confidenceScore)
	i := -1
	for !session.IsDone() {
		i++
		item := session.NextItem()
		questionID, ok := catQuestionToLib[item.ID]
		if !ok {
			return nil, fmt.Errorf("unknown question id: %s", item.ID)
		}

		if i > len(answers)-1 {
			return &PlaySessionResult{
				NextAnswerID: &questionID,
			}, nil
		}

		answer := answers[i]
		if answer.QuestionID != questionID {
			return nil, fmt.Errorf("invalid answer id: %d vs %d", answer.QuestionID, questionID)
		}

		libQuestion, ok := kettel.ItemsLib[questionID]
		if !ok {
			return nil, fmt.Errorf("unknown lib question id: %d", questionID)
		}

		if answer.VariantID < 0 || answer.VariantID > len(libQuestion.RawVariantKeys)-1 {
			return nil, fmt.Errorf("invalid variant id: %d", answer.VariantID)
		}
		raw := libQuestion.RawVariantKeys[answer.VariantID]

		session.RecordResponse(item.ID, raw)
	}

	result := session.Result()

	return &PlaySessionResult{
		IsFinished: true,
		IsSure:     result.SE <= confidenceScore,
		ResultSten: &result.Sten,
	}, nil
}
