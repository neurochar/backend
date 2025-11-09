package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/samber/lo"
)

type ListUsersHandlerOut struct {
	Items []OutUser `json:"items"`
	Total uint64    `json:"total"`
}

func (ctrl *Controller) ListUsersHandler(c *fiber.Ctx) error {
	const op = "ListUsersHandler"

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

	roleID := c.QueryInt("role_id", 0)

	query := c.Query("query", "")

	listOptions := &userUC.UserListOptions{}

	if roleID > 0 {
		listOptions.RoleID = lo.ToPtr(uint64(roleID))
	}

	if query != "" {
		listOptions.Query = lo.ToPtr(query)
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	}

	items, total, err := ctrl.userFacade.Common.FindPagedList(c.Context(), listOptions, listParams)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out := ListUsersHandlerOut{
		Items: make([]OutUser, 0, len(items)),
		Total: total,
	}

	for _, item := range items {
		outItem, err := OutUserDTO(c, ctrl.fileUC, item.Account, item.ProfileDTO, item.Role, false)
		if err != nil {
			return err
		}

		out.Items = append(out.Items, *outItem)
	}

	return c.JSON(out)
}
