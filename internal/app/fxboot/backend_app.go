package fxboot

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neurochar/backend/internal/app"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/app/fxboot/invoking"
	"github.com/neurochar/backend/internal/app/fxboot/providing"
	deliveryCommon "github.com/neurochar/backend/internal/delivery/common"
	publicGRPC "github.com/neurochar/backend/internal/delivery/grpc/public"
	publicHTTPServer "github.com/neurochar/backend/internal/delivery/httpgw/server"
	"github.com/neurochar/backend/internal/domain/alert"
	"github.com/neurochar/backend/internal/domain/crm"
	emailingModule "github.com/neurochar/backend/internal/domain/emailing"
	"github.com/neurochar/backend/internal/domain/file"
	"github.com/neurochar/backend/internal/domain/tenant"
	"github.com/neurochar/backend/internal/domain/testing"
	"github.com/neurochar/backend/internal/domain/user"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/storage"
	"github.com/neurochar/backend/internal/infra/storage/s3d"
	"github.com/neurochar/backend/internal/jobs"
	"github.com/neurochar/backend/pkg/backoff"
	"github.com/neurochar/backend/pkg/pgclient"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	storageMigrations "github.com/neurochar/backend/internal/infra/storage/migrations"
)

func BackendAppGetOptionsMap(appID app.ID, cfg config.Config) OptionsMap {
	return OptionsMap{
		Providing: map[ProvidingID]fx.Option{
			ProvidingAppID: fx.Provide(func() app.ID {
				return appID
			}),
			ProvidingIDFXTimeouts: fx.Options(
				fx.StartTimeout(time.Second*time.Duration(cfg.BackendApp.Base.StartTimeoutSec)),
				fx.StopTimeout(time.Second*time.Duration(cfg.BackendApp.Base.StopTimeoutSec)),
			),
			ProvidingIDConfig: fx.Provide(func() config.Config {
				return cfg
			}),
			ProvidingIDLogger: fx.Provide(func(cfg config.Config) *slog.Logger {
				return providing.NewLogger(
					cfg.BackendApp.Name,
					cfg.BackendApp.Version,
					cfg.BackendApp.Base.UseLogger,
					cfg.BackendApp.Base.IsProd,
				)
			}),
			ProvidingIDFXLogger: fx.WithLogger(func(cfg config.Config) fxevent.Logger {
				return providing.NewFXLogger(cfg.BackendApp.Base.UseFxLogger)
			}),
			ProvidingIDImageProc: fx.Provide(providing.NewImageProc),
			ProvidingIDDBClients: fx.Provide(
				func(logger *slog.Logger, cfg config.Config, shutdown fx.Shutdowner) db.MasterClient {
					return providing.NewDBClients(
						cfg.Postgres.Master.DSN,
						cfg.BackendApp.Base.LogSQLQueries,
						logger,
						shutdown,
					)
				},
			),
			ProvidingIDOpenAIClient:    fx.Provide(providing.NewOpenAIClient),
			ProvidingIDBackoff:         fx.Provide(providing.NewBackoff),
			ProvidingIDStorageClient:   fx.Provide(providing.NewStorageClient),
			ProvidingIDEmailing:        fx.Provide(providing.NewEmailing),
			ProvidingIDDeliveryCommon:  deliveryCommon.FxModule,
			ProvidingPublicGRPCServer:  providing.PublicGRPCServer,
			ProvidingPublicHTTPGateway: providing.PublicHTTPGateway,
			ProvidingIDJobsController:  fx.Provide(jobs.NewController),
			ProvidingIDFileModule:      file.FxModule,
			ProvidingIDUserModule:      user.FxModule,
			ProvidingIDEmailingModule:  emailingModule.FxModule,
			ProvidingIDAlertModule:     alert.FxModule,
			ProvidingIDTenantModule:    tenant.FxModule,
			ProvidingIDCRMModule:       crm.FxModule,
			ProvidingIDTestingModule:   testing.FxModule,
		},
		Invokes: []fx.Option{
			fx.Invoke(BackendAppInitInvoke),
		},
	}
}

type BackendInvokeInput struct {
	fx.In

	LC               fx.Lifecycle
	Shutdowner       fx.Shutdowner
	Invokes          []invoking.InvokeInit `group:"InvokeInit"`
	Logger           *slog.Logger
	Cfg              config.Config
	DBMasterClient   db.MasterClient
	BackoffCtrl      *backoff.Controller
	S3Client         *s3.Client
	StorageClient    storage.Client
	PublicGRPCServer *publicGRPC.PublicServer
	PublicHTTPServer *publicHTTPServer.Server
	JobsController   *jobs.Controller
}

