package room

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) entitiesToDTO(
	ctx context.Context,
	items []*entity.Room,
	dtoOpts *usecase.RoomDTOOptions,
) ([]*usecase.RoomDTO, error) {
	const op = "entitiesToDTO"

	tenantAccountsMap := make(map[uuid.UUID]*tenantUC.AccountDTO, 0)
	tenantAccountsIDs := make([]uuid.UUID, 0)

	candidatesMap := make(map[uuid.UUID]*candidateUC.CandidateDTO, 0)
	candidatesIDs := make([]uuid.UUID, 0)

	profilesMap := make(map[uuid.UUID]*usecase.ProfileDTO, 0)
	profilesIDs := make([]uuid.UUID, 0)

	for _, item := range items {
		if item.CreatedBy != nil {
			tenantAccountsIDs = append(tenantAccountsIDs, *item.CreatedBy)
		}

		if item.CandidateID != nil {
			candidatesIDs = append(candidatesIDs, *item.CandidateID)
		}

		if item.ProfileID != nil {
			profilesIDs = append(profilesIDs, *item.ProfileID)
		}
	}

	if (dtoOpts == nil || dtoOpts.FetchCreatedBy) && len(tenantAccountsIDs) > 0 {
		accountsList, err := uc.tenantAccountUC.FindList(ctx, &tenantUC.AccountListOptions{
			FilterIDs: &tenantAccountsIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		tenantAccountsMap = lo.SliceToMap(accountsList, func(item *tenantUC.AccountDTO) (uuid.UUID, *tenantUC.AccountDTO) {
			return item.Account.ID, item
		})
	}

	if (dtoOpts == nil || dtoOpts.FetchCandidate) && len(candidatesIDs) > 0 {
		candidatesList, err := uc.candidateUC.FindList(ctx, &candidateUC.CandidateListOptions{
			FilterIDs: &candidatesIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		candidatesMap = lo.SliceToMap(
			candidatesList,
			func(item *candidateUC.CandidateDTO) (uuid.UUID, *candidateUC.CandidateDTO) {
				return item.Candidate.ID, item
			},
		)
	}

	if (dtoOpts == nil || dtoOpts.FetchProfile) && len(profilesIDs) > 0 {
		profilesList, err := uc.profileUC.FindList(ctx, &usecase.ProfileListOptions{
			FilterIDs: &profilesIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		profilesMap = lo.SliceToMap(
			profilesList,
			func(item *usecase.ProfileDTO) (uuid.UUID, *usecase.ProfileDTO) {
				return item.Profile.ID, item
			},
		)
	}

	out := make([]*usecase.RoomDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.RoomDTO{
			Room: item,
		}

		if (dtoOpts == nil || dtoOpts.FetchCreatedBy) && item.CreatedBy != nil {
			account, ok := tenantAccountsMap[*item.CreatedBy]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("account not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.CreatedBy = account
		}

		if (dtoOpts == nil || dtoOpts.FetchCandidate) && item.CandidateID != nil {
			candidate, ok := candidatesMap[*item.CandidateID]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("candidate not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.CandidateDTO = candidate
		}

		if (dtoOpts == nil || dtoOpts.FetchProfile) && item.ProfileID != nil {
			profile, ok := profilesMap[*item.ProfileID]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("profile not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.ProfileDTO = profile
		}

		out = append(out, resItem)
	}

	return out, nil
}

var MetricTextMap = map[decimal.Decimal]string{
	decimal.MustNew(100, 0): "Полное совпадение с целевым значением метрики.",
	decimal.MustNew(90, 0):  "Очень близко к целевому уровню. Отличия минимальны.",
	decimal.MustNew(80, 0):  "Высокая степень соответствия метрики ожидаемому значению.",
	decimal.MustNew(70, 0):  "Хорошее соответствие. Есть небольшое отклонение, но оно некритично.",
	decimal.MustNew(60, 0):  "Умеренное совпадение. Отклонение заметно, но допустимо.",
	decimal.MustNew(50, 0):  "Среднее соответствие. Метрика совпадает с целевой примерно наполовину.",
	decimal.MustNew(40, 0):  "Отклонение выше среднего. Метрика заметно уходит от целевого уровня.",
	decimal.MustNew(30, 0):  "Низкое совпадение. Уровень метрики значительно отличается от целевого.",
	decimal.MustNew(20, 0):  "Очень низкое совпадение. Метрика почти полностью отличается от целевой.",
	decimal.MustNew(10, 0):  "Крайне низкое совпадение. Практически противоположное выражение метрики.",
	decimal.MustNew(0, 0):   "Максимальное отличие от целевого значения.",
}

type MetricMatch struct {
	Match    decimal.Decimal
	Priority int
}

func (uc *UsecaseImpl) сalcOverallMatch(metrics []MetricMatch) (decimal.Decimal, error) {
	if len(metrics) == 0 {
		return decimal.Zero, nil
	}

	sumWeighted := decimal.Zero
	sumWeights := decimal.Zero

	for _, m := range metrics {
		if m.Priority <= 0 {
			continue
		}
		if m.Priority > 3 {
			m.Priority = 3
		}

		w, err := decimal.New(int64(m.Priority), 0)
		if err != nil {
			return decimal.Zero, err
		}

		weighted, err := m.Match.Mul(w)
		if err != nil {
			return decimal.Zero, err
		}

		sumWeighted, err = sumWeighted.Add(weighted)
		if err != nil {
			return decimal.Zero, err
		}

		sumWeights, err = sumWeights.Add(w)
		if err != nil {
			return decimal.Zero, err
		}
	}

	if sumWeights.IsZero() {
		return decimal.Zero, nil
	}

	overall, err := sumWeighted.Quo(sumWeights)
	if err != nil {
		return decimal.Zero, err
	}

	overall = overall.Round(2)

	return overall, nil
}

var ProfileMatchText = map[int]string{
	100: "Профиль практически полностью совпадает с целевым. Максимальная совместимость.",
	90:  "Очень высокое совпадение профиля. Отличия минимальны и несущественны.",
	80:  "Высокое совпадение. Большинство параметров профиля соответствуют целевой модели.",
	70:  "Хорошее совпадение профиля. Есть отдельные отличия, но общая картина близка.",
	60:  "Умеренное совпадение. Профиль в целом подходит, но заметны отклонения.",
	50:  "Среднее совпадение. Часть характеристик совпадает, часть заметно отличается.",
	40:  "Ниже среднего. Значимые различия между профилем и целевой моделью.",
	30:  "Низкое совпадение. Многие важные черты отличаются от целевого значения.",
	20:  "Очень низкое совпадение. Профиль почти не соответствует целевой модели.",
	10:  "Крайне низкое совпадение. Профиль выражен противоположно целевому.",
	0:   "Полное несовпадение профиля с целевой моделью.",
}

func (uc *UsecaseImpl) roundToProfileBucket(d decimal.Decimal) (int, error) {
	is95 := d.Cmp(decimal.MustNew(95, 0))
	if is95 >= 0 {
		return 100, nil
	}

	ten := decimal.MustNew(10, 0)

	div, err := d.Quo(ten)
	if err != nil {
		return 0, err
	}

	floored := div.Floor(0)

	result, err := floored.Mul(ten)
	if err != nil {
		return 0, err
	}

	i, _, ok := result.Int64(0)
	if !ok {
		return 0, appErrors.ErrInternal
	}

	return int(i), nil
}

func (uc *UsecaseImpl) processRoom(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "processRoom"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		roomDTO, err := uc.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		}, nil)
		if err != nil {
			return err
		}

		candidateGender := crmEntity.CandidateGenderUnknown
		var candidateBirthday *time.Time

		if roomDTO.CandidateDTO != nil {
			candidateGender = roomDTO.CandidateDTO.Candidate.CandidateGender
			candidateBirthday = roomDTO.CandidateDTO.Candidate.CandidateBirthday
		}

		room := roomDTO.Room

		room.Result = &entity.RoomResult{
			Techniques: make(map[uint64]entity.RoomResultTechnique),
			Traits:     make(map[uint64]entity.RoomResultTraitItem),
		}

		techniqueToAnswers := make(map[uint64]map[uint64]any, 0)

		for i, techniqueData := range room.TechniqueData {
			if _, ok := techniqueToAnswers[techniqueData.TechniqueID]; !ok {
				techniqueToAnswers[techniqueData.TechniqueID] = make(map[uint64]any, 0)
			}

			val, ok := room.CandidateAnswerData[uint64(i)]
			if !ok {
				return appErrors.ErrInternal
			}

			techniqueToAnswers[techniqueData.TechniqueID][uint64(i)] = val
		}

		for techniqueID, answers := range techniqueToAnswers {
			technique, ok := techniques.TechniquesLib[techniqueID]
			if !ok {
				return appErrors.ErrInternal
			}

			roomResultTechnique, err := technique.CountResult(
				room.PersonalityTraitsMap,
				room.TechniqueData,
				answers,
				candidateGender,
				candidateBirthday,
			)
			if err != nil {
				return err
			}

			room.Result.Techniques[techniqueID] = roomResultTechnique
		}

		metricMatches := make([]MetricMatch, 0)

		for traitID, traitConfig := range room.PersonalityTraitsMap {
			traitValueSum := decimal.Zero
			count := decimal.Zero

			for _, roomResultTechnique := range room.Result.Techniques {
				for resTraitID, traitTechResult := range roomResultTechnique {
					if resTraitID != traitID {
						continue
					}

					traitValueSum, err = traitValueSum.Add(traitTechResult.Result)
					if err != nil {
						return appErrors.ErrInternal.WithWrap(err)
					}

					count, err = count.Add(decimal.One)
					if err != nil {
						return appErrors.ErrInternal.WithWrap(err)
					}
				}
			}

			traitValue, err := traitValueSum.Quo(count)
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			traitTarget, err := decimal.NewFromInt64(int64(traitConfig.Target), 0, 0)
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			match, err := traitTarget.Sub(traitValue)
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			match = match.Abs()

			match, err = match.Mul(decimal.MustNew(10, 0))
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			match, err = decimal.Hundred.Sub(match)
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			tip := MetricTextMap[match]

			room.Result.Traits[traitID] = entity.RoomResultTraitItem{
				TotalResult: traitValue,
				Match:       match,
				Tip:         tip,
			}

			metricMatches = append(metricMatches, MetricMatch{
				Priority: traitConfig.Priority,
				Match:    match,
			})
		}

		totalMatch, err := uc.сalcOverallMatch(metricMatches)
		if err != nil {
			return err
		}

		room.Result.TotalMatch = totalMatch

		bucket, _ := uc.roundToProfileBucket(totalMatch)
		room.Result.TotalMatchTip = ProfileMatchText[bucket]

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
