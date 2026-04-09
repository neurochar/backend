package interceptor

import (
	"context"
	"log/slog"
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"google.golang.org/grpc"
)

func InterceptorLogger(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		reqData := &tools.LogRequestData{
			Processor: "grpc",
			Method:    info.FullMethod,
		}

		enrich := tools.GetEnrich(ctx)
		if enrich != nil {
			if enrich.RequestIP != nil {
				reqData.IP = enrich.RequestIP.String()
			}

			reqData.IPChain = enrich.RequestIPChain

			reqData.RequestID = enrich.RequestID
		}

		ctx = tools.LogSetRequest(ctx, reqData)

		res, err := handler(ctx, req)

		duration := time.Since(start)

		respData := &tools.LogResponseData{
			Processor:  "grpc",
			DurationMS: duration.Milliseconds(),
		}

		if err != nil {
			respData.Error = err.Error()
			errStr := appErrors.ToJSONStruct(err, true, false)
			respData.AppError = &errStr
		}

		ctx = tools.LogSetResponse(ctx, respData)
		tools.LogGRPCContext(ctx, logger)

		return res, err
	}
}
