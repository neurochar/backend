package jobs

import (
	"context"
	"time"
)

const (
	CommandProcessFilesToDelete = "process_files_to_delete"

	CommandProcessUnusedFiles = "process_unused_files"

	CommandProcessEmailsToSend = "process_emails_to_send"

	CommandProcessEmailsToDelete = "process_emails_to_delete"
)

func (ctrl *Controller) RegisterProcessFilesToDelete(timeout time.Duration, failedTimeout time.Duration) {
	ctrl.registerFn(CommandProcessFilesToDelete, func(ctx context.Context) (time.Duration, error) {
		return ctrl.processFilesToDelete(ctx, timeout, failedTimeout)
	})
}

func (ctrl *Controller) RegisterProcessUnusedFiles(timeout time.Duration, failedTimeout time.Duration, unusedTTL time.Duration) {
	ctrl.registerFn(CommandProcessUnusedFiles, func(ctx context.Context) (time.Duration, error) {
		return ctrl.processUnusedFiles(ctx, timeout, failedTimeout, unusedTTL)
	})
}

func (ctrl *Controller) RegisterProcessEmailsToSend(timeout time.Duration, failedTimeout time.Duration) {
	ctrl.registerFn(CommandProcessEmailsToSend, func(ctx context.Context) (time.Duration, error) {
		return ctrl.processEmailsToSend(ctx, timeout, failedTimeout)
	})
}

func (ctrl *Controller) RegisterProcessEmailsToDelete(timeout time.Duration, failedTimeout time.Duration, unusedTTL time.Duration) {
	ctrl.registerFn(CommandProcessEmailsToDelete, func(ctx context.Context) (time.Duration, error) {
		return ctrl.processEmailsToDelete(ctx, timeout, failedTimeout, unusedTTL)
	})
}
