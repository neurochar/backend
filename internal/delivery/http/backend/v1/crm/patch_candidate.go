package crm

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
)

type PatchCandidateHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	CandidateName    *string `json:"candidateName" validate:"required,min=1,max=150"`
	CandidateSurname *string `json:"candidateSurname" validate:"required,min=1,max=150"`
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
