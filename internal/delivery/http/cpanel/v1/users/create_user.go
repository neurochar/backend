package users

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/neurochar/backend/pkg/validation"

	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
)

type CreateUserHandlerIn struct {
	Account InAccountCreate `json:"account"`
	Profile InProfile       `json:"profile"`
}

func (ctrl *Controller) CreateUserHandler(c *fiber.Ctx) error {
	const op = "CreateProfileHandler"

	in := &CreateUserHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	var photo100x100FileID *uuid.UUID

	if in.Profile.Photo100x100FileID != "" {
		parseID, err := uuid.Parse(in.Profile.Photo100x100FileID)
		if err != nil {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
		}

		photo100x100FileID = &parseID
	}

	requestIP := net.ParseIP(middleware.GetRealIP(c))

	user, err := ctrl.userFacade.Common.CreateUser(
		c.Context(),
		userUC.CreateUserInput{
			Account: userUC.AccountDataInput{
				Email:           in.Account.Email,
				Password:        in.Account.Password,
				RoleID:          in.Account.RoleID,
				IsEmailVerified: in.Account.IsEmailVerified,
			},
			Profile: userUC.ProfileDataInput{
				Name:               in.Profile.Name,
				Surname:            in.Profile.Surname,
				Photo100x100FileID: photo100x100FileID,
			},

			IsSendPassword: in.Account.IsSendPassword,
		},
		requestIP,
		true,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutUserDTO(c, ctrl.fileUC, user.Account, user.ProfileDTO, user.Role, true)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return c.JSON(out)
}
