package crm

import (
	"time"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type CreateCandidateHandlerIn struct {
	CandidateName     string     `json:"candidateName" validate:"required,min=1,max=150"`
	CandidateSurname  string     `json:"candidateSurname" validate:"required,min=1,max=150"`
	CandidateGender   uint8      `json:"candidateGender"`
	CandidateBirthday *time.Time `json:"candidateBirthday"`
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

	candidateGender, err := crmEntity.CandidateGenderFromUint8(in.CandidateGender)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	authData := backendMiddleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	candidateDTO, err := ctrl.crmFacade.Candidate.CreateByDTO(
		c.Context(),
		authData.TenantID,
		crmUC.CreateCandidateDataInput{
			CandidateName:     in.CandidateName,
			CandidateSurname:  in.CandidateSurname,
			CandidateGender:   candidateGender,
			CandidateBirthday: in.CandidateBirthday,
			CreatedBy:         &authData.AccountID,
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
