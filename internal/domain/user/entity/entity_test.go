package entity

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount_Success(t *testing.T) {
	t.Parallel()

	email := "test@example.com"
	password := "StrongPass1!"
	roleID := uint64(1)

	account, err := NewAccount(email, password, roleID, true)
	require.NoError(t, err)
	require.NotNil(t, account)

	assert.NotEqual(t, uuid.Nil, account.ID)
	assert.Equal(t, roleID, account.RoleID)
	assert.Equal(t, "test@example.com", account.Email)
	assert.True(t, account.IsEmailVerified)
	assert.False(t, account.IsBlocked)
	assert.NotEmpty(t, account.PasswordHash)
	assert.NotEqual(t, password, account.PasswordHash)
	assert.True(t, account.IsConfirmed())
	assert.False(t, account.CreatedAt.IsZero())
	assert.False(t, account.UpdatedAt.IsZero())
	assert.Nil(t, account.DeletedAt)
}

func TestNewAccount_InvalidEmail(t *testing.T) {
	t.Parallel()

	_, err := NewAccount("invalid-email", "StrongPass1!", 1, false)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAccountInvalidEmail)
}

func TestNewAccount_WeakPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
	}{
		{name: "too short", password: "Ab1"},
		{name: "only letters", password: "abcdefgh"},
		{name: "only digits", password: "12345678"},
		{name: "only special", password: "!@#$%^&*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAccount("test@example.com", tt.password, 1, false)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrAccountInvalidPassword)
		})
	}
}

func TestAccount_SetEmail(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetEmail("  Test@Example.COM  ")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", account.Email)
}

func TestAccount_SetEmail_Invalid(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetEmail("not-an-email")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAccountInvalidEmail)
}

func TestAccount_SetPasswordAndVerify(t *testing.T) {
	t.Parallel()

	account := &Account{}
	password := "StrongPass1!"
	err := account.SetPassword(password)
	require.NoError(t, err)

	assert.NotEmpty(t, account.PasswordHash)
	assert.NotEqual(t, password, account.PasswordHash)
	assert.True(t, account.VerifyPassword(password))
	assert.False(t, account.VerifyPassword("WrongPass1!"))
}

func TestAccount_SetPassword_Weak(t *testing.T) {
	t.Parallel()

	t.Run("too short", func(t *testing.T) {
		account := &Account{}
		err := account.SetPassword("Ab1")
		assert.ErrorIs(t, err, ErrAccountInvalidPassword)
	})

	t.Run("only letters", func(t *testing.T) {
		account := &Account{}
		err := account.SetPassword("abcdefgh")
		assert.ErrorIs(t, err, ErrAccountInvalidPassword)
	})
}

func TestAccount_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	account := &Account{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), account.Version())
}

func TestAccount_IsConfirmed(t *testing.T) {
	t.Parallel()

	assert.True(t, (&Account{IsEmailVerified: true}).IsConfirmed())
	assert.False(t, (&Account{IsEmailVerified: false}).IsConfirmed())
}

func TestAccount_SetLastLoginAt(t *testing.T) {
	t.Parallel()

	account := &Account{}
	assert.Nil(t, account.LastLoginAt)

	now := time.Now()
	account.SetLastLoginAt(&now)
	require.NotNil(t, account.LastLoginAt)
	assert.Equal(t, now.Truncate(time.Microsecond), *account.LastLoginAt)

	account.SetLastLoginAt(nil)
	assert.Nil(t, account.LastLoginAt)
}

func TestAccount_SetLastRequestAt(t *testing.T) {
	t.Parallel()

	account := &Account{}
	assert.Nil(t, account.LastRequestAt)

	now := time.Now()
	account.SetLastRequestAt(&now)
	require.NotNil(t, account.LastRequestAt)
	assert.Equal(t, now.Truncate(time.Microsecond), *account.LastRequestAt)

	account.SetLastRequestAt(nil)
	assert.Nil(t, account.LastRequestAt)
}

func TestAccount_SetLastRequestIP(t *testing.T) {
	t.Parallel()

	account := &Account{}
	ip := net.ParseIP("192.168.1.1")
	account.SetLastRequestIP(&ip)
	require.NotNil(t, account.LastRequestIP)
	assert.True(t, ip.Equal(*account.LastRequestIP))
}

func TestNewAccountCode_Success(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	requestIP := net.ParseIP("10.0.0.1")

	code, err := NewAccountCode(accountID, AccountCodeTypeEmailVerification, requestIP)
	require.NoError(t, err)
	require.NotNil(t, code)

	assert.NotEqual(t, uuid.Nil, code.ID)
	assert.Equal(t, accountID, code.AccountID)
	assert.Equal(t, AccountCodeTypeEmailVerification, code.Type)
	assert.True(t, code.IsActive)
	assert.Len(t, code.Code, 8)
	assert.Equal(t, requestIP, code.RequestIP)
	assert.Equal(t, 0, code.Attempts)
}

func TestAccountCode_GenerateNumericCode(t *testing.T) {
	t.Parallel()

	code := &AccountCode{}
	err := code.GenerateNumericCode(6)
	require.NoError(t, err)
	assert.Len(t, code.Code, 6)
}

func TestAccountCode_VerifyCode(t *testing.T) {
	t.Parallel()

	code := &AccountCode{Code: "12345678"}
	assert.True(t, code.VerifyCode("12345678"))
	assert.False(t, code.VerifyCode("87654321"))
	assert.False(t, code.VerifyCode(""))
}

