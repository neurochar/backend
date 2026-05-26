package room

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques/kettel"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques/kettel/cat"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/convert"
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
	Priority entity.TraitPriority
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
	roomDTO *usecase.RoomDTO,
) error {
	const op = "processRoom"

	roomDTO.Room.Result = &entity.RoomResult{
		Techniques: make(map[uint64]entity.RoomResultTechnique),
		Traits:     make(map[uint64]entity.RoomResultTraitItem),
	}

	roomDTO.Room.Result.Techniques[1] = make(entity.RoomResultTechnique, 0)

	for traitID, traitStatus := range roomDTO.Room.TraitsStatuses {
		if traitStatus.Sten != nil {
			dSten, err := decimal.NewFromInt64(int64(*traitStatus.Sten), 0, 0)
			if err != nil {
				return appErrors.ErrInternal
			}

			roomDTO.Room.Result.Techniques[1][traitID] = entity.RoomResultTechniquesItem{
				Result: dSten,
			}
		}
	}

	metricMatches := make([]MetricMatch, 0)

	for traitID, traitConfig := range roomDTO.Room.PersonalityTraitsMap {
		traitValueSum := decimal.Zero
		count := decimal.Zero

		for _, roomResultTechnique := range roomDTO.Room.Result.Techniques {
			for resTraitID, traitTechResult := range roomResultTechnique {
				if resTraitID != traitID {
					continue
				}

				var err error
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

		roomDTO.Room.Result.Traits[traitID] = entity.RoomResultTraitItem{
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

	roomDTO.Room.Result.TotalMatch = totalMatch

	bucket, _ := uc.roundToProfileBucket(totalMatch)
	roomDTO.Room.Result.TotalMatchTip = ProfileMatchText[bucket]

	totalMatchInt := 0
	i64, _, ok := totalMatch.Int64(0)
	if ok {
		totalMatchInt = int(i64)
	}

	roomDTO.Room.ResultIndex = lo.ToPtr(totalMatchInt)

	llmReq := &usecase.GenerateRoomResultsRequest{
		Job: usecase.GenerateRoomResultsRequestJob{
			Role:        roomDTO.ProfileDTO.Profile.Name,
			Description: roomDTO.ProfileDTO.Profile.Description,
		},
		Candidate: usecase.GenerateRoomResultsRequestCandidate{
			Sex: roomDTO.CandidateDTO.Candidate.CandidateGender,
		},
		PsyTestResult: usecase.GenerateRoomResultsRequestPsyTestResult{
			SumResult: float64(totalMatchInt),
			Traits:    make([]usecase.GenerateRoomResultsRequestTrait, 0, len(roomDTO.Room.TechniqueData)),
		},
	}

	if roomDTO.CandidateDTO.Candidate.CandidateBirthday != nil {
		llmReq.Candidate.Age = roomDTO.CandidateDTO.Candidate.CalcAge(time.Now())
	}

	if roomDTO.CandidateDTO.Resume != nil &&
		roomDTO.CandidateDTO.Resume.Resume != nil &&
		roomDTO.CandidateDTO.Resume.Resume.AnalyzeData != nil {

		llmReq.Candidate.Resume = lo.ToPtr(roomDTO.CandidateDTO.Resume.Resume.AnalyzeData.AnonymizedText)
	}

	for traitID, trait := range roomDTO.Room.Result.Traits {

		teqTrait, err := uc.personalityTraitUC.FindOneByID(ctx, traitID)
		if err != nil {
			continue
		}

		traitTotalRes := 0
		i64, _, ok := trait.TotalResult.Int64(0)
		if ok {
			traitTotalRes = int(i64)
		}

		roomTrait, ok := roomDTO.Room.PersonalityTraitsMap[traitID]
		if !ok {
			continue
		}

		llmReq.PsyTestResult.Traits = append(llmReq.PsyTestResult.Traits, usecase.GenerateRoomResultsRequestTrait{
			Name:           teqTrait.GetName(),
			Description:    teqTrait.GetDescription(),
			LeftStateName:  teqTrait.GetLeftStateName(),
			RightStateName: teqTrait.GetRightStateName(),
			Result:         traitTotalRes,
			Priority:       roomTrait.Priority,
			Target:         roomTrait.Target,
		})
	}

	llmResp, err := uc.repoLLM.GenerateRoomResults(ctx, llmReq)
	if err != nil {
		uc.logger.ErrorContext(ctx, "repoLLM.GenerateRoomResults", slog.Any("error", err))
		roomDTO.Room.Result = nil
		roomDTO.Room.ResultIndex = nil
		return appErrors.ErrInternal.WithWrap(err)
	} else if llmResp != nil {
		roomDTO.Room.Result.Analyze = llmResp.Analyze

		if roomDTO.Room.Result.Analyze != nil {
			score := roomDTO.Room.Result.Analyze.PersonalityFit.Score

			score = max(0, min(100, score))

			total, _, _ := roomDTO.Room.Result.TotalMatch.Int64(0)
			totalInt := int(total)

			const maxDiff = 30

			score = max(totalInt-maxDiff, min(totalInt+maxDiff, score))
			score = max(0, min(100, score))

			roomDTO.Room.Result.Analyze.PersonalityFit.Score = score

			roomDTO.Room.ResultIndex = lo.ToPtr(roomDTO.Room.Result.Analyze.PersonalityFit.Score)
		}
	} else {
		roomDTO.Room.Result = nil
		roomDTO.Room.ResultIndex = nil
		return appErrors.ErrInternal.WithHints("empty llm req")
	}

	return nil
}

var (
	ErrQenerateNextQuestionLastAnswerNotFilled = appErrors.ErrBadRequest.Extend("last answer not filled")
	ErrQenerateNextQuestionAllFinished         = appErrors.ErrBadRequest.Extend("all finished")
)

func (uc *UsecaseImpl) generateNextQuestionForRoom(
	room *entity.Room,
) error {
	if len(room.TechniqueData) > 0 {
		lastIndex := len(room.TechniqueData) - 1
		if _, ok := room.CandidateAnswerData[uint64(lastIndex)]; !ok {
			return ErrQenerateNextQuestionLastAnswerNotFilled
		}
	}

	allTraits := make([]uint64, 0, len(room.PersonalityTraitsMap))
	for traitID := range room.PersonalityTraitsMap {
		allTraits = append(allTraits, traitID)
	}

	notFinishedTraits := make([]uint64, 0, len(room.TraitsStatuses))
	for traitID, trait := range room.TraitsStatuses {
		if trait.Sten == nil {
			notFinishedTraits = append(notFinishedTraits, traitID)
		}
	}

	if len(notFinishedTraits) == 0 {
		return ErrQenerateNextQuestionAllFinished
	}

	useTraitID := lo.Sample(notFinishedTraits)

	if room.TraitsStatuses[useTraitID].UseCat {
		traitAnswers := []cat.SessionAnswer{}
		for i, technique := range room.TechniqueData {
			if technique.TechniqueID != 1 {
				continue
			}

			teqItem, err := technique.ItemData.GetItem()
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			kettelQuestion, ok := kettel.ItemsLib[teqItem.GetID()]
			if !ok {
				return appErrors.ErrInternal
			}

			if kettelQuestion.TraitID != useTraitID {
				continue
			}

			answer, ok := room.CandidateAnswerData[uint64(i)]
			if !ok {
				return appErrors.ErrInternal
			}

			answerVariantID, ok := convert.ToInt(answer)
			if !ok {
				return appErrors.ErrInternal
			}

			traitAnswers = append(traitAnswers, cat.SessionAnswer{
				QuestionID: kettelQuestion.ID,
				VariantID:  answerVariantID,
			})
		}

		result, err := uc.catCtrl.PlaySession(useTraitID, traitAnswers)
		if err != nil {
			return appErrors.ErrInternal.WithWrap(err)
		}

		if result.IsFinished {
			return nil
		} else if result.NextAnswerID != nil {
			room.TechniqueData = append(room.TechniqueData, entity.RoomTechniqueDataItem{
				TechniqueID: 1,
				ItemData:    kettel.NewKettelItem(*result.NextAnswerID),
			})
		}
	} else {
		technique, ok := techniques.TechniquesLib[1]
		if !ok {
			return appErrors.ErrInternal
		}

		traitsItems, err := technique.TraitItems(useTraitID)
		if err != nil {
			return appErrors.ErrInternal.WithWrap(err)
		}

		questions := make([]uint64, 0, len(traitsItems))
		for _, technique := range room.TechniqueData {
			if technique.TechniqueID != 1 {
				continue
			}

			teqItem, err := technique.ItemData.GetItem()
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			kettelQuestion, ok := kettel.ItemsLib[teqItem.GetID()]
			if !ok {
				return appErrors.ErrInternal
			}

			if kettelQuestion.TraitID != useTraitID {
				continue
			}

			questions = append(questions, kettelQuestion.ID)
		}

		notAsked, _ := lo.Difference(traitsItems, questions)

		if len(notAsked) == 0 {
			return nil
		}

		useQuestionID := lo.Sample(notAsked)

		room.TechniqueData = append(room.TechniqueData, entity.RoomTechniqueDataItem{
			TechniqueID: 1,
			ItemData:    kettel.NewKettelItem(useQuestionID),
		})
	}

	return nil
}

var (
	ErrAnswerQuestionForRoomAllAnswered          = appErrors.ErrBadRequest.Extend("all answered")
	ErrAnswerQuestionForRoomInvalidQuestionIndex = appErrors.ErrBadRequest.Extend("invalid question index")
)

func (uc *UsecaseImpl) answerQuestionForRoom(
	ctx context.Context,
	room *entity.Room,
	questionIndex int32,
	answer any,
) error {
	if len(room.TechniqueData) == 0 {
		return appErrors.ErrInternal.Extend("no technique data")
	}

	lastIndex := len(room.TechniqueData) - 1
	if _, ok := room.CandidateAnswerData[uint64(lastIndex)]; ok {
		return ErrAnswerQuestionForRoomAllAnswered
	}

	if int(questionIndex) != lastIndex {
		return ErrAnswerQuestionForRoomInvalidQuestionIndex
	}

	teqItem, err := room.TechniqueData[lastIndex].ItemData.GetItem()
	if err != nil {
		return appErrors.ErrInternal.WithWrap(err)
	}

	kettelQuestion, ok := kettel.ItemsLib[teqItem.GetID()]
	if !ok {
		return appErrors.ErrInternal
	}

	useTraitID := kettelQuestion.TraitID

	var answerInt int
	switch teqItem.GetType() {
	case entity.TechniqueItemTypeQuestionWithVariantsSignleAnswer:
		answerInt, ok = convert.ToInt(answer)
		if !ok {
			return appErrors.ErrBadRequest
		}
	}

	err = teqItem.ValidateAnswer(answerInt)
	if err != nil {
		return err
	}

	if room.CandidateAnswerData == nil {
		room.CandidateAnswerData = make(map[uint64]any, 0)
	}
	room.CandidateAnswerData[uint64(lastIndex)] = answerInt

	if _, ok := room.TraitsStatuses[useTraitID]; !ok {
		return appErrors.ErrInternal
	}

	room.TraitsStatuses[useTraitID].AnsweredCount++

	continueToClassic := false
	if room.TraitsStatuses[useTraitID].UseCat {
		traitAnswers := []cat.SessionAnswer{}
		for i, technique := range room.TechniqueData {
			if technique.TechniqueID != 1 {
				continue
			}

			teqItem, err := technique.ItemData.GetItem()
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			kettelQuestion, ok := kettel.ItemsLib[teqItem.GetID()]
			if !ok {
				return appErrors.ErrInternal
			}

			if kettelQuestion.TraitID != useTraitID {
				continue
			}

			answer, ok := room.CandidateAnswerData[uint64(i)]
			if !ok {
				return appErrors.ErrInternal
			}

			answerVariantID, ok := convert.ToInt(answer)
			if !ok {
				return appErrors.ErrInternal
			}

			traitAnswers = append(traitAnswers, cat.SessionAnswer{
				QuestionID: kettelQuestion.ID,
				VariantID:  answerVariantID,
			})

		}

		result, err := uc.catCtrl.PlaySession(useTraitID, traitAnswers)
		if err != nil {
			return appErrors.ErrInternal.WithWrap(err)
		}

		if result.IsFinished {
			if result.IsSure {
				room.TraitsStatuses[useTraitID].Sten = result.ResultSten
			} else {
				room.TraitsStatuses[useTraitID].UseCat = false
				continueToClassic = true
			}
		}
	}

	if !room.TraitsStatuses[useTraitID].UseCat || continueToClassic {
		technique, ok := techniques.TechniquesLib[1]
		if !ok {
			return appErrors.ErrInternal
		}

		answers := map[uint64]any{}
		for i, technique := range room.TechniqueData {
			if technique.TechniqueID != 1 {
				continue
			}

			teqItem, err := technique.ItemData.GetItem()
			if err != nil {
				return appErrors.ErrInternal.WithWrap(err)
			}

			kettelQuestion, ok := kettel.ItemsLib[teqItem.GetID()]
			if !ok {
				return appErrors.ErrInternal
			}

			if kettelQuestion.TraitID != useTraitID {
				continue
			}

			answer, ok := room.CandidateAnswerData[uint64(i)]
			if !ok {
				return appErrors.ErrInternal
			}

			answers[kettelQuestion.ID] = answer
		}

		traitsItems, err := technique.TraitItems(useTraitID)
		if err != nil {
			return err
		}

		if len(answers) == len(traitsItems) {
			dto, err := uc.entitiesToDTO(ctx, []*entity.Room{room}, &usecase.RoomDTOOptions{
				FetchCandidate: true,
			})
			if err != nil {
				return err
			}
			if len(dto) == 0 {
				return appErrors.ErrInternal
			}

			roomDTO := dto[0]

			candidateGender := crmEntity.CandidateGenderUnspecified
			var candidateBirthday *time.Time

			if roomDTO.CandidateDTO != nil {
				candidateGender = roomDTO.CandidateDTO.Candidate.CandidateGender
				candidateBirthday = roomDTO.CandidateDTO.Candidate.CandidateBirthday
			}

			roomResultTechnique, err := technique.CountTraitSten(
				useTraitID,
				answers,
				candidateGender,
				candidateBirthday,
			)
			if err != nil {
				return err
			}

			room.TraitsStatuses[useTraitID].Sten = lo.ToPtr(roomResultTechnique)
		}
	}

	return nil
}
