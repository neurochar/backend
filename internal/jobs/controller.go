// Package jobs - cron and ticker jobs
package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/loghandler"

	emailingUC "github.com/neurochar/backend/internal/domain/emailing/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
)

const timeoutOnPanic = 5 * time.Second

type fnEntity struct {
	name   string
	fn     func(context.Context) (time.Duration, error)
	ctx    context.Context
	cancel context.CancelFunc
	stop   chan struct{}
	going  int32
}

// Controller - cron and ticker jobs controller
type Controller struct {
	pkg        string
	cfg        config.Config
	activeFns  int32
	fns        map[string]*fnEntity
	logger     *slog.Logger
	fileUC     fileUC.Usecase
	emailingUC emailingUC.Usecase
}

// NewController - constructor for Controller
func NewController(
	cfg config.Config,
	logger *slog.Logger,
	fileUC fileUC.Usecase,
	emailingUC emailingUC.Usecase,
) *Controller {
	return &Controller{
		pkg:        "jobs.Controller",
		cfg:        cfg,
		fns:        map[string]*fnEntity{},
		logger:     logger,
		fileUC:     fileUC,
		emailingUC: emailingUC,
	}
}

func (ctrl *Controller) registerFn(name string, fn func(context.Context) (time.Duration, error)) {
	if item, ok := ctrl.fns[name]; ok {
		item.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())

	ctrl.fns[name] = &fnEntity{
		name:   name,
		fn:     fn,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start - start job
func (ctrl *Controller) Start(name string) error {
	item, ok := ctrl.fns[name]
	if !ok {
		return fmt.Errorf("job %s not found", name)
	}

	if !atomic.CompareAndSwapInt32(&item.going, 0, 1) {
		return nil
	}

	item.stop = make(chan struct{})

	go func(item *fnEntity) {
		defer func() {
			close(item.stop)
			atomic.StoreInt32(&item.going, 0)
		}()

		defer func() {
			atomic.AddInt32(&ctrl.activeFns, -1)
		}()

		atomic.AddInt32(&ctrl.activeFns, 1)

		for {
			ctx := loghandler.SetContextData(context.Background(), "job.name", item.name)
			ctx = loghandler.SetContextData(ctx, "job.id", uuid.New())

			if ctrl.cfg.CronjobApp.Base.LogJob {
				ctrl.logger.InfoContext(ctx, "start job")
			}

			start := time.Now()

			timeout, err := func() (d time.Duration, fnErr error) {
				defer func() {
					if r := recover(); r != nil {
						ctrl.logger.ErrorContext(
							ctx,
							"job task panic",
							slog.Any("job.task.panic", r),
							slog.String("stack", string(debug.Stack())),
						)

						d = timeoutOnPanic
						fnErr = nil
					}
				}()

				d, fnErr = item.fn(ctx)
				return
			}()

			timeSince := time.Since(start)

			if err != nil {
				ctrl.logger.ErrorContext(
					ctx,
					"failed job",
					slog.Any("error", err),
					slog.Duration("duration", timeSince),
				)
			} else if ctrl.cfg.CronjobApp.Base.LogJob {
				ctrl.logger.InfoContext(
					ctx,
					"end job",
					slog.Float64("duration_sec", timeSince.Seconds()),
				)
			}

			select {
			case <-item.ctx.Done():
				return
			case <-time.After(timeout):
			}

		}
	}(item)

	return nil
}

// StartAll - start all jobs
func (ctrl *Controller) StartAll() {
	for _, item := range ctrl.fns {
		// nolint
		_ = ctrl.Start(item.name)
	}
}

// Stop - stop job
func (ctrl *Controller) Stop(ctx context.Context, name string) error {
	item, ok := ctrl.fns[name]
	if !ok {
		return fmt.Errorf("job %s not found", name)
	}

	if atomic.LoadInt32(&item.going) == 0 {
		return nil
	}

	item.cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-item.stop:
			return nil
		}
	}
}

// StopAll - stop all jobs
func (ctrl *Controller) StopAll(ctx context.Context) error {
	for _, item := range ctrl.fns {
		item.cancel()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if atomic.LoadInt32(&ctrl.activeFns) == 0 {
				return nil
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}
