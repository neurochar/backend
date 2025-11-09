package jobs

import (
	"context"
	"time"
)

func (ctrl *Controller) processEmailsToDelete(
	ctx context.Context,
	timeout time.Duration,
	failedTimeout time.Duration,
	unusedTTL time.Duration,
) (time.Duration, error) {
	anyJobDone, err := ctrl.emailingUC.JobProcessItemsToDelete(ctx, unusedTTL)
	if err != nil {
		return failedTimeout, err
	}

	if anyJobDone {
		return 0, nil
	}

	return timeout, nil
}
