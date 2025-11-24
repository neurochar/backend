package candidate

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
)

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateCandidateDataInput,
) (*usecase.CandidateDTO, error) {
	const op = "CreateByDTO"

	var authorAccountID *uuid.UUID
	authData := auth.GetAuthData(ctx)
	if authData != nil {
		authorAccountID = &authData.AccountID
	}

	candidate, err := entity.NewCandidate(
		tenantID,
		authorAccountID,
		in.CandidateName,
		in.CandidateSurname,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = candidate.SetCandidateGender(in.CandidateGender)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = candidate.SetCandidateBirthday(in.CandidateBirthday)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repo.Create(ctx, candidate)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	candidateDTO, err := uc.entitiesToDTO(ctx, []*entity.Candidate{candidate}, nil)
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
	in usecase.PatchCandidateDataInput,
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

		if in.CandidateName != nil {
			err := candidate.SetCandidateName(*in.CandidateName)
			if err != nil {
				return err
			}
		}

		if in.CandidateSurname != nil {
			err := candidate.SetCandidateSurname(*in.CandidateSurname)
			if err != nil {
				return err
			}
		}

		if in.CandidateGender != nil {
			err := candidate.SetCandidateGender(*in.CandidateGender)
			if err != nil {
				return err
			}
		}

		if in.CandidateBirthday != nil {
			candidate.SetCandidateBirthday(*in.CandidateBirthday)
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

func (uc *UsecaseImpl) Update(ctx context.Context, item *entity.Candidate) error {
	const op = "Update"

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
