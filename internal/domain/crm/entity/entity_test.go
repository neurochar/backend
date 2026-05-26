package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCandidateGenderFromUint8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   uint8
		want    CandidateGender
		wantErr bool
	}{
		{name: "unspecified", input: 0, want: CandidateGenderUnspecified, wantErr: false},
		{name: "male", input: 1, want: CandidateGenderMale, wantErr: false},
		{name: "female", input: 2, want: CandidateGenderFemale, wantErr: false},
		{name: "unknown", input: 99, want: CandidateGenderUnspecified, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CandidateGenderFromUint8(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrCandidateGenderUnknown)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewCandidate_Success(t *testing.T) {
	t.Parallel()

	tenantID := uuid.New()
	createdBy := uuid.New()

	candidate, err := NewCandidate(tenantID, &createdBy, "John", "Doe")
	require.NoError(t, err)
	require.NotNil(t, candidate)

	assert.NotEqual(t, uuid.Nil, candidate.ID)
	assert.Equal(t, tenantID, candidate.TenantID)
	assert.Equal(t, "John", candidate.CandidateName)
	assert.Equal(t, "Doe", candidate.CandidateSurname)
	assert.Equal(t, &createdBy, candidate.CreatedBy)
	assert.Equal(t, CandidateGenderUnspecified, candidate.CandidateGender)
	assert.Nil(t, candidate.CandidateBirthday)
}

func TestNewCandidate_InvalidName(t *testing.T) {
	t.Parallel()

	_, err := NewCandidate(uuid.New(), nil, "  ", "Doe")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCandidateInvalidName)
}

func TestNewCandidate_InvalidSurname(t *testing.T) {
	t.Parallel()

	_, err := NewCandidate(uuid.New(), nil, "John", "  ")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCandidateInvalidSurname)
}

func TestCandidate_SetCandidateName(t *testing.T) {
	t.Parallel()

	candidate := &Candidate{}
	err := candidate.SetCandidateName("  Alice  ")
	require.NoError(t, err)
	assert.Equal(t, "Alice", candidate.CandidateName)
}

func TestCandidate_SetCandidateName_Empty(t *testing.T) {
	t.Parallel()

	candidate := &Candidate{}
	err := candidate.SetCandidateName("  ")
	assert.ErrorIs(t, err, ErrCandidateInvalidName)
}

func TestCandidate_SetCandidateSurname(t *testing.T) {
	t.Parallel()

	candidate := &Candidate{}
	err := candidate.SetCandidateSurname("  Smith  ")
	require.NoError(t, err)
	assert.Equal(t, "Smith", candidate.CandidateSurname)
}

func TestCandidate_SetCandidateSurname_Empty(t *testing.T) {
	t.Parallel()

	candidate := &Candidate{}
	err := candidate.SetCandidateSurname("  ")
	assert.ErrorIs(t, err, ErrCandidateInvalidSurname)
}

func TestCandidate_SetCandidateGender(t *testing.T) {
	t.Parallel()

	candidate := &Candidate{}
	err := candidate.SetCandidateGender(CandidateGenderMale)
	require.NoError(t, err)
	assert.Equal(t, CandidateGenderMale, candidate.CandidateGender)
}

func TestCandidate_SetCandidateBirthday(t *testing.T) {
	t.Parallel()

	candidate := &Candidate{}
	birthday := time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)
	err := candidate.SetCandidateBirthday(&birthday)
	require.NoError(t, err)
	require.NotNil(t, candidate.CandidateBirthday)
	assert.Equal(t, birthday, *candidate.CandidateBirthday)
}

func TestCandidate_CalcAge(t *testing.T) {
	t.Parallel()

	t.Run("nil birthday", func(t *testing.T) {
		candidate := &Candidate{}
		assert.Nil(t, candidate.CalcAge(time.Now()))
	})

	t.Run("exact birthday", func(t *testing.T) {
		candidate := &Candidate{CandidateBirthday: loPtr(time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC))}
		now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
		age := candidate.CalcAge(now)
		require.NotNil(t, age)
		assert.Equal(t, 34, *age)
	})

	t.Run("before birthday this year", func(t *testing.T) {
		candidate := &Candidate{CandidateBirthday: loPtr(time.Date(1990, 12, 25, 0, 0, 0, 0, time.UTC))}
		now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
		age := candidate.CalcAge(now)
		require.NotNil(t, age)
		assert.Equal(t, 33, *age)
	})

	t.Run("after birthday this year", func(t *testing.T) {
		candidate := &Candidate{CandidateBirthday: loPtr(time.Date(1990, 3, 1, 0, 0, 0, 0, time.UTC))}
		now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
		age := candidate.CalcAge(now)
		require.NotNil(t, age)
		assert.Equal(t, 34, *age)
	})

	t.Run("leap year birthday", func(t *testing.T) {
		candidate := &Candidate{CandidateBirthday: loPtr(time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC))}
		now := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
		age := candidate.CalcAge(now)
		require.NotNil(t, age)
		assert.Equal(t, 24, *age)
	})
}

