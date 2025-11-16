package auth

import (
	"errors"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/neurochar/backend/pkg/validation"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

type RefreshHandlerIn struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (ctrl *Controller) RefreshHandler(c *fiber.Ctx) error {
	const op = "RefreshHandler"

	in := &RefreshHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	ip := middleware.GetRealIP(c)

	requestIP := net.ParseIP(ip)

	claims, err := ctrl.tenantUserFacade.Auth.ParseRefreshToken(in.RefreshToken, true)
	if err != nil {
		if errors.Is(err, tenantUserUC.ErrInvalidToken) {
			return appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
		}
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	authDTO, err := ctrl.tenantUserFacade.Auth.GenerateNewClaims(c.Context(), claims, requestIP)
	if err != nil {
		if errors.Is(err, appErrors.ErrUnauthorized) {
			return appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
		}
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	accessJWT, err := ctrl.tenantUserFacade.Auth.IssueAccessJWT(authDTO.AccessClaims)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	refreshJWT, err := ctrl.tenantUserFacade.Auth.IssueRefreshJWT(authDTO.RefreshClaims)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutTokensDTO(c, refreshJWT, uint64(ctrl.cfg.Auth.RefreshTokenLifetimeHrs)*3600, accessJWT)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
