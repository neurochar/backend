package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/delivery/common"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

// RequestIP - middleware for request ip
func RequestIP(serverIPs []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestIP := parseRealIP(c, serverIPs)

		c.Locals("requestIP", requestIP)

		ctxData, ctxKey := loghandler.SetData(c.Context(), "request.ip", requestIP)
		c.Locals(ctxKey, ctxData)

		return c.Next()
	}
}

func parseRealIP(c *fiber.Ctx, serverIPs []string) string {
	ips := strings.Split(c.IP(), ",")
	clientIPs := make([]string, 0, len(ips))
	for i := range ips {
		v := strings.TrimSpace(ips[i])
		if v != "" {
			clientIPs = append(clientIPs, v)
		}
	}

	realIP, err := common.ParseRealIP(clientIPs, serverIPs)
	if err != nil {
		return c.Context().RemoteIP().String()
	}

	return realIP
}

func GetRealIP(c *fiber.Ctx) string {
	ip := c.Locals("requestIP")
	ipStr, ok := ip.(string)
	if ok {
		return ipStr
	}

	return c.IP()
}
