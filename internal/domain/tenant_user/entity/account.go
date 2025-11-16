package entity

import (
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/emailnormalize"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
)

var ErrAccountInvalidEmail = appErrors.ErrBadRequest.Extend("invalid email").WithTextCode("INVALID_EMAIL")

var ErrAccountInvalidPassword = appErrors.ErrBadRequest.Extend("invalid password").WithTextCode("INVALID_PASSWORD")

var ErrAccountProfileInvalidName = appErrors.ErrBadRequest.Extend("invalid name").WithTextCode("INVALID_NAME")

var ErrAccountProfileInvalidSurname = appErrors.ErrBadRequest.Extend("invalid surname").WithTextCode("INVALID_SURNAME")

type Account struct {
	ID                         uuid.UUID
	TenantID                   uuid.UUID
	RoleID                     uint64
	Email                      string
	PasswordHash               string
	IsConfirmed                bool
	IsEmailVerified            bool
	IsBlocked                  bool
	LastLoginAt                *time.Time
	LastRequestAt              *time.Time
	LastRequestIP              *net.IP
	ProfileName                string
	ProfileSurname             string
	ProfilePhoto100x100FileID  *uuid.UUID
	ProfilePhotoOriginalFileID *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Account) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Account) FilesIDs() []uuid.UUID {
	result := make([]uuid.UUID, 0, 2)

	if item.ProfilePhoto100x100FileID != nil {
		result = append(result, *item.ProfilePhoto100x100FileID)
	}

	if item.ProfilePhotoOriginalFileID != nil {
		result = append(result, *item.ProfilePhotoOriginalFileID)
	}

	return result
}

func (item *Account) SetPassword(password string) error {
	if len(password) < 8 {
		return ErrAccountInvalidPassword
	}

	var hasLetter, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	count := 0
	if hasLetter {
		count++
	}
	if hasDigit {
		count++
	}
	if hasSpecial {
		count++
	}

	if count < 2 {
		return ErrAccountInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	item.PasswordHash = string(hash)

	return nil
}

func (item *Account) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(item.PasswordHash), []byte(password))
	return err == nil
}

func (item *Account) SetEmail(email string) error {
	email = strings.TrimSpace(email)

	err := validate.Var(email, "required,email")
	if err != nil {
		return ErrAccountInvalidEmail
	}

	res, err := emailnormalize.Normalize(email)
	if err != nil {
		return ErrAccountInvalidEmail
	}

	item.Email = res.NormalizedAddress

	return nil
}

func (item *Account) SetRoleID(value uint64) error {
	item.RoleID = value

	return nil
}

func (item *Account) SetLastLoginAt(value *time.Time) {
	if value != nil {
		value = lo.ToPtr(value.Truncate(time.Microsecond))
	}

	item.LastLoginAt = value
}

func (item *Account) SetLastRequestAt(value *time.Time) {
	if value != nil {
		value = lo.ToPtr(value.Truncate(time.Microsecond))
	}

	item.LastRequestAt = value
}

func (item *Account) SetLastRequestIP(value *net.IP) {
	item.LastRequestIP = value
}

func (item *Account) SetProfileName(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrAccountProfileInvalidName
	}

	item.ProfileName = value

	return nil
}

func (item *Account) SetProfileSurname(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrAccountProfileInvalidSurname
	}

	item.ProfileSurname = value

	return nil
}

func NewAccount(
	tenantID uuid.UUID,
	email string,
	password string,
	roleID uint64,
	isConfirmed bool,
	isEmailVerified bool,
) (*Account, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	account := &Account{
		ID:              uuid.New(),
		TenantID:        tenantID,
		RoleID:          roleID,
		IsConfirmed:     isConfirmed,
		IsEmailVerified: isEmailVerified,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
	}

	err := account.SetEmail(email)
	if err != nil {
		return nil, err
	}

	err = account.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return account, nil
}