func TestAccountCode_IsAlive(t *testing.T) {
	t.Parallel()

	t.Run("active and fresh", func(t *testing.T) {
		code := &AccountCode{IsActive: true, CreatedAt: time.Now()}
		assert.True(t, code.IsAlive(time.Hour))
	})

	t.Run("inactive", func(t *testing.T) {
		code := &AccountCode{IsActive: false, CreatedAt: time.Now()}
		assert.False(t, code.IsAlive(time.Hour))
	})

	t.Run("expired", func(t *testing.T) {
		code := &AccountCode{IsActive: true, CreatedAt: time.Now().Add(-2 * time.Hour)}
		assert.False(t, code.IsAlive(time.Hour))
	})
}

func TestAccountCode_AddAttempt(t *testing.T) {
	t.Parallel()

	code := &AccountCode{Attempts: 0}
	code.AddAttempt()
	assert.Equal(t, 1, code.Attempts)
	code.AddAttempt()
	assert.Equal(t, 2, code.Attempts)
}

func TestAccountCode_Deactivate(t *testing.T) {
	t.Parallel()

	code := &AccountCode{IsActive: true}
	code.Deactivate()
	assert.False(t, code.IsActive)
}

func TestNewAdminSession(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	ip := net.ParseIP("10.0.0.1")

	session := NewSession(accountID, ip)
	require.NotNil(t, session)

	assert.NotEqual(t, uuid.Nil, session.ID)
	assert.Equal(t, accountID, session.AccountID)
	assert.Equal(t, ip, session.LastRequestIP)
	assert.False(t, session.CreatedAt.IsZero())
	assert.False(t, session.UpdatedAt.IsZero())
	assert.Nil(t, session.DeletedAt)
}

func TestAdminSession_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	session := &AdminSession{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), session.Version())
}

func TestAdminSession_IsAlive(t *testing.T) {
	t.Parallel()

	t.Run("alive", func(t *testing.T) {
		session := &AdminSession{LastRequestAt: time.Now()}
		assert.True(t, session.IsAlive(time.Hour))
	})

	t.Run("expired", func(t *testing.T) {
		session := &AdminSession{LastRequestAt: time.Now().Add(-2 * time.Hour)}
		assert.False(t, session.IsAlive(time.Hour))
	})
}

func TestAdminSession_SetLastRequestAt(t *testing.T) {
	t.Parallel()

	session := &AdminSession{}
	now := time.Now()
	session.SetLastRequestAt(now)
	assert.Equal(t, now.Truncate(time.Microsecond), session.LastRequestAt)
}

func TestAdminSession_SetLastRequestIP(t *testing.T) {
	t.Parallel()

	session := &AdminSession{}
	ip := net.ParseIP("10.0.0.1")
	session.SetLastRequestIP(ip)
	assert.Equal(t, ip, session.LastRequestIP)
}

func TestNewProfile_Success(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	name := "John"
	surname := "Doe"

	profile, err := NewProfile(accountID, name, surname)
	require.NoError(t, err)
	require.NotNil(t, profile)

	assert.Equal(t, accountID, profile.AccountID)
	assert.Equal(t, name, profile.Name)
	assert.Equal(t, surname, profile.Surname)
	assert.Nil(t, profile.Photo100x100FileID)
}

func TestNewProfile_EmptyName(t *testing.T) {
	t.Parallel()

	_, err := NewProfile(uuid.New(), "", "Doe")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrProfileInvalidName)
}

func TestNewProfile_EmptySurname(t *testing.T) {
	t.Parallel()

	_, err := NewProfile(uuid.New(), "John", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrProfileInvalidSurname)
}

func TestProfile_SetName(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetName("  Alice  ")
	require.NoError(t, err)
	assert.Equal(t, "Alice", profile.Name)
}

func TestProfile_SetName_Empty(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetName("  ")
	assert.ErrorIs(t, err, ErrProfileInvalidName)
}

func TestProfile_SetSurname(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetSurname("  Smith  ")
	require.NoError(t, err)
	assert.Equal(t, "Smith", profile.Surname)
}

func TestProfile_SetSurname_Empty(t *testing.T) {
	t.Parallel()

	profile := &Profile{}
	err := profile.SetSurname("  ")
	assert.ErrorIs(t, err, ErrProfileInvalidSurname)
}

func TestProfile_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	profile := &Profile{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), profile.Version())
}

func TestRightType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "bool", RightTypeBool.String())
	assert.Equal(t, "int", RightTypeInt.String())
	assert.Equal(t, "", RightType(0).String())
	assert.Equal(t, "", RightType(99).String())
}

func TestNewRole_Success(t *testing.T) {
	t.Parallel()

	role, err := NewRole("admin")
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, "admin", role.Name)
	assert.False(t, role.IsSystem)
	assert.False(t, role.IsSuper)
}

func TestNewRole_EmptyName(t *testing.T) {
	t.Parallel()

	_, err := NewRole("")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNameEmpty)
}

func TestRole_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	role := &Role{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), role.Version())
}

func TestNewRoleToRight(t *testing.T) {
	t.Parallel()

	item := NewRoleToRight(1, 2, 3)
	require.NotNil(t, item)

	assert.Equal(t, uint64(1), item.RoleID)
	assert.Equal(t, uint64(2), item.RightID)
	assert.Equal(t, 3, item.Value)
	assert.False(t, item.CreatedAt.IsZero())
}
