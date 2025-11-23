package registration

import "github.com/google/uuid"

type OutTenant struct {
	ID     uuid.UUID `json:"id"`
	TextID string    `json:"textID"`
	URL    string    `json:"url"`
}
