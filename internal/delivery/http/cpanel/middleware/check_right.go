package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	userConstants "github.com/neurochar/backend/internal/domain/user/constants"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

func (ctrl *Controller) MiddlewareCheckRight(rightKey string, shouldBeVal int) func(*fiber.Ctx) error {
	var right *userEntity.Right
	for _, r := range userConstants.Rights {
		if r.Key == rightKey {
			right = r
			break
		}
	}
	if right == nil {
		panic(fmt.Sprintf("right %s not found", rightKey))
	}

	return func(c *fiber.Ctx) error {
		authData := GetAuthData(c)
		if authData == nil {
			return appErrors.ErrUnauthorized
		}

		roleRight, err := authData.Role.GetRightByKey(rightKey)
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		if roleRight.Value != shouldBeVal {
			return appErrors.ErrForbidden
		}

		return c.Next()
	}
}
