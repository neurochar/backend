package registration

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/registration/v1"
)

func (ctrl *Controller) CheckRegistration(
	ctx context.Context,
	req *desc.CheckRegistrationRequest,
) (*desc.CheckRegistrationResponse, error) {
	const op = "CheckRegistration"

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	registration, err := ctrl.tenantFacade.Registration.FindOneByID(ctx, id, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if registration.IsFinished {
		return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", ctrl.pkg, op)
	}

	return &desc.CheckRegistrationResponse{}, nil
}
