package server

import (
	"time"
)

type Config struct {
	Addr             string
	UseLogger        bool
	CorsAllowOrigins []string
	TrustedProxies   []string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
}

type Option func(*Config)

func WithReadTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.ReadTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.WriteTimeout = d
	}
}

func WithIdleTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.IdleTimeout = d
	}
}

func WithLogger(enabled bool) Option {
	return func(c *Config) {
		c.UseLogger = enabled
	}
}

func WithCORS(origins []string) Option {
	return func(c *Config) {
		c.CorsAllowOrigins = origins
	}
}

func WithTrustedProxies(proxies []string) Option {
	return func(c *Config) {
		c.TrustedProxies = proxies
	}
}
