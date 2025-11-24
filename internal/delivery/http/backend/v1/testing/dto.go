package testing

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
