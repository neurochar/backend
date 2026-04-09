package registration

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/registration/v1"
)

func (ctrl *Controller) StartRegistration(
	ctx context.Context,
	req *desc.StartRegistrationRequest,
) (*desc.StartRegistrationResponse, error) {
	const op = "StartRegistration"

	ip := tools.GetRealIP(ctx)

	if ip != nil {
		backoffSession, ok := ctrl.backoff.GetIfExists(ip.String(), backoffConfigRegistrationGroupID)

		if ok && !backoffSession.IsAllowed() {
			tryAfter := backoffSession.NextAllowedUntilSeconds()

			return nil, appErrors.Chainf(
				appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
				"%s.%s", ctrl.pkg, op,
			)
		}
	}

	_, err := ctrl.tenantFacade.Registration.CreateByDTO(
		ctx,
		tenantUC.CreateRegistrationIn{
			Email:     req.Email,
			RequestIP: ip,
		},
		ip,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	backoffSession := ctrl.backoff.GetOrCreate(ip.String(), backoffConfigRegistrationGroupID)
	backoffSession.AddBackoff()

	return &desc.StartRegistrationResponse{}, nil
}
