package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/dto"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

type PatchAccountHandlerIn struct {
	dto.UpdateMeta

	Email           *string `json:"email" validate:"omitempty,email"`
	Password        *string `json:"password" validate:"omitempty"`
	RoleID          *uint64 `json:"roleID" validate:"omitempty"`
	IsEmailVerified *bool   `json:"isEmailVerified" validate:"omitempty"`
	IsBlocked       *bool   `json:"isBlocked" validate:"omitempty"`
}

func (ctrl *Controller) PatchAccountHandler(c *fiber.Ctx) error {
	const op = "PatchAccountHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	in := &PatchAccountHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	removeSessions := in.Password != nil && *in.Password != ""

	err = ctrl.userFacade.Common.PatchAccountByDTO(
		c.Context(),
		id,
		userUC.PatchAccountDataInput{
			Version: in.Version,

			Email:           in.Email,
			Password:        in.Password,
			RoleID:          in.RoleID,
			IsEmailVerified: in.IsEmailVerified,
			IsBlocked:       in.IsBlocked,
		},
		removeSessions,
		in.SkipVersionCheck,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
