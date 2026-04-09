package users_tenant

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) ListAccounts(
	ctx context.Context,
	req *desc.ListAccountsRequest,
) (*desc.ListAccountsResponse, error) {
	const op = "ListAccounts"

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

	listOptions := &tenantUC.AccountListOptions{
		FilterTenantID: &authData.TenantUserClaims().TenantID,
		Sort: []uctypes.SortOption[tenantUC.AccountListOptionsSortField]{
			{
				Field:  tenantUC.AccountListOptionsSortFieldCreatedAt,
				IsDesc: false,
			},
		},
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  limit,
		Offset: offset,
	}

	items, total, err := ctrl.tenantFacade.Account.FindPagedList(
		ctx,
		listOptions,
		listParams,
		&tenantUC.AccountDTOOptions{
			FetchPhotoFiles: true,
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.ListAccountsResponse{
		Items: lo.Map(items, func(item *tenantUC.AccountDTO, _ int) *typesv1.AccountTenant {
			return mapper.TenantAccountToPb(item, ctrl.fileUC, false)
		}),
		Total: int32(total),
	}, nil
}
