package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthCtxKey string

const (
	AuthCtxKeyToken AuthCtxKey = "auth_token"
)

func Authorization() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		if accessToken == "" {
			return c.Next()
		}

		c.Locals(AuthCtxKeyToken, accessToken)

		return c.Next()
	}
}
