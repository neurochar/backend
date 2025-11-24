package fxboot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/app"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/app/fxboot/invoking"
	"github.com/neurochar/backend/internal/app/fxboot/providing"
	cpanelHTTP "github.com/neurochar/backend/internal/delivery/http/cpanel"
	"github.com/neurochar/backend/internal/domain/alert"
	alertUC "github.com/neurochar/backend/internal/domain/alert/usecase"
	"github.com/neurochar/backend/internal/domain/crm"
	emailingModule "github.com/neurochar/backend/internal/domain/emailing"
	"github.com/neurochar/backend/internal/domain/file"
	"github.com/neurochar/backend/internal/domain/tenant"
	"github.com/neurochar/backend/internal/domain/testing"
	"github.com/neurochar/backend/internal/domain/user"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/storage"
	"github.com/neurochar/backend/internal/infra/storage/s3d"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/pgclient"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	storageMigrations "github.com/neurochar/backend/internal/infra/storage/migrations"
)

func CPanelAppGetOptionsMap(appID app.ID, cfg config.Config) OptionsMap {
	return OptionsMap{
		Providing: map[ProvidingID]fx.Option{
			ProvidingAppID: fx.Provide(func() app.ID {
				return appID
			}),
			ProvidingIDFXTimeouts: fx.Options(
				fx.StartTimeout(time.Second*time.Duration(cfg.CPanelApp.Base.StartTimeoutSec)),
				fx.StopTimeout(time.Second*time.Duration(cfg.CPanelApp.Base.StopTimeoutSec)),
			),
			ProvidingIDConfig: fx.Provide(func() config.Config {
				return cfg
			}),
			ProvidingIDLogger: fx.Provide(func(cfg config.Config) *slog.Logger {
				return providing.NewLogger(
					cfg.CPanelApp.Name,
					cfg.CPanelApp.Version,
					cfg.CPanelApp.Base.UseLogger,
					cfg.CPanelApp.Base.IsProd,
				)
			}),
			ProvidingIDFXLogger: fx.WithLogger(func(cfg config.Config) fxevent.Logger {
				return providing.NewFXLogger(cfg.CPanelApp.Base.UseFxLogger)
			}),
			ProvidingIDImageProc: fx.Provide(providing.NewImageProc),
			ProvidingIDDBClients: fx.Provide(
				func(logger *slog.Logger, cfg config.Config, shutdown fx.Shutdowner) db.MasterClient {
					return providing.NewDBClients(
						cfg.Postgres.Master.DSN,
						cfg.CPanelApp.Base.LogSQLQueries,
						logger,
						shutdown,
					)
				},
			),
			ProvidingIDBackoff:       fx.Provide(providing.NewBackoff),
			ProvidingIDStorageClient: fx.Provide(providing.NewStorageClient),
			ProvidingIDEmailing:      fx.Provide(providing.NewEmailing),
			ProvidingHTTPFiberServer: fx.Provide(
				func(logger *slog.Logger, cfg config.Config, alertUsecase alertUC.Usecase) *fiber.App {
					httpConfig := cpanelHTTP.HTTPConfig{
						AppTitle:         cfg.Global.ProjectName,
						UnderProxy:       cfg.CPanelApp.HTTP.UnderProxy,
						UseLogger:        cfg.CPanelApp.Base.UseLogger && cfg.CPanelApp.Base.LogHTTP,
						BodyLimit:        -1,
						CorsAllowOrigins: cfg.CPanelApp.HTTP.CorsAllowOrigins,
						ServerIPs:        []string{cfg.Global.ServerIP},
					}

					return cpanelHTTP.NewHTTPFiber(httpConfig, logger, alertUsecase)
				},
			),
			ProvidingIDDeliveryHTTP:   cpanelHTTP.FxModule,
			ProvidingIDFileModule:     file.FxModule,
			ProvidingIDUserModule:     user.FxModule,
			ProvidingIDEmailingModule: emailingModule.FxModule,
			ProvidingIDAlertModule:    alert.FxModule,
			ProvidingIDTenantModule:   tenant.FxModule,
			ProvidingIDCRMModule:      crm.FxModule,
			ProvidingIDTestingModule:  testing.FxModule,
		},
		Invokes: []fx.Option{
			fx.Invoke(CPanelAppInitInvoke),
		},
	}
}

