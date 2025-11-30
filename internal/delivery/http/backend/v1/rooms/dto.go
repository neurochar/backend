package rooms

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
)

type OutProfile struct {
	Version int64 `json:"_version,omitempty"`

	ID                   uuid.UUID                                 `json:"id"`
	TenantID             uuid.UUID                                 `json:"tenantID"`
	Name                 string                                    `json:"name"`
	PersonalityTraitsMap testingEntity.ProfilePersonalityTraitsMap `json:"personalityTraitsMap"`
}

func OutProfileDTO(
	c *fiber.Ctx,
	profileDTO *testingUC.ProfileDTO,
) (*OutProfile, error) {
	out := &OutProfile{
		Version:  profileDTO.Profile.Version(),
		ID:       profileDTO.Profile.ID,
		TenantID: profileDTO.Profile.TenantID,

		Name:                 profileDTO.Profile.Name,
		PersonalityTraitsMap: profileDTO.Profile.PersonalityTraitsMap,
	}

	return out, nil
}

type OutRoomCandidate struct {
	ID               uuid.UUID `json:"id"`
	CandidateName    string    `json:"candidateName"`
	CandidateSurname string    `json:"candidateSurname"`
}

type OutRoomProfile struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type OutListRoom struct {
	Version int64 `json:"_version,omitempty"`

	ID        uuid.UUID                    `json:"id"`
	Status    testingEntity.RoomStatusType `json:"status"`
	TenantID  uuid.UUID                    `json:"tenantID"`
	Candidate *OutRoomCandidate            `json:"candidate"`
	Profile   *OutRoomProfile              `json:"profile"`
}

func OutListRoomDTO(
	c *fiber.Ctx,
	roomDTO *testingUC.RoomDTO,
) (*OutListRoom, error) {
	out := &OutListRoom{
		Version:  roomDTO.Room.Version(),
		ID:       roomDTO.Room.ID,
		TenantID: roomDTO.Room.TenantID,
	}

	if roomDTO.CandidateDTO != nil {
		out.Candidate = &OutRoomCandidate{
			ID:               roomDTO.CandidateDTO.Candidate.ID,
			CandidateName:    roomDTO.CandidateDTO.Candidate.CandidateName,
			CandidateSurname: roomDTO.CandidateDTO.Candidate.CandidateSurname,
		}
	}

	if roomDTO.ProfileDTO != nil {
		out.Profile = &OutRoomProfile{
			ID:   roomDTO.ProfileDTO.Profile.ID,
			Name: roomDTO.ProfileDTO.Profile.Name,
		}
	}

	return out, nil
}

type OutRoom struct {
	Version int64 `json:"_version,omitempty"`

	ID                   uuid.UUID                                 `json:"id"`
	Status               testingEntity.RoomStatusType              `json:"status"`
	TenantID             uuid.UUID                                 `json:"tenantID"`
	Candidate            *OutRoomCandidate                         `json:"candidate"`
	Profile              *OutRoomProfile                           `json:"profile"`
	PersonalityTraitsMap testingEntity.ProfilePersonalityTraitsMap `json:"personalityTraitsMap"`
}

func OutRoomDTO(
	c *fiber.Ctx,
	roomDTO *testingUC.RoomDTO,
) (*OutRoom, error) {
	out := &OutRoom{
		Version:              roomDTO.Room.Version(),
		ID:                   roomDTO.Room.ID,
		TenantID:             roomDTO.Room.TenantID,
		PersonalityTraitsMap: roomDTO.Room.PersonalityTraitsMap,
	}

	if roomDTO.CandidateDTO != nil {
		out.Candidate = &OutRoomCandidate{
			ID:               roomDTO.CandidateDTO.Candidate.ID,
			CandidateName:    roomDTO.CandidateDTO.Candidate.CandidateName,
			CandidateSurname: roomDTO.CandidateDTO.Candidate.CandidateSurname,
		}
	}

	if roomDTO.ProfileDTO != nil {
		out.Profile = &OutRoomProfile{
			ID:   roomDTO.ProfileDTO.Profile.ID,
			Name: roomDTO.ProfileDTO.Profile.Name,
		}
	}

	return out, nil
}
