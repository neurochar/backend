package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonalityTraitType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "bipolar", PersonalityTraitTypeBipolar.String())
	assert.Equal(t, "", PersonalityTraitTypeUnspecified.String())
	assert.Equal(t, "", PersonalityTraitType(99).String())
}

func TestPersonalityTraitBipolar(t *testing.T) {
	t.Parallel()

	trait := &PersonalityTraitBipolar{
		ID:             42,
		Name:           "Extraversion",
		Description:    "Tendency to be outgoing",
		LeftStateName:  "Introvert",
		RightStateName: "Extravert",
	}

	assert.Equal(t, uint64(42), trait.GetID())
	assert.Equal(t, PersonalityTraitTypeBipolar, trait.GetType())
	assert.Equal(t, "Extraversion", trait.GetName())
	assert.Equal(t, "Tendency to be outgoing", trait.GetDescription())
	assert.Equal(t, "Introvert", trait.GetLeftStateName())
	assert.Equal(t, "Extravert", trait.GetRightStateName())
}

func TestNewProfile_Success(t *testing.T) {
	t.Parallel()

	tenantID := uuid.New()
	createdBy := uuid.New()
	traitsMap := ProfilePersonalityTraitsMap{
		1: {Priority: TraitPriorityHigh, Target: 70},
	}

	profile, err := NewProfile(tenantID, &createdBy, "Developer Profile", "Description", traitsMap)
	require.NoError(t, err)
	require.NotNil(t, profile)

	assert.NotEqual(t, uuid.Nil, profile.ID)
	assert.Equal(t, tenantID, profile.TenantID)
	assert.Equal(t, "Developer Profile", profile.Name)
	assert.Equal(t, "Description", profile.Description)
	assert.Equal(t, &createdBy, profile.CreatedBy)
	assert.Equal(t, traitsMap, profile.PersonalityTraitsMap)
}

func TestNewProfile_EmptyName(t *testing.T) {
	t.Parallel()

	_, err := NewProfile(uuid.New(), nil, "  ", "desc", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrProfileInvalidName)
}

func TestNewProfile_NilTraitsMap(t *testing.T) {
	t.Parallel()

	profile, err := NewProfile(uuid.New(), nil, "Test", "desc", nil)
	require.NoError(t, err)
	require.NotNil(t, profile)
	assert.NotNil(t, profile.PersonalityTraitsMap)
	assert.Empty(t, profile.PersonalityTraitsMap)
}

func TestProfile_SetName(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetName("  New Name  ")
	require.NoError(t, err)
	assert.Equal(t, "New Name", profile.Name)
}

func TestProfile_SetName_Empty(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetName("  ")
	assert.ErrorIs(t, err, ErrProfileInvalidName)
}

func TestProfile_SetDescription(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetDescription("  Some Description  ")
	require.NoError(t, err)
	assert.Equal(t, "Some Description", profile.Description)
}

func TestProfile_SetPersonalityTraitsMap(t *testing.T) {
	t.Parallel()

	t.Run("non-nil map", func(t *testing.T) {
		profile := &Profile{}
		m := ProfilePersonalityTraitsMap{1: {Priority: TraitPriorityLow, Target: 50}}
		err := profile.SetPersonalityTraitsMap(m)
		require.NoError(t, err)
		assert.Equal(t, m, profile.PersonalityTraitsMap)
	})

	t.Run("nil map", func(t *testing.T) {
		profile := &Profile{}
		err := profile.SetPersonalityTraitsMap(nil)
		require.NoError(t, err)
		assert.NotNil(t, profile.PersonalityTraitsMap)
		assert.Empty(t, profile.PersonalityTraitsMap)
	})
}

func TestProfile_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	profile := &Profile{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), profile.Version())
}

func TestRoomStatusType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "not_started", RoomStatusTypeNotStarted.String())
	assert.Equal(t, "finished", RoomStatusTypeFinished.String())
	assert.Equal(t, "", RoomStatusTypeUnspecified.String())
	assert.Equal(t, "", RoomStatusTypeStarted.String())
	assert.Equal(t, "", RoomStatusType(99).String())
}

func TestNewRoom_Success(t *testing.T) {
	t.Parallel()

	tenantID := uuid.New()
	createdBy := uuid.New()
	candidateID := uuid.New()
	profileID := uuid.New()

	room, err := NewRoom(tenantID, &createdBy, candidateID, profileID)
	require.NoError(t, err)
	require.NotNil(t, room)

	assert.NotEqual(t, uuid.Nil, room.ID)
	assert.Equal(t, tenantID, room.TenantID)
	assert.Equal(t, RoomStatusTypeNotStarted, room.Status)
	assert.Equal(t, &createdBy, room.CreatedBy)
	require.NotNil(t, room.CandidateID)
	assert.Equal(t, candidateID, *room.CandidateID)
	require.NotNil(t, room.ProfileID)
	assert.Equal(t, profileID, *room.ProfileID)
	assert.Nil(t, room.StartedAt)
	assert.Nil(t, room.FinishedAt)
	assert.False(t, room.IsProcessed)
}

func TestRoom_SetCandidateID(t *testing.T) {
	t.Parallel()

	room := &Room{}
	id := uuid.New()
	err := room.SetCandidateID(&id)
	require.NoError(t, err)
	require.NotNil(t, room.CandidateID)
	assert.Equal(t, id, *room.CandidateID)
}

