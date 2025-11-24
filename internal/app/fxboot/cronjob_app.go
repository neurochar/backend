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
	"github.com/neurochar/backend/internal/domain/alert"
	"github.com/neurochar/backend/internal/domain/crm"
	emailingModule "github.com/neurochar/backend/internal/domain/emailing"
	"github.com/neurochar/backend/internal/domain/file"
	"github.com/neurochar/backend/internal/domain/tenant"
	"github.com/neurochar/backend/internal/domain/testing"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/storage"
	"github.com/neurochar/backend/internal/infra/storage/s3d"
	"github.com/neurochar/backend/internal/jobs"
	"github.com/neurochar/backend/pkg/pgclient"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	storageMigrations "github.com/neurochar/backend/internal/infra/storage/migrations"
)

// CronjobAppGetOptionsMap returns fx.Options for cronjob
func CronjobAppGetOptionsMap(appID app.ID, cfg config.Config) OptionsMap {
	return OptionsMap{
		Providing: map[ProvidingID]fx.Option{
			ProvidingAppID: fx.Provide(func() app.ID {
				return appID
			}),
			ProvidingIDFXTimeouts: fx.Options(
				fx.StartTimeout(time.Second*time.Duration(cfg.CronjobApp.Base.StartTimeoutSec)),
				fx.StopTimeout(time.Second*time.Duration(cfg.CronjobApp.Base.StopTimeoutSec)),
			),
			ProvidingIDConfig: fx.Provide(func() config.Config {
				return cfg
			}),
			ProvidingIDLogger: fx.Provide(func(cfg config.Config) *slog.Logger {
				return providing.NewLogger(
					cfg.CronjobApp.Name,
					cfg.CronjobApp.Version,
					cfg.CronjobApp.Base.UseLogger,
					cfg.CronjobApp.Base.IsProd,
				)
			}),
			ProvidingIDFXLogger: fx.WithLogger(func(cfg config.Config) fxevent.Logger {
				return providing.NewFXLogger(cfg.CronjobApp.Base.UseFxLogger)
			}),
			ProvidingIDImageProc: fx.Provide(providing.NewImageProc),
			ProvidingIDDBClients: fx.Provide(
				func(logger *slog.Logger, cfg config.Config, shutdown fx.Shutdowner) db.MasterClient {
					return providing.NewDBClients(
						cfg.Postgres.Master.DSN,
						cfg.CronjobApp.Base.LogSQLQueries,
						logger,
						shutdown,
					)
				},
			),
			ProvidingIDEmailing:       fx.Provide(providing.NewEmailing),
			ProvidingIDJobsController: fx.Provide(jobs.NewController),
			ProvidingIDStorageClient:  fx.Provide(providing.NewStorageClient),
			ProvidingIDFileModule:     file.FxModule,
			ProvidingIDEmailingModule: emailingModule.FxModule,
			ProvidingIDAlertModule:    alert.FxModule,
			ProvidingIDTenantModule:   tenant.FxModule,
			ProvidingIDCRMModule:      crm.FxModule,
			ProvidingIDTestingModule:  testing.FxModule,
		},
		Invokes: []fx.Option{
			fx.Invoke(CronjobAppInitInvoke),
		},
	}
}

type CronjobInvokeInput struct {
	fx.In

	LC             fx.Lifecycle
	Shutdowner     fx.Shutdowner
	Invokes        []invoking.InvokeInit `group:"InvokeInit"`
	Logger         *slog.Logger
	Cfg            config.Config
	DBMasterClient db.MasterClient
	S3Client       *s3.Client
	StorageClient  storage.Client
	JobsController *jobs.Controller
}

// CronjobAppInitInvoke - app init
func CronjobAppInitInvoke(
	in CronjobInvokeInput,
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
					err := invokeItem.StartBeforeOpen(ctx)
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

			// Закрываем postgress
			in.DBMasterClient.Close()
			in.Logger.InfoContext(ctx, "closing db clients")

			return nil
		},
	})
}
