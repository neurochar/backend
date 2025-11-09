// Package entity contains account entity
package entity

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type AccountCodeType uint8

const (
	AccountCodeTypeEmailVerification AccountCodeType = 1
	AccountCodeTypePasswordRecovery  AccountCodeType = 2
)

type AccountCode struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Type      AccountCodeType
	IsActive  bool
	Code      string
	RequestIP net.IP
	Attempts  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (item *AccountCode) GenerateNumericCode(n int) error {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
	nBig, err := rand.Int(rand.Reader, max)
	if err != nil {
		return appErrors.ErrInternal.WithWrap(err)
	}

	item.Code = fmt.Sprintf(fmt.Sprintf("%%0%dd", n), nBig)

	return nil
}

func (item *AccountCode) VerifyCode(code string) bool {
	return item.Code == code
}

func (item *AccountCode) IsAlive(maxDuration time.Duration) bool {
	return item.IsActive && time.Since(item.CreatedAt) < maxDuration
}

func (item *AccountCode) AddAttempt() {
	item.Attempts++
}

func (item *AccountCode) Deactivate() {
	item.IsActive = false
}

// NewAccountCode - constructor for new account code.
// By default generates numeric code in 8 digits
func NewAccountCode(accountID uuid.UUID, codeType AccountCodeType, requestIP net.IP) (*AccountCode, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &AccountCode{
		ID:        uuid.New(),
		AccountID: accountID,
		Type:      codeType,
		IsActive:  true,
		RequestIP: requestIP,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := item.GenerateNumericCode(8)
	if err != nil {
		return nil, err
	}

	return item, nil
}
