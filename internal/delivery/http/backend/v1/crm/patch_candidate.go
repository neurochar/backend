package crm

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/null"
	"github.com/neurochar/backend/pkg/validation"

	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
)

type PatchCandidateHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	CandidateName     *string           `json:"candidateName" validate:"required,min=1,max=150"`
	CandidateSurname  *string           `json:"candidateSurname" validate:"required,min=1,max=150"`
	CandidateGender   *uint8            `json:"candidateGender"`
	CandidateBirthday null.NullableTime `json:"candidateBirthday"`
}

func (ctrl *Controller) PatchCandidateHandler(c *fiber.Ctx) error {
	const op = "PatchCandidateHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	in := &PatchCandidateHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	usecaseInput := crmUC.PatchCandidateDataInput{
		Version: in.Version,

		CandidateName:    in.CandidateName,
		CandidateSurname: in.CandidateSurname,
	}

	if in.CandidateGender != nil {
		candidateGender, err := crmEntity.CandidateGenderFromUint8(*in.CandidateGender)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
		}

		usecaseInput.CandidateGender = &candidateGender
	}

	if in.CandidateBirthday.IsSet {
		var value *time.Time
		if in.CandidateBirthday.IsValid {
			value = &in.CandidateBirthday.Time
		}
		usecaseInput.CandidateBirthday = &value
	}

	err = ctrl.crmFacade.Candidate.PatchByDTO(
		c.Context(),
		id,
		usecaseInput,
		in.SkipVersionCheck,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
