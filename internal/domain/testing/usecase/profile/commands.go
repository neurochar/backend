package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
)

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateProfileDataInput,
) (*usecase.ProfileDTO, error) {
	const op = "CreateByDTO"

	var authorAccountID *uuid.UUID
	authData := auth.GetAuthData(ctx)
	if authData != nil && authData.IsTenantUser() {
		authorAccountID = &authData.TenantUserClaims().AccountID
	}

	err := uc.personalityTraitUC.ValidatePersonalityTraitsMap(in.PersonalityTraitsMap)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	candidate, err := entity.NewProfile(
		tenantID,
		authorAccountID,
		in.Name,
		in.Description,
		in.PersonalityTraitsMap,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repo.Create(ctx, candidate)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	candidateDTO, err := uc.entitiesToDTO(ctx, []*entity.Profile{candidate}, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(candidateDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty candidate dto")
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return candidateDTO[0], nil
}

func (uc *UsecaseImpl) PatchByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchProfileDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchByDTO"

	if in.PersonalityTraitsMap != nil {
		err := uc.personalityTraitUC.ValidatePersonalityTraitsMap(*in.PersonalityTraitsMap)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}
	}

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		candidate, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != candidate.TenantID {
				return appErrors.ErrForbidden
			}
		}

		if !skipVersionCheck && candidate.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, candidate.Version()).
				WithDetail("last_updated_at", false, candidate.UpdatedAt)
		}

		if in.Name != nil {
			err := candidate.SetName(*in.Name)
			if err != nil {
				return err
			}
		}

		if in.Description != nil {
			err := candidate.SetDescription(*in.Description)
			if err != nil {
				return err
			}
		}

		if in.PersonalityTraitsMap != nil {
			err := candidate.SetPersonalityTraitsMap(*in.PersonalityTraitsMap)
			if err != nil {
				return err
			}
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

func (uc *UsecaseImpl) Update(ctx context.Context, item *entity.Profile) error {
	const op = "Update"

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) GenerateProfileDescriptionByName(ctx context.Context, name string) (string, error) {
	const op = "GenerateProfileDescriptionByName"

	if name == "" {
		return "", appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", uc.pkg, op)
	}

	res, err := uc.llmRepo.GenerateProfileDescriptionByName(ctx, name)
	if err != nil {
		if errors.Is(err, usecase.ErrLLMInvalidResponse) {
			return "", appErrors.Chainf(usecase.ErrLLMInvalidResponse, "%s.%s", uc.pkg, op)
		}

		if errors.Is(err, usecase.ErrLLMBadRequest) {
			return "", appErrors.Chainf(usecase.ErrLLMBadRequest.WithTextCode("INVALID_PROFILE_NAME"), "%s.%s", uc.pkg, op)
		}

		return "", appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return res, nil
}

func (uc *UsecaseImpl) GenerateProfileTraitsMapByDescription(
	ctx context.Context,
	req *usecase.GenerateProfileTraitsMapByDescriptionRequest,
) (*usecase.GenerateProfileTraitsMapByDescriptionResponse, error) {
	const op = "GenerateProfileTraitsMapByDescription"

	if req.Description == "" || req.Role == "" {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", uc.pkg, op)
	}

	res, err := uc.llmRepo.GenerateProfileTraitsMapByDescription(ctx, req)
	if err != nil {
		if errors.Is(err, usecase.ErrLLMInvalidResponse) {
			return nil, appErrors.Chainf(usecase.ErrLLMInvalidResponse, "%s.%s", uc.pkg, op)
		}

		if errors.Is(err, usecase.ErrLLMBadRequest) {
			return nil, appErrors.Chainf(usecase.ErrLLMBadRequest.WithTextCode("INVALID_PROFILE"), "%s.%s", uc.pkg, op)
		}

		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return res, nil
}