func TestCandidate_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	candidate := &Candidate{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), candidate.Version())
}

func TestCandidateResumeStatusFromUint8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input uint8
		want  CandidateResumeStatus
	}{
		{name: "unspecified", input: 0, want: CandidateResumeStatusUnspecified},
		{name: "new", input: 1, want: CandidateResumeStatusNew},
		{name: "to_process", input: 2, want: CandidateResumeStatusToProcess},
		{name: "processing", input: 3, want: CandidateResumeStatusProcessing},
		{name: "processed", input: 10, want: CandidateResumeStatusProcessed},
		{name: "error", input: 99, want: CandidateResumeStatusProcessError},
		{name: "unknown", input: 50, want: CandidateResumeStatusUnspecified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CandidateResumeStatusFromUint8(tt.input))
		})
	}
}

func TestCandidateResumeFileTypeFromUint8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input uint8
		want  CandidateResumeFileType
	}{
		{name: "unspecified", input: 0, want: CandidateResumeFileTypeUnspecified},
		{name: "pdf", input: 1, want: CandidateResumeFileTypePdf},
		{name: "word", input: 2, want: CandidateResumeFileTypeWord},
		{name: "unknown", input: 99, want: CandidateResumeFileTypeUnspecified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CandidateResumeFileTypeFromUint8(tt.input))
		})
	}
}

func TestNewCandidateResume(t *testing.T) {
	t.Parallel()

	tenantID := uuid.New()
	fileID := uuid.New()
	fileHash := "sha256:abc123"

	resume, err := NewCandidateResume(tenantID, fileID, fileHash, CandidateResumeFileTypePdf)
	require.NoError(t, err)
	require.NotNil(t, resume)

	assert.NotEqual(t, uuid.Nil, resume.ID)
	assert.Equal(t, tenantID, resume.TenantID)
	assert.Equal(t, CandidateResumeStatusNew, resume.Status)
	assert.Equal(t, fileID, resume.FileID)
	assert.Equal(t, fileHash, resume.FileHash)
	assert.Equal(t, CandidateResumeFileTypePdf, resume.FileType)
	assert.Nil(t, resume.CandidateID)
	assert.Nil(t, resume.AnalyzeData)
	assert.Nil(t, resume.ErrorText)
}

func TestCandidateResume_SetStatus(t *testing.T) {
	t.Parallel()

	resume := &CandidateResume{}
	err := resume.SetStatus(CandidateResumeStatusProcessed)
	require.NoError(t, err)
	assert.Equal(t, CandidateResumeStatusProcessed, resume.Status)
}

func TestCandidateResume_SetCandidateID(t *testing.T) {
	t.Parallel()

	resume := &CandidateResume{}
	id := uuid.New()
	err := resume.SetCandidateID(&id)
	require.NoError(t, err)
	require.NotNil(t, resume.CandidateID)
	assert.Equal(t, id, *resume.CandidateID)
}

func TestCandidateResume_SetAnalyzeData(t *testing.T) {
	t.Parallel()

	resume := &CandidateResume{}
	data := &CandidateResumeAnalyzeData{AnonymizedText: "some text", DataVersion: 1}
	err := resume.SetAnalyzeData(data)
	require.NoError(t, err)
	require.NotNil(t, resume.AnalyzeData)
	assert.Equal(t, "some text", resume.AnalyzeData.AnonymizedText)
	assert.Equal(t, int64(1), resume.AnalyzeData.DataVersion)
}

func TestCandidateResume_SetErrorText(t *testing.T) {
	t.Parallel()

	resume := &CandidateResume{}
	errText := "processing failed"
	err := resume.SetErrorText(&errText)
	require.NoError(t, err)
	require.NotNil(t, resume.ErrorText)
	assert.Equal(t, "processing failed", *resume.ErrorText)
}

func TestCandidateResume_FilesIDs(t *testing.T) {
	t.Parallel()

	fileID := uuid.New()
	resume := &CandidateResume{FileID: fileID}
	files := resume.FilesIDs()
	assert.Len(t, files, 1)
	assert.Equal(t, fileID, files[0])
}

func TestCandidateResume_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	resume := &CandidateResume{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), resume.Version())
}

func loPtr[T any](v T) *T {
	return &v
}