type CPanelInvokeInput struct {
	fx.In

	LC              fx.Lifecycle
	Shutdowner      fx.Shutdowner
	Invokes         []invoking.InvokeInit `group:"InvokeInit"`
	Logger          *slog.Logger
	Cfg             config.Config
	DBMasterClient  db.MasterClient
	BackoffCtrl     *backoff.Controller
	S3Client        *s3.Client
	StorageClient   storage.Client
	HttpFiberServer *fiber.App
}

// CPanelAppInitInvoke - app init
func CPanelAppInitInvoke(
	in CPanelInvokeInput,
) {
	in.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Тестирование соединения с мастером postgress
			err := pgclient.TestConnection(
				ctx,
				in.DBMasterClient,
				in.Logger,
				in.Cfg.Postgres.MaxAttempts,
				in.Cfg.Postgres.AttemptSleepSeconds,
			)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to test master db connection", slog.Any("error", err))
				return err
			}

			in.Logger.InfoContext(
				ctx,
				"successfully connected to Postgress",
				slog.String("serverID", in.DBMasterClient.ServerID()),
			)

			// Миграции goose
			err = db.UpMigrations(in.Cfg.Postgres.Master.DSN, in.Cfg.Postgres.MigrationsPath, in.Logger)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to run migrations", slog.Any("error", err))
				return err
			}

			// Тестирование соединения с s3
			err = s3d.PingS3Client(ctx, in.S3Client)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to ping s3", slog.Any("error", err))
				return err
			}
			in.Logger.InfoContext(ctx, "connected to s3")

			// Миграции хранилища
			if in.Cfg.Storage.UpMigrations {
				createdAny, err := storageMigrations.UpBuckets(ctx, in.StorageClient)
				if err != nil {
					in.Logger.ErrorContext(ctx, "failed to migrate storage", slog.Any("error", err))
					return err
				}

				if createdAny {
					in.Logger.InfoContext(ctx, "storage buckets created")
				} else {
					in.Logger.InfoContext(ctx, "storage buckets already exist")
				}
			} else {
				in.Logger.InfoContext(ctx, "storage migrations skipped")
			}

			// Запускаем invoke функции до открытия
			for _, invokeItem := range in.Invokes {
				if invokeItem.StartBeforeOpen != nil {
					err := invokeItem.StartBeforeOpen(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn start before open", slog.Any("error", err))
						return err
					}
				}
			}

			// Запускаем http
			if in.Cfg.CPanelApp.HTTP.Port > 0 {
				in.Logger.InfoContext(ctx, "starting http server", slog.Int("port", in.Cfg.CPanelApp.HTTP.Port))
				go func() {
					if err := in.HttpFiberServer.Listen(fmt.Sprintf(":%d", in.Cfg.CPanelApp.HTTP.Port)); err != nil {
						in.Logger.ErrorContext(ctx, "failed to start fiber", slog.Any("error", err))
						err := in.Shutdowner.Shutdown()
						if err != nil {
							in.Logger.ErrorContext(ctx, "failed to shutdown", slog.Any("error", err))
						}
					}
				}()
			}

			// Запускаем invoke функции после открытия
			for _, invokeItem := range in.Invokes {
				if invokeItem.StartAfterOpen != nil {
					err := invokeItem.StartAfterOpen(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn start after open", slog.Any("error", err))
						return err
					}
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, invokeItem := range in.Invokes {
				if invokeItem.Stop != nil {
					err := invokeItem.Stop(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn stop", slog.Any("error", err))
						return err
					}
				}
			}

			// Останавливаем http
			if in.Cfg.CPanelApp.HTTP.Port > 0 {
				in.Logger.Info("stopping http fiber")
				err := in.HttpFiberServer.ShutdownWithTimeout(time.Duration(in.Cfg.CPanelApp.HTTP.StopTimeoutSec) * time.Second)
				if err != nil {
					in.Logger.ErrorContext(ctx, "failed to stop fiber", slog.Any("error", err))
				}
			}

			// Закрываем postgress
			in.DBMasterClient.Close()
			in.Logger.InfoContext(ctx, "closing db clients")

			// Останавливаем backoff
			in.BackoffCtrl.Stop(ctx)

			return nil
		},
	})
}
