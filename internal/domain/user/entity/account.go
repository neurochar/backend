// Package entity contains account entity
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

// ErrAccountInvalidEmail - error for invalid email
var ErrAccountInvalidEmail = appErrors.ErrBadRequest.Extend("invalid email").WithTextCode("INVALID_EMAIL")

// ErrAccountInvalidPassword - error for invalid password
var ErrAccountInvalidPassword = appErrors.ErrBadRequest.Extend("invalid password").WithTextCode("INVALID_PASSWORD")

// Account - account entity
type Account struct {
	ID              uuid.UUID
	RoleID          uint64
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	IsBlocked       bool
	LastLoginAt     *time.Time
	LastRequestAt   *time.Time
	LastRequestIP   *net.IP

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Version - get version
func (item *Account) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

// IsConfirmed - check if account is confirmed
func (item *Account) IsConfirmed() bool {
	return item.IsEmailVerified
}

// SetPassword - set password
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

// VerifyPassword - verify password
func (item *Account) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(item.PasswordHash), []byte(password))
	return err == nil
}

// SetEmail - set email
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

// NewAccount - constructor for new account
func NewAccount(email string, password string, roleID uint64, isEmailVerified bool) (*Account, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	account := &Account{
		ID:              uuid.New(),
		RoleID:          roleID,
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
