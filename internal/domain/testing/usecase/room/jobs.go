package room

import (
	"context"
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) JobProcessRoomsResults(ctx context.Context) (bool, error) {
	const op = "JobProcessRoomsResults"

	anyJonDone := false

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.FindList(ctx, &usecase.RoomListOptions{
			FilterStatus:      lo.ToPtr(entity.RoomStatusTypeFinished),
			FilterIsProcessed: lo.ToPtr(false),
			FilterProcessTries: &uctypes.CompareOption[int]{
				Value: 5,
				Type:  uctypes.CompareLessOrEqual,
			},
			FilterNeedProcessAt: &usecase.RoomListOptionsFilterNeedProcessAt{
				SelectNull: true,
				CompareValue: uctypes.CompareOption[time.Time]{
					Value: time.Now(),
					Type:  uctypes.CompareLessOrEqual,
				},
			},
		}, &uctypes.QueryGetListParams{
			Limit:               1,
			ForUpdateSkipLocked: true,
		}, nil)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		item := items[0]
		item.Room.ProcessTries++

		err = uc.processRoom(ctx, item)
		if err != nil {
			item.Room.ProcessError = lo.ToPtr(err.Error())
			item.Room.NeedProcessAt = lo.ToPtr(time.Now().Add(time.Duration(item.Room.ProcessTries*5) * time.Minute))
		} else {
			item.Room.ProcessError = nil
			item.Room.IsProcessed = true
		}

		err = uc.repo.Update(ctx, item.Room)
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
