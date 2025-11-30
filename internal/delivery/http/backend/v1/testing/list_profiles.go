package testing

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
)

type ListProfilesHandlerOut struct {
	Items []OutProfile `json:"items"`
	Total uint64       `json:"total"`
}

func (ctrl *Controller) ListProfilesHandler(c *fiber.Ctx) error {
	const op = "ListProfilesHandler"

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

	search := c.Query("search")

	authData := middleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	listOptions := &testingUC.ProfileListOptions{
		FilterTenantID: &authData.TenantID,
		Sort: []uctypes.SortOption[testingUC.ProfileListOptionsSortField]{
			{
				Field:  testingUC.ProfileListOptionsSortFieldCreatedAt,
				IsDesc: true,
			},
		},
	}

	if search != "" {
		listOptions.SearchQuery = &search
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}

	items, total, err := ctrl.testingFacade.Profile.FindPagedList(
		c.Context(),
		listOptions,
		listParams,
		&testingUC.ProfileDTOOptions{},
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out := ListProfilesHandlerOut{
		Items: make([]OutProfile, 0, len(items)),
		Total: total,
	}

	for _, item := range items {
		outItem, err := OutProfileDTO(c, item)
		if err != nil {
			return err
		}

		out.Items = append(out.Items, *outItem)
	}

	return c.JSON(out)
}
