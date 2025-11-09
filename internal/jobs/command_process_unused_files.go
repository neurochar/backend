package jobs

import (
	"context"
	"time"
)

func (ctrl *Controller) processUnusedFiles(
	ctx context.Context,
	timeout time.Duration,
	failedTimeout time.Duration,
	unusedTTL time.Duration,
) (time.Duration, error) {
	anyJobDone, err := ctrl.fileUC.JobProcessUnusedFiles(ctx, unusedTTL)
	if err != nil {
		return failedTimeout, err
	}

	if anyJobDone {
		return 0, nil
	}

	return timeout, nil
}
