package entity

import (
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var CandidateResumeStatusUnknown = appErrors.ErrBadRequest.Extend("resume status unknown")

type CandidateResumeStatus uint8

const (
	CandidateResumeStatusNew        CandidateResumeStatus = 0
	CandidateResumeStatusToProcess  CandidateResumeStatus = 1
	CandidateResumeStatusProcessing CandidateResumeStatus = 2
	CandidateResumeStatusProcessed  CandidateResumeStatus = 10
)

func CandidateResumeStatusFromUint8(value uint8) (CandidateResumeStatus, error) {
	switch value {
	case 0:
		return CandidateResumeStatusNew, nil
	case 1:
		return CandidateResumeStatusToProcess, nil
	case 2:
		return CandidateResumeStatusProcessing, nil
	case 10:
		return CandidateResumeStatusProcessed, nil
	default:
		return 0, CandidateResumeStatusUnknown
	}
}

type CandidateResumeAnalyzeData struct {
	AnonymizedText string
}

type CandidateResume struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Status      CandidateResumeStatus
	CandidateID *uuid.UUID
	FileID      uuid.UUID
	FileHash    string
	AnalyzeData *CandidateResumeAnalyzeData
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (item *CandidateResume) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *CandidateResume) SetStatus(value CandidateResumeStatus) error {
	item.Status = value

	return nil
}

func (item *CandidateResume) SetCandidateID(value *uuid.UUID) error {
	item.CandidateID = value

	return nil
}

func (item *CandidateResume) SetAnalyzeData(value *CandidateResumeAnalyzeData) error {
	item.AnalyzeData = value

	return nil
}

func NewCandidateResume(
	tenantID uuid.UUID,
	fileID uuid.UUID,
	fileHash string,
) (*CandidateResume, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	account := &CandidateResume{
		ID:        uuid.New(),
		TenantID:  tenantID,
		FileID:    fileID,
		FileHash:  fileHash,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	return account, nil
}
