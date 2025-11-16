package entity

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant/constants"
	"github.com/samber/lo"
)

var (
	ErrTenantInvalidTextID = appErrors.ErrBadRequest.Extend("invalid textID")
	ErrTenantInvalidName   = appErrors.ErrBadRequest.Extend("invalid name")
)

type Tenant struct {
	ID       uuid.UUID
	TextID   string
	IsDemo   bool
	IsActive bool
	Name     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Version - get version
func (item *Tenant) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

var subdomainRegexp = regexp.MustCompile(`^(?i)[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

func (item *Tenant) SetTextID(value string) error {
	value = strings.ToLower(strings.TrimSpace(value))

	if lo.Contains(constants.TenantTextIDBlacklist, value) {
		return ErrTenantInvalidTextID.WithTextCode("TEXTID_FORBIDDEN")
	}

	if len(value) < 3 {
		return ErrTenantInvalidTextID.WithTextCode("TEXTID_TOO_SHORT")
	}

	if len(value) > 63 {
		return ErrTenantInvalidTextID.WithTextCode("TEXTID_TOO_LONG")
	}

	if !subdomainRegexp.MatchString(value) {
		return ErrTenantInvalidTextID.WithTextCode("TEXTID_INVALID")
	}

	item.TextID = value
	return nil
}

func (item *Tenant) SetIsActive(value bool) {
	item.IsActive = value
}

func (item *Tenant) SetName(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrTenantInvalidName
	}

	item.Name = value

	return nil
}

func (item *Tenant) GetDomain(mainDomain string) string {
	return fmt.Sprintf("%s.%s", item.TextID, mainDomain)
}

func (item *Tenant) GetUrl(mainDomain string, scheme string) string {
	return fmt.Sprintf("%s://%s.%s", scheme, item.TextID, mainDomain)
}

func NewTenant(textID string, name string, isDemo bool) (*Tenant, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Tenant{
		ID:     uuid.New(),
		IsDemo: isDemo,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := item.SetTextID(textID)
	if err != nil {
		return nil, err
	}

	err = item.SetName(name)
	if err != nil {
		return nil, err
	}

	return item, nil
}
