package room

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
)

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateRoomDataInput,
) (*usecase.RoomDTO, error) {
	const op = "CreateByDTO"

	var authorAccountID *uuid.UUID
	authData := auth.GetAuthData(ctx)
	if authData != nil {
		authorAccountID = &authData.AccountID
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

		room.TechniqueData = make([]entity.RoomTechniqueDataItem, 0)

		for _, technique := range techniques.TechniquesLib {
			items := technique.GetItemsByPersonalityTraits(room.PersonalityTraitsMap)

			for _, item := range items {
				room.TechniqueData = append(room.TechniqueData, entity.RoomTechniqueDataItem{
					TechniqueID: technique.GetID(),
					ItemData:    item,
				})
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
		candidate, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckRights(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || authData.TenantID != candidate.TenantID {
				return appErrors.ErrForbidden
			}
		}

		if !skipVersionCheck && candidate.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, candidate.Version()).
				WithDetail("last_updated_at", false, candidate.UpdatedAt)
		}

		err = uc.repo.Update(ctx, candidate)
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
