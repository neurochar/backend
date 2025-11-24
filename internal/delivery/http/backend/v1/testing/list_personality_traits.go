package testing

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

type ListPersonalityTraitsHandlerOut struct {
	Items []entity.PersonalityTrait `json:"items"`
}

func (ctrl *Controller) ListPersonalityTraitsHandler(c *fiber.Ctx) error {
	const op = "ListPersonalityTraitsHandler"

	items, err := ctrl.testingFacade.PersonalityTrait.FindList(
		c.Context(),
		nil,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out := ListPersonalityTraitsHandlerOut{
		Items: make([]entity.PersonalityTrait, 0, len(items)),
	}

	// nolint
	for _, item := range items {
		out.Items = append(out.Items, item)
	}

	return c.JSON(out)
}
