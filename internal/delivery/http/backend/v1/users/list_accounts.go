package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
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

	authData := middleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	listOptions := &tenantUC.AccountListOptions{
		FilterTenantID: &authData.TenantID,
		Sort: []uctypes.SortOption[tenantUC.AccountListOptionsSortField]{
			{
				Field:  tenantUC.AccountListOptionsSortFieldCreatedAt,
				IsDesc: false,
			},
		},
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}

	items, total, err := ctrl.tenantFacade.Account.FindPagedList(
		c.Context(),
		listOptions,
		listParams,
		&tenantUC.AccountDTOOptions{
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