// BackendAppInitInvoke - app init
func BackendAppInitInvoke(
	in BackendInvokeInput,
) {
	ctxWork, cancel := context.WithCancel(context.Background())

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

			// Миграции хранилища
			if in.Cfg.Storage.UpMigrations {
				// Тестирование соединения с s3
				err = s3d.PingS3Client(ctx, in.S3Client)
				if err != nil {
					in.Logger.ErrorContext(ctx, "failed to ping s3", slog.Any("error", err))
					return err
				}
				in.Logger.InfoContext(ctx, "connected to s3")

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

			// Регистрация кронов
			in.JobsController.RegisterProcessFilesToDelete(
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessFilesToDelete.TimeoutMillisec)*time.Millisecond,
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessFilesToDelete.FailedTimeoutMillisec)*time.Millisecond,
			)

			in.JobsController.RegisterProcessUnusedFiles(
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessUnusedFiles.TimeoutMillisec)*time.Millisecond,
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessUnusedFiles.FailedTimeoutMillisec)*time.Millisecond,
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessUnusedFiles.UnusedTtlMin)*time.Minute,
			)

			in.JobsController.RegisterProcessEmailsToSend(
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessEmailsToSend.TimeoutMillisec)*time.Millisecond,
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessEmailsToSend.FailedTimeoutMillisec)*time.Millisecond,
			)

			in.JobsController.RegisterProcessEmailsToDelete(
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessEmailsToDelete.TimeoutMillisec)*time.Millisecond,
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessEmailsToDelete.FailedTimeoutMillisec)*time.Millisecond,
				time.Duration(in.Cfg.CronjobApp.Jobs.ProcessEmailsToDelete.TtlMin)*time.Minute,
			)

			// Запускаем invoke функции до открытия
			for _, invokeItem := range in.Invokes {
				if invokeItem.StartBeforeOpen != nil {
					err := invokeItem.StartBeforeOpen(ctxWork)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn start before open", slog.Any("error", err))
						return err
					}
				}
			}

			if in.Cfg.CronjobApp.Jobs.Autostart {
				in.JobsController.StartAll()
				in.Logger.InfoContext(ctx, "jobs started")
			} else {
				in.Logger.InfoContext(ctx, "jobs autostart skipped")
			}

			// Запускаем public gRPC
			servePublicGRPC, err := in.PublicGRPCServer.Server().Listen()
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to start gRPC public server", slog.Any("error", err))
			}
			in.Logger.InfoContext(ctx, "started gRPC public server", slog.Int("port", in.Cfg.BackendApp.GRPC.Port))

			go func() {
				if err := servePublicGRPC(); err != nil {
					in.Logger.ErrorContext(ctx, "failed to serve gRPC public server", slog.Any("error", err))
				}
			}()

			// Запускаем public HTTP gateway
			servePublicHTTP, err := in.PublicHTTPServer.Listen()
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to start HTTP public server", slog.Any("error", err))
			}
			in.Logger.InfoContext(ctx, "started HTTP public server", slog.Int("port", in.Cfg.BackendApp.HTTP.Port))

			go func() {
				if err := servePublicHTTP(); err != nil {
					in.Logger.ErrorContext(ctx, "failed to serve HTTP public server", slog.Any("error", err))
				}
			}()

			// Запускаем invoke функции после открытия
			for _, invokeItem := range in.Invokes {
				if invokeItem.StartAfterOpen != nil {
					err := invokeItem.StartAfterOpen(ctxWork)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn start after open", slog.Any("error", err))
						return err
					}
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Останавливаем public HTTP gateway
			in.Logger.InfoContext(ctx, "stopping public HTTP server")
			err := in.PublicHTTPServer.Shutdown(ctx)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to shutdown public HTTP gateway", slog.Any("error", err))
			}

			// Останавливаем public gRPC
			in.Logger.InfoContext(ctx, "stopping public gRPC server")
			in.PublicGRPCServer.Server().Shutdown()

			for _, invokeItem := range in.Invokes {
				if invokeItem.Stop != nil {
					err := invokeItem.Stop(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn stop", slog.Any("error", err))
					}
				}
			}

			cancel()

			// Останавливаем jobs
			err = in.JobsController.StopAll(ctx)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to stop jobs", slog.Any("error", err))
			} else {
				in.Logger.InfoContext(ctx, "jobs stopped")
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
