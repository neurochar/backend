package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/pkg/validation"
	"github.com/samber/lo"
)

type UpdateMyPasswordHandlerIn struct {
	CurrentPassword string `json:"currentPassword" validate:""`
	NewPassword     string `json:"newPassword" validate:""`
	NewPassword2    string `json:"newPassword2" validate:""`
}

func (ctrl *Controller) UpdateMyPasswordHandler(c *fiber.Ctx) error {
	const op = "UpdateMyPasswordHandler"

	in := &UpdateMyPasswordHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	authData := middleware.GetAuthData(c)

	if authData == nil {
		return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", ctrl.pkg, op)
	}

	if !authData.Account.VerifyPassword(in.CurrentPassword) {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithTextCode("CURRENT_PASSWORD_INCORRECT").WithHints("current password is incorrect"),
			"%s.%s", ctrl.pkg, op)
	}

	if in.NewPassword != in.NewPassword2 {
		return appErrors.Chainf(ErrPasswordsMismatch, "%s.%s", ctrl.pkg, op)
	}

	err := ctrl.userFacade.Common.PatchAccountByDTO(
		c.Context(),
		authData.Account.ID,
		userUC.PatchAccountDataInput{
			Password: lo.ToPtr(in.NewPassword),
		},
		true,
		true,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