func TestRoom_SetProfileID(t *testing.T) {
	t.Parallel()

	room := &Room{}
	id := uuid.New()
	err := room.SetProfileID(&id)
	require.NoError(t, err)
	require.NotNil(t, room.ProfileID)
	assert.Equal(t, id, *room.ProfileID)
}

func TestRoom_SetResultIndex(t *testing.T) {
	t.Parallel()

	room := &Room{}
	idx := 5
	err := room.SetResultIndex(&idx)
	require.NoError(t, err)
	require.NotNil(t, room.ResultIndex)
	assert.Equal(t, 5, *room.ResultIndex)
}

func TestRoom_SetPersonalityTraitsMap(t *testing.T) {
	t.Parallel()

	t.Run("non-nil", func(t *testing.T) {
		room := &Room{}
		m := ProfilePersonalityTraitsMap{1: {Priority: TraitPriorityMedium, Target: 60}}
		err := room.SetPersonalityTraitsMap(m)
		require.NoError(t, err)
		assert.Equal(t, m, room.PersonalityTraitsMap)
	})

	t.Run("nil", func(t *testing.T) {
		room := &Room{}
		err := room.SetPersonalityTraitsMap(nil)
		require.NoError(t, err)
		assert.NotNil(t, room.PersonalityTraitsMap)
		assert.Empty(t, room.PersonalityTraitsMap)
	})
}

func TestRoom_SetTechniqueData(t *testing.T) {
	t.Parallel()

	t.Run("non-nil", func(t *testing.T) {
		room := &Room{}
		data := []RoomTechniqueDataItem{{TechniqueID: 1}}
		err := room.SetTechniqueData(data)
		require.NoError(t, err)
		assert.Len(t, room.TechniqueData, 1)
	})

	t.Run("nil", func(t *testing.T) {
		room := &Room{}
		err := room.SetTechniqueData(nil)
		require.NoError(t, err)
		assert.NotNil(t, room.TechniqueData)
		assert.Empty(t, room.TechniqueData)
	})
}

func TestRoom_Duration(t *testing.T) {
	t.Parallel()

	t.Run("both nil", func(t *testing.T) {
		room := &Room{}
		assert.Nil(t, room.Duration())
	})

	t.Run("started not nil, finished nil", func(t *testing.T) {
		now := time.Now()
		room := &Room{StartedAt: &now, FinishedAt: nil}
		assert.Nil(t, room.Duration())
	})

	t.Run("both set", func(t *testing.T) {
		start := time.Now()
		finish := start.Add(30 * time.Minute)
		room := &Room{StartedAt: &start, FinishedAt: &finish}
		dur := room.Duration()
		require.NotNil(t, dur)
		assert.Equal(t, 30*time.Minute, *dur)
	})
}

func TestRoom_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	room := &Room{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), room.Version())
}

func TestTechniqueItemType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "question_with_variants_single_answer", TechniqueItemTypeQuestionWithVariantsSignleAnswer.String())
	assert.Equal(t, "", TechniqueItemTypeUnspecified.String())
	assert.Equal(t, "", TechniqueItemType(99).String())
}

func TestRoomResult(t *testing.T) {
	t.Parallel()

	result := &RoomResult{
		TotalMatchTip: "Good match",
		Techniques:    make(map[uint64]RoomResultTechnique),
		Traits:        make(map[uint64]RoomResultTraitItem),
	}

	assert.Equal(t, "Good match", result.TotalMatchTip)
	assert.Equal(t, "Good match", result.TotalMatchTip)
	assert.Empty(t, result.Techniques)
	assert.Empty(t, result.Traits)
	assert.Nil(t, result.Analyze)
}

func TestRoomResultTechniquesItem(t *testing.T) {
	t.Parallel()

	item := RoomResultTechniquesItem{}
	assert.Zero(t, item.Result)
}

func TestRoomResultTraitItem(t *testing.T) {
	t.Parallel()

	item := RoomResultTraitItem{
		Tip: "Strong trait",
	}
	assert.Zero(t, item.TotalResult)
	assert.Zero(t, item.Match)
	assert.Equal(t, "Strong trait", item.Tip)
}

func TestRoomResultAnalyze(t *testing.T) {
	t.Parallel()

	analyze := &RoomResultAnalyze{
		HiringDecision:     RoomResultAnalyzeHiringDecisionHire,
		ConfidenceScore:    0.85,
		MainRecommendation: "Proceed with hiring",
		PersonalityFit: RoomResultAnalyzePersonalityFit{
			Score:      8,
			Summary:    "Excellent fit",
			KeyMatches: []string{"Leadership", "Teamwork"},
			KeyGaps:    []string{"Experience"},
		},
		Risks:       []string{"Onboarding time"},
		ActionItems: []string{"Schedule interview"},
	}

	assert.Equal(t, RoomResultAnalyzeHiringDecisionHire, analyze.HiringDecision)
	assert.Equal(t, 0.85, analyze.ConfidenceScore)
	assert.Len(t, analyze.PersonalityFit.KeyMatches, 2)
	assert.Len(t, analyze.PersonalityFit.KeyGaps, 1)
	assert.Len(t, analyze.Risks, 1)
	assert.Len(t, analyze.ActionItems, 1)
}
