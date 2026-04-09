package tools

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

type LogRequestData struct {
	RequestID *uuid.UUID `json:"request_id"`
	Processor string     `json:"processor"`
	Method    string     `json:"method"`
	URI       string     `json:"uri,omitempty"`
	Referer   string     `json:"referer,omitempty"`
	IP        string     `json:"ip,omitempty"`
	IPChain   []string   `json:"ip_chain,omitempty"`
}

type LogResponseData struct {
	Processor  string                `json:"processor"`
	DurationMS int64                 `json:"duration_ms"`
	Code       int                   `json:"code"`
	AppError   *appErrors.JSONStruct `json:"app_error,omitempty"`
	Error      string                `json:"error,omitempty"`
}

func LogSetRequest(
	ctx context.Context,
	reqData *LogRequestData,
) context.Context {
	return loghandler.SetContextData(ctx, "request", reqData)
}

func LogSetResponse(
	ctx context.Context,
	resData *LogResponseData,
) context.Context {
	return loghandler.SetContextData(ctx, "response", resData)
}

func LogHTTPContext(
	ctx context.Context,
	logger *slog.Logger,
) {
	data, ok := loghandler.GetData(ctx)
	if !ok {
		return
	}

	resData, ok := data["response"].(*LogResponseData)
	if !ok {
		return
	}

	if resData.Code >= 400 {
		logger.ErrorContext(ctx, "http")
	} else {
		logger.InfoContext(ctx, "http")
	}
}

func LogGRPCContext(
	ctx context.Context,
	logger *slog.Logger,
) {
	data, ok := loghandler.GetData(ctx)
	if !ok {
		return
	}

	resData, ok := data["response"].(*LogResponseData)
	if !ok {
		return
	}

	if resData.Error != "" {
		logger.ErrorContext(ctx, "grpc")
	} else {
		logger.InfoContext(ctx, "grpc")
	}
}
