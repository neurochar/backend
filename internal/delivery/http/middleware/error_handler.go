// Package middleware contains middleware for http handlers
package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	alertUC "github.com/neurochar/backend/internal/domain/alert/usecase"
)

type errorJSON struct {
	Code     int            `json:"code"`
	TextCode string         `json:"textCode"`
	Hints    []string       `json:"hints"`
	Details  map[string]any `json:"details"`
}

// ErrorHandler - обработчик ошибок
func ErrorHandler(appTitle string, logger *slog.Logger, alertUsecase alertUC.Usecase) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		code := 500
		jsonRes := errorJSON{
			TextCode: "INTERNAL_ERROR",
			Hints:    []string{},
			Details:  map[string]any{},
		}

		if appError, ok := appErrors.ExtractError(err); ok {
			code = int(appError.Meta().Code)
			jsonRes.TextCode = appError.Meta().TextCode
			jsonRes.Hints = appError.Hints()
			jsonRes.Details = appError.Details(false)
		} else {
			switch errTyped := err.(type) {
			case *fiber.Error:
				code = errTyped.Code
				switch {
				case code == 405:
					jsonRes.TextCode = "METHOD_NOT_ALLOWED"
				case code >= 400 && code < 500:
					jsonRes.TextCode = "BAD_REQUEST"
				}
				jsonRes.Hints = []string{errTyped.Message}
			default:
			}
		}

		jsonRes.Code = code

		if code >= 500 {
			go func() {
				/*
					defer func() {
						if r := recover(); r != nil {
							logger.ErrorContext(
								loghandler.WithSource(c.Context()),
								"panic inside http error handler",
								slog.Any("error", r),
								slog.Any("trackeback", string(debug.Stack())),
							)
						}
					}()

					uri := ""
					if c.Context() != nil && c.Context().URI() != nil {
						uri = string(c.Context().URI().QueryString())
					}

					message := fmt.Sprintf("%s\n\n%s: %s?%s", appTitle, c.Method(), c.Path(), uri)

					appErr, ok := appErrors.ExtractError(err)
					if ok {
						structErr := appErrors.ToJSONStruct(appErr, true, true)

						b, err := json.Marshal(structErr)
						if err != nil {
							b = []byte(err.Error())
						}

						message += fmt.Sprintf("\n\nError: %s", string(b))
					} else {
						if err != nil {
							message += fmt.Sprintf("\n\nError: %s", err.Error())
						}
					}

					if c.Body() != nil {
						message += fmt.Sprintf("\n\nBody: %s", string(c.Body()))
					}

					err = alertUsecase.SendAlert(c.Context(), message)
					if err != nil {
						logger.Error("failed to send alert", slog.Any("error", err))
					}
				*/
			}()
		}

		return c.Status(code).JSON(jsonRes)
	}
}
