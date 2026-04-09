package server

import "google.golang.org/grpc"

type Config struct {
	Addr              string
	UseLogger         bool
	TrustedProxies    []string
	ExtraInterceptors []grpc.UnaryServerInterceptor
}

type Option func(*Config)

func WithLogger(enabled bool) Option {
	return func(c *Config) {
		c.UseLogger = enabled
	}
}

func WithTrustedProxies(proxies []string) Option {
	return func(c *Config) {
		c.TrustedProxies = proxies
	}
}

func WithExtraInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(c *Config) {
		c.ExtraInterceptors = interceptors
	}
}
