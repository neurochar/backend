package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

type ListAccountsHandlerOut struct {
	Items []OutAccount `json:"items"`
	Total uint64       `json:"total"`
}

func (ctrl *Controller) ListAccountsHandler(c *fiber.Ctx) error {
	const op = "ListAccountsHandler"

	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 1
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	auth := middleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	tenantID, err := uuid.Parse(auth.TenantID.String())
	if err != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	listOptions := &tenantUserUC.AccountListOptions{
		FilterTenantID: &tenantID,
		Sort: []uctypes.SortOption[tenantUserUC.AccountListOptionsSortField]{
			{
				Field:  tenantUserUC.AccountListOptionsSortFieldCreatedAt,
				IsDesc: false,
			},
		},
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}

	items, total, err := ctrl.tenantUserFacade.Account.FindPagedList(
		c.Context(),
		listOptions,
		listParams,
		&tenantUserUC.AccountDTOOptions{
			FetchPhotoFiles: true,
		},
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out := ListAccountsHandlerOut{
		Items: make([]OutAccount, 0, len(items)),
		Total: total,
	}

	for _, item := range items {
		outItem, err := OutAccountDTO(c, false, ctrl.fileUC, item)
		if err != nil {
			return err
		}

		out.Items = append(out.Items, *outItem)
	}

	return c.JSON(out)
}
