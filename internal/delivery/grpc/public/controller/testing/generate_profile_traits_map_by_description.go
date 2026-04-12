package testing

import (
	"context"
	"errors"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
)

func (ctrl *Controller) GenerateProfileTraitsMapByDescription(
	ctx context.Context,
	req *desc.GenerateProfileTraitsMapByDescriptionRequest,
) (*desc.GenerateProfileTraitsMapByDescriptionResponse, error) {
	const op = "GenerateProfileTraitsMapByDescription"

	ip := tools.GetRealIP(ctx)

	err := ctrl.limiter.Get(limiter.DefaultName).Register(ctx, &limiter.RegisterKey{
		IP: ip,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if ip != nil {
		backoffSession, ok := ctrl.backoff.GetIfExists(ip.String(), backoffConfigLLMGroupID)

		if ok && !backoffSession.IsAllowed() {
			tryAfter := backoffSession.NextAllowedUntilSeconds()

			return nil, appErrors.Chainf(
				appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
				"%s.%s", ctrl.pkg, op,
			)
		}
	}

	resp, err := ctrl.testingFacade.Profile.GenerateProfileTraitsMapByDescription(ctx, &testingUC.GenerateProfileTraitsMapByDescriptionRequest{
		Description: req.Description,
		Role:        req.Name,
	})
	if err != nil {
		if !errors.Is(err, appErrors.ErrInternal) {
			backoffSession := ctrl.backoff.GetOrCreate(ip.String(), backoffConfigLLMGroupID)
			backoffSession.AddBackoff()
		}

		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	backoffSession := ctrl.backoff.GetOrCreate(ip.String(), backoffConfigLLMGroupID)
	backoffSession.AddBackoff()

	return &desc.GenerateProfileTraitsMapByDescriptionResponse{
		Traits: mapper.ProfilePersonalityTraitsMapToPb(resp.TraitsMap),
	}, nil
}
