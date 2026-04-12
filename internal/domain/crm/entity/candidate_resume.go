package entity

import (
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var CandidateResumeStatusUnknown = appErrors.ErrBadRequest.Extend("resume status unknown")

type CandidateResumeStatus uint8

const (
	CandidateResumeStatusUnspecified  CandidateResumeStatus = 0
	CandidateResumeStatusNew          CandidateResumeStatus = 1
	CandidateResumeStatusToProcess    CandidateResumeStatus = 2
	CandidateResumeStatusProcessing   CandidateResumeStatus = 3
	CandidateResumeStatusProcessed    CandidateResumeStatus = 10
	CandidateResumeStatusProcessError CandidateResumeStatus = 99
)

func CandidateResumeStatusFromUint8(value uint8) CandidateResumeStatus {
	switch value {
	case 1:
		return CandidateResumeStatusNew
	case 2:
		return CandidateResumeStatusToProcess
	case 3:
		return CandidateResumeStatusProcessing
	case 10:
		return CandidateResumeStatusProcessed
	case 99:
		return CandidateResumeStatusProcessError
	default:
		return CandidateResumeStatusUnspecified
	}
}

type CandidateResumeFileType uint8

const (
	CandidateResumeFileTypeUnspecified CandidateResumeFileType = 0
	CandidateResumeFileTypePdf         CandidateResumeFileType = 1
	CandidateResumeFileTypeWord        CandidateResumeFileType = 2
)

func CandidateResumeFileTypeFromUint8(value uint8) CandidateResumeFileType {
	switch value {
	case 1:
		return CandidateResumeFileTypePdf
	case 2:
		return CandidateResumeFileTypeWord
	default:
		return CandidateResumeFileTypeUnspecified
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
	FileType    CandidateResumeFileType
	AnalyzeData *CandidateResumeAnalyzeData
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (item *CandidateResume) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *CandidateResume) FilesIDs() []uuid.UUID {
	return []uuid.UUID{
		item.FileID,
	}
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
	fileType CandidateResumeFileType,
) (*CandidateResume, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	account := &CandidateResume{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Status:    CandidateResumeStatusNew,
		FileID:    fileID,
		FileHash:  fileHash,
		FileType:  fileType,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	return account, nil
}
