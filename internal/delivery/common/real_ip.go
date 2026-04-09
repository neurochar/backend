package common

import (
	"fmt"
	"net"

	"github.com/samber/lo"
)

func ParseRealIP(ips []string, serverIPs []string) (string, error) {
	for i := len(ips) - 1; i >= 0; i-- {
		ip := ips[i]
		parsed := net.ParseIP(ip)
		if parsed == nil {
			continue
		}
		if !lo.Contains(serverIPs, ip) && !IsPrivateIP(parsed) {
			return ip, nil
		}
	}

	if len(ips) > 0 {
		return ips[len(ips)-1], nil
	}

	return "", fmt.Errorf("empty ip list")
}

var privateBlocks = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.0/8",
	"::1/128",
}

func IsPrivateIP(ip net.IP) bool {
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
