package jobs

import (
	"context"
	"time"
)

func (ctrl *Controller) processFilesToDelete(
	ctx context.Context,
	timeout time.Duration,
	failedTimeout time.Duration,
) (time.Duration, error) {
	anyJobDone, err := ctrl.fileUC.JobProcessFilesToDelete(ctx)
	if err != nil {
		return failedTimeout, err
	}

	if anyJobDone {
		return 0, nil
	}

	return timeout, nil
}
