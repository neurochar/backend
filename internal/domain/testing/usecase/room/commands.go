package room

import (
	"context"
	"errors"
	"net/netip"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateRoomDataInput,
) (*usecase.RoomDTO, error) {
	const op = "CreateByDTO"

	var authorAccountID *uuid.UUID
	authData := auth.GetAuthData(ctx)
	if authData != nil && authData.IsTenantUser() {
		authorAccountID = &authData.TenantUserClaims().AccountID
	}

	var room *entity.Room

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		candidateDTO, err := uc.candidateUC.FindOneByID(ctx, in.CandidateID, &uctypes.QueryGetOneParams{
			ForShare: true,
		}, &candidateUC.CandidateDTOOptions{})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrBadRequest.WithHints("candidateID invalid").WithTextCode("CANDIDATE_ID_INVALID")
			}
			return err
		}

		if candidateDTO.Candidate.TenantID != tenantID {
			return appErrors.ErrForbidden
		}

		profileDTO, err := uc.profileUC.FindOneByID(ctx, in.ProfileID, &uctypes.QueryGetOneParams{
			ForShare: true,
		}, &usecase.ProfileDTOOptions{})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrBadRequest.WithHints("profileID invalid").WithTextCode("PROFILE_ID_INVALID")
			}
			return err
		}

		if profileDTO.Profile.TenantID != tenantID {
			return appErrors.ErrForbidden
		}

		room, err = entity.NewRoom(
			tenantID,
			authorAccountID,
			candidateDTO.Candidate.ID,
			profileDTO.Profile.ID,
		)
		if err != nil {
			return err
		}

		err = room.SetPersonalityTraitsMap(profileDTO.Profile.PersonalityTraitsMap)
		if err != nil {
			return err
		}

		room.TraitsStatuses = make(entity.RoomTraitsStatuses, len(room.PersonalityTraitsMap))
		for traitID := range room.PersonalityTraitsMap {
			room.TraitsStatuses[traitID] = &entity.RoomTraitsStatusesItem{
				AnsweredCount: 0,
				UseCat:        true,
			}
		}

		err = uc.repo.Create(ctx, room)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	roomDTO, err := uc.entitiesToDTO(ctx, []*entity.Room{room}, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(roomDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty room dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return roomDTO[0], nil
}

func (uc *UsecaseImpl) PatchByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchRoomDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		room, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != room.TenantID {
				return appErrors.ErrForbidden
			}
		}

		if !skipVersionCheck && room.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, room.Version()).
				WithDetail("last_updated_at", false, room.UpdatedAt)
		}

		err = uc.repo.Update(ctx, room)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Update(ctx context.Context, item *entity.Room) error {
	const op = "Update"

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Start(
	ctx context.Context,
	id uuid.UUID,
	requestIP *netip.Addr,
) (*usecase.RoomDTO, error) {
	const op = "Start"

	var room *entity.Room

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error
		room, err = uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != room.TenantID {
				return appErrors.ErrForbidden
			}
		}

		if room.Status != entity.RoomStatusTypeNotStarted {
			return appErrors.ErrBadRequest.WithHints("room already started")
		}

		err = uc.generateNextQuestionForRoom(room)
		if err != nil {
			return err
		}

		room.Status = entity.RoomStatusTypeStarted
		room.StartedAt = lo.ToPtr(time.Now())

		err = uc.repo.Update(ctx, room)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	roomDTO, err := uc.entitiesToDTO(ctx, []*entity.Room{room}, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(roomDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty room dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return roomDTO[0], nil
}

func (uc *UsecaseImpl) Answer(
	ctx context.Context,
	id uuid.UUID,
	questionIndex int32,
	answer any,
	requestIP *netip.Addr,
) (*usecase.RoomDTO, error) {
	const op = "Answer"

	var room *entity.Room

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error
		room, err = uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != room.TenantID {
				return appErrors.ErrForbidden
			}
		}

		if room.Status != entity.RoomStatusTypeStarted {
			return appErrors.ErrBadRequest.WithHints("room not in started status")
		}

		err = uc.answerQuestionForRoom(ctx, room, questionIndex, answer)
		if err != nil {
			return err
		}

		err = uc.generateNextQuestionForRoom(room)
		if err != nil {
			if errors.Is(err, ErrQenerateNextQuestionAllFinished) {
				room.Status = entity.RoomStatusTypeFinished
				room.FinishedIP = requestIP
				room.FinishedAt = lo.ToPtr(time.Now())
			} else {
				return err
			}
		}

		err = uc.repo.Update(ctx, room)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	roomDTO, err := uc.entitiesToDTO(ctx, []*entity.Room{room}, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(roomDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty room dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return roomDTO[0], nil
}

func (uc *UsecaseImpl) Process(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "Process"

	room, err := uc.repo.FindOneByID(ctx, id, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckTenantAccess(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != room.TenantID {
			return appErrors.ErrForbidden
		}
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		roomDTO, err := uc.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		}, nil)
		if err != nil {
			return err
		}

		err = uc.processRoom(ctx, roomDTO)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = uc.repo.Update(ctx, roomDTO.Room)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
