package crm

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type CreateCandidateHandlerIn struct {
	CandidateName    string `json:"candidateName" validate:"required,min=1,max=150"`
	CandidateSurname string `json:"candidateSurname" validate:"required,min=1,max=150"`
}

func (ctrl *Controller) CreateCandidateHandler(c *fiber.Ctx) error {
	const op = "CreateCandidateHandler"

	in := &CreateCandidateHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	authData := backendMiddleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	candidateDTO, err := ctrl.crmFacade.Candidate.CreateByDTO(
		c.Context(),
		authData.TenantID,
		crmUC.CreateCandidateDataInput{
			CandidateName:    in.CandidateName,
			CandidateSurname: in.CandidateSurname,
			CreatedBy:        &authData.AccountID,
		},
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutCandidateDTO(
		c,
		candidateDTO,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
