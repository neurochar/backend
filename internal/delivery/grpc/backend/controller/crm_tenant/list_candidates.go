package crm_tenant

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/delivery/grpc/backend/mapper"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	"github.com/samber/lo"
)

func (ctrl *controller) ListCandidates(
	ctx context.Context,
	req *desc.ListCandidatesRequest,
) (*desc.ListCandidatesResponse, error) {
	const op = "ListCandidates"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	} else if limit < 1 {
		limit = 1
	}

	offset := req.GetOffset()

	listOptions := &crmUC.CandidateListOptions{
		FilterTenantID: &authData.TenantUserClaims().TenantID,
		Sort: []uctypes.SortOption[crmUC.CandidateListOptionsSortField]{
			{
				Field:  crmUC.CandidateListOptionsSortFieldCreatedAt,
				IsDesc: true,
			},
		},
	}

	if req.SearchQuery != nil {
		listOptions.SearchQuery = req.SearchQuery
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  limit,
		Offset: offset,
	}

	items, total, err := ctrl.crmFacade.Candidate.FindPagedList(
		ctx,
		listOptions,
		listParams,
		&crmUC.CandidateDTOOptions{},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.ListCandidatesResponse{
		Items: lo.Map(items, func(item *crmUC.CandidateDTO, _ int) *desc.Candidate {
			return mapper.CandidateDTOToPb(item)
		}),
		Total: total,
	}, nil
}
