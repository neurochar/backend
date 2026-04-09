package server

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
)

type ResponseWriter struct {
	http.ResponseWriter
	Code int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

var GatewayMarshaler = &runtime.JSONPb{
	MarshalOptions: protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: false,
	},
	UnmarshalOptions: protojson.UnmarshalOptions{
		DiscardUnknown: true,
	},
}

func ChainMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func getAllForwardedIPs(r *http.Request) []string {
	var result []string

	vals := r.Header.Values("X-Forwarded-For")

	for _, v := range vals {
		parts := strings.Split(v, ",")
		for _, ip := range parts {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				result = append(result, ip)
			}
		}
	}

	return result
}

func getRequestID(r *http.Request) *uuid.UUID {
	vals := r.Header.Values("X-Request-ID")
	if len(vals) > 0 && vals[0] != "" {
		val, err := uuid.Parse(vals[0])
		if err == nil {
			return &val
		}
	}

	return lo.ToPtr(uuid.New())
}

func getAuthorization(r *http.Request) string {
	vals := r.Header.Values("Authorization")
	if len(vals) > 0 {
		return vals[0]
	}

	return ""
}
