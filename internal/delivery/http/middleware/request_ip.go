package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

// RequestIP - middleware for request ip
func RequestIP() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestIP := GetRealIP(c)

		c.Locals("requestIP", requestIP)

		ctxData, ctxKey := loghandler.SetData(c.Context(), "request.ip", requestIP)
		c.Locals(ctxKey, ctxData)

		return c.Next()
	}
}

func GetRealIP(c *fiber.Ctx) string {
	ips := strings.Split(c.IP(), ",")
	for i := range ips {
		ips[i] = strings.TrimSpace(ips[i])
	}

	forwardedBy := c.Get("X-Forwarded-By")

	for i := len(ips) - 1; i >= 0; i-- {
		ip := ips[i]
		parsed := net.ParseIP(ip)
		if parsed == nil {
			continue
		}
		if !IsPrivateIP(parsed) && (forwardedBy == "" || !strings.Contains(forwardedBy, ip)) {
			return ip
		}
	}

	if len(ips) > 0 {
		return ips[len(ips)-1]
	}

	return c.IP()
}

func IsPrivateIP(ip net.IP) bool {
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
	}

	for _, block := range privateBlocks {
		_, subnet, err := net.ParseCIDR(block)
		if err != nil {
			continue
		}
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}
