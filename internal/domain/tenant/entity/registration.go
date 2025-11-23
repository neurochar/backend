package entity

import (
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/emailnormalize"
)

var ErrRegistrationInvalidEmail = appErrors.ErrBadRequest.Extend("invalid email").WithTextCode("INVALID_EMAIL")

type Registration struct {
	ID         uuid.UUID
	Email      string
	Tariff     uint64
	IsFinished bool
	TenantID   *uuid.UUID
	RequestIP  *net.IP

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Version - get version
func (item *Registration) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Registration) SetEmail(email string) error {
	email = strings.TrimSpace(email)

	err := validate.Var(email, "required,email")
	if err != nil {
		return ErrRegistrationInvalidEmail
	}

	res, err := emailnormalize.Normalize(email)
	if err != nil {
		return ErrRegistrationInvalidEmail
	}

	item.Email = res.NormalizedAddress

	return nil
}

func (item *Registration) SetRequestIP(ip *net.IP) error {
	item.RequestIP = ip

	return nil
}

func NewRegistration(email string, tariff uint64) (*Registration, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Registration{
		ID:     uuid.New(),
		Tariff: tariff,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := item.SetEmail(email)
	if err != nil {
		return nil, err
	}

	return item, nil
}
