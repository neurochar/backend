package usecase

import (
	"context"
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) JobProcessItemsToSend(ctx context.Context) (bool, error) {
	const op = "JobProcessItemsToSend"

	anyJonDone := false

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(ctx, &ListOptions{
			FilterSentAtCompare: &uctypes.CompareOption[*time.Time]{
				Value: nil,
				Type:  uctypes.CompareEqual,
			},
		}, &uctypes.QueryGetListParams{
			Limit:               1,
			ForUpdateSkipLocked: true,
		})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		item := items[0]

		err = uc.emailing.Send(item.MessageData)
		if err != nil {
			return err
		}

		item.SentAt = lo.ToPtr(time.Now().Truncate(time.Microsecond))

		err = uc.repo.Update(ctx, item)
		if err != nil {
			return err
		}

		anyJonDone = true

		return nil
	})
	if err != nil {
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return anyJonDone, nil
}

func (uc *UsecaseImpl) JobProcessItemsToDelete(ctx context.Context, ttl time.Duration) (bool, error) {
	const op = "JobProcessItemsToDelete"

	anyJonDone := false

	maxTime := time.Now().Add(-ttl)

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(ctx, &ListOptions{
			FilterSentAtCompare: &uctypes.CompareOption[*time.Time]{
				Value: lo.ToPtr(maxTime),
				Type:  uctypes.CompareLessOrEqual,
			},
		}, &uctypes.QueryGetListParams{
			ForUpdateSkipLocked: true,
			Limit:               1,
		})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		err = uc.repo.DeleteByID(ctx, items[0].ID)
		if err != nil {
			return err
		}

		anyJonDone = true

		return nil
	})
	if err != nil {
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return anyJonDone, nil
}
