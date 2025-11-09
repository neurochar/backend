package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Cors - middleware for cors
func Cors(corsAllowOrigins []string) func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(corsAllowOrigins, ", "),
		AllowCredentials: true,
	})
}
