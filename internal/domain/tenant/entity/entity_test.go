package entity

import (
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantAccount_NewAccount_Success(t *testing.T) {
	t.Parallel()

	tenantID := uuid.New()
	email := "test@example.com"
	password := "StrongPass1!"
	roleID := uint64(1)

	account, err := NewAccount(tenantID, email, password, false, roleID, true, true)
	require.NoError(t, err)
	require.NotNil(t, account)

	assert.NotEqual(t, uuid.Nil, account.ID)
	assert.Equal(t, tenantID, account.TenantID)
	assert.Equal(t, roleID, account.RoleID)
	assert.Equal(t, "test@example.com", account.Email)
	assert.True(t, account.IsConfirmed)
	assert.True(t, account.IsEmailVerified)
	assert.False(t, account.IsBlocked)
	assert.NotEmpty(t, account.PasswordHash)
}

func TestTenantAccount_NewAccount_InvalidEmail(t *testing.T) {
	t.Parallel()

	_, err := NewAccount(uuid.New(), "bad", "StrongPass1!", false, 1, false, false)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAccountInvalidEmail)
}

func TestTenantAccount_NewAccount_WeakPassword(t *testing.T) {
	t.Parallel()

	_, err := NewAccount(uuid.New(), "test@example.com", "short", false, 1, false, false)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAccountInvalidPassword)
}

func TestTenantAccount_NewAccount_SkipPasswordCheck(t *testing.T) {
	t.Parallel()

	account, err := NewAccount(uuid.New(), "test@example.com", "weak", true, 1, false, false)
	require.NoError(t, err)
	require.NotNil(t, account)
	assert.NotEmpty(t, account.PasswordHash)
}

func TestTenantAccount_SetEmail(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetEmail("  USER@Example.COM  ")
	require.NoError(t, err)
	assert.Equal(t, "user@example.com", account.Email)
}

func TestTenantAccount_SetEmail_Invalid(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetEmail("not-email")
	assert.ErrorIs(t, err, ErrAccountInvalidEmail)
}

func TestTenantAccount_SetPassword(t *testing.T) {
	t.Parallel()

	t.Run("valid with check", func(t *testing.T) {
		account := &Account{}
		err := account.SetPassword("StrongPass1!", false)
		require.NoError(t, err)
		assert.NotEmpty(t, account.PasswordHash)
		assert.True(t, account.VerifyPassword("StrongPass1!"))
	})

	t.Run("skip check", func(t *testing.T) {
		account := &Account{}
		err := account.SetPassword("weak", true)
		require.NoError(t, err)
		assert.NotEmpty(t, account.PasswordHash)
		assert.True(t, account.VerifyPassword("weak"))
	})

	t.Run("weak with check", func(t *testing.T) {
		account := &Account{}
		err := account.SetPassword("short", false)
		assert.ErrorIs(t, err, ErrAccountInvalidPassword)
	})
}

func TestTenantAccount_VerifyPassword(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetPassword("StrongPass1!", false)
	require.NoError(t, err)

	assert.True(t, account.VerifyPassword("StrongPass1!"))
	assert.False(t, account.VerifyPassword("wrong"))
	assert.False(t, account.VerifyPassword(""))
}

func TestTenantAccount_FilesIDs(t *testing.T) {
	t.Parallel()

	t.Run("both nil", func(t *testing.T) {
		account := &Account{}
		assert.Empty(t, account.FilesIDs())
	})

	t.Run("only 100x100", func(t *testing.T) {
		fileID := uuid.New()
		account := &Account{ProfilePhoto100x100FileID: &fileID}
		files := account.FilesIDs()
		assert.Len(t, files, 1)
		assert.Equal(t, fileID, files[0])
	})

	t.Run("both set", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		account := &Account{
			ProfilePhoto100x100FileID:  &id1,
			ProfilePhotoOriginalFileID: &id2,
		}
		files := account.FilesIDs()
		assert.Len(t, files, 2)
		assert.Contains(t, files, id1)
		assert.Contains(t, files, id2)
	})
}

func TestTenantAccount_SetLastLoginAt(t *testing.T) {
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

func TestTenantAccount_SetLastRequestAt(t *testing.T) {
	t.Parallel()

	account := &Account{}
	now := time.Now()
	account.SetLastRequestAt(&now)
	require.NotNil(t, account.LastRequestAt)
	assert.Equal(t, now.Truncate(time.Microsecond), *account.LastRequestAt)
}

func TestTenantAccount_SetLastRequestIP(t *testing.T) {
	t.Parallel()

	account := &Account{}
	ip := netip.MustParseAddr("10.0.0.1")
	account.SetLastRequestIP(&ip)
	require.NotNil(t, account.LastRequestIP)
	assert.Equal(t, ip, *account.LastRequestIP)
}

func TestTenantAccount_SetRoleID(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetRoleID(42)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), account.RoleID)
}

func TestTenantAccount_SetProfileName(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetProfileName("  John  ")
	require.NoError(t, err)
	assert.Equal(t, "John", account.ProfileName)
}

func TestTenantAccount_SetProfileName_Empty(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetProfileName("  ")
	assert.ErrorIs(t, err, ErrAccountProfileInvalidName)
}

func TestTenantAccount_SetProfileSurname(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetProfileSurname("  Doe  ")
	require.NoError(t, err)
	assert.Equal(t, "Doe", account.ProfileSurname)
}

func TestTenantAccount_SetProfileSurname_Empty(t *testing.T) {
	t.Parallel()

	account := &Account{}
	err := account.SetProfileSurname("  ")
	assert.ErrorIs(t, err, ErrAccountProfileInvalidSurname)
}

func TestTenantAccount_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	account := &Account{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), account.Version())
}

func TestTenantAccountCode_NewAccountCode(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	ip := netip.MustParseAddr("10.0.0.1")

	code, err := NewAccountCode(accountID, AccountCodeTypePasswordRecovery, &ip)
	require.NoError(t, err)
	require.NotNil(t, code)

	assert.Equal(t, accountID, code.AccountID)
	assert.Equal(t, AccountCodeTypePasswordRecovery, code.Type)
	assert.True(t, code.IsActive)
	assert.Len(t, code.Code, 8)
}

func TestTenantAccountCode_GenerateNumericCode(t *testing.T) {
	t.Parallel()

	code := &AccountCode{}
	err := code.GenerateNumericCode(4)
	require.NoError(t, err)
	assert.Len(t, code.Code, 4)
}

func TestTenantAccountCode_VerifyCode(t *testing.T) {
	t.Parallel()

	code := &AccountCode{Code: "12345678"}
	assert.True(t, code.VerifyCode("12345678"))
	assert.False(t, code.VerifyCode("00000000"))
}

func TestTenantAccountCode_IsAlive(t *testing.T) {
	t.Parallel()

	assert.True(t, (&AccountCode{IsActive: true, CreatedAt: time.Now()}).IsAlive(time.Hour))
	assert.False(t, (&AccountCode{IsActive: false, CreatedAt: time.Now()}).IsAlive(time.Hour))
	assert.False(t, (&AccountCode{IsActive: true, CreatedAt: time.Now().Add(-2 * time.Hour)}).IsAlive(time.Hour))
}

func TestTenantAccountCode_AddAttempt(t *testing.T) {
	t.Parallel()

	code := &AccountCode{}
	code.AddAttempt()
	code.AddAttempt()
	code.AddAttempt()
	assert.Equal(t, 3, code.Attempts)
}

func TestTenantAccountCode_Deactivate(t *testing.T) {
	t.Parallel()

	code := &AccountCode{IsActive: true}
	code.Deactivate()
	assert.False(t, code.IsActive)
}

func TestNewRegistration_Success(t *testing.T) {
	t.Parallel()

	reg, err := NewRegistration("test@example.com", 1)
	require.NoError(t, err)
	require.NotNil(t, reg)

	assert.NotEqual(t, uuid.Nil, reg.ID)
	assert.Equal(t, "test@example.com", reg.Email)
	assert.Equal(t, uint64(1), reg.Tariff)
	assert.False(t, reg.IsFinished)
	assert.Nil(t, reg.TenantID)
}

func TestNewRegistration_InvalidEmail(t *testing.T) {
	t.Parallel()

	_, err := NewRegistration("bad", 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRegistrationInvalidEmail)
}

func TestRegistration_SetEmail(t *testing.T) {
	t.Parallel()

	reg := &Registration{}
	err := reg.SetEmail("  Test@Test.com  ")
	require.NoError(t, err)
	assert.Equal(t, "test@test.com", reg.Email)
}

func TestRegistration_SetRequestIP(t *testing.T) {
	t.Parallel()

	reg := &Registration{}
	ip := netip.MustParseAddr("10.0.0.1")
	err := reg.SetRequestIP(&ip)
	require.NoError(t, err)
	require.NotNil(t, reg.RequestIP)
	assert.Equal(t, ip, *reg.RequestIP)
}

func TestRegistration_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	reg := &Registration{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), reg.Version())
}

func TestTenantRole(t *testing.T) {
	t.Parallel()

	role := &Role{ID: 1, Rank: 10, TextID: "admin"}
	assert.Equal(t, uint64(1), role.ID)
	assert.Equal(t, 10, role.Rank)
	assert.Equal(t, "admin", role.TextID)
}

func TestNewSession(t *testing.T) {
	t.Parallel()

	accountID := uuid.New()
	ip := netip.MustParseAddr("10.0.0.1")
	now := time.Now()

	session := NewSession(accountID, &ip, now, time.Hour)
	require.NotNil(t, session)

	assert.NotEqual(t, uuid.Nil, session.ID)
	assert.Equal(t, accountID, session.AccountID)
	assert.Equal(t, ip, *session.CreateRequestIP)
	assert.Equal(t, ip, *session.RefreshTokenRequestIP)
	assert.Equal(t, uint64(1), session.RefreshVersion)
	assert.NotEqual(t, uuid.Nil, session.RefreshToken)
	assert.True(t, session.RefreshTokenIssuedAt.Equal(now.Truncate(time.Microsecond)))
	assert.True(t, session.RefreshTokenExpiresAt.After(session.RefreshTokenIssuedAt))
}

func TestSession_GenerateNewRefresh(t *testing.T) {
	t.Parallel()

	session := &Session{RefreshVersion: 0}
	now := time.Now()
	ip := netip.MustParseAddr("10.0.0.1")

	session.GenerateNewRefresh(now, 30*time.Minute, &ip)

	assert.Equal(t, uint64(1), session.RefreshVersion)
	assert.NotEqual(t, uuid.Nil, session.RefreshToken)
	assert.Equal(t, ip, *session.RefreshTokenRequestIP)
	assert.True(t, session.RefreshTokenExpiresAt.Sub(session.RefreshTokenIssuedAt) == 30*time.Minute)
}

func TestSession_IsAlive(t *testing.T) {
	t.Parallel()

	t.Run("alive", func(t *testing.T) {
		session := &Session{RefreshTokenExpiresAt: time.Now().Add(time.Hour)}
		assert.False(t, session.IsAlive(time.Now()))
	})

	t.Run("expired", func(t *testing.T) {
		session := &Session{RefreshTokenExpiresAt: time.Now().Add(-time.Hour)}
		assert.True(t, session.IsAlive(time.Now()))
	})
}

func TestSession_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	session := &Session{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), session.Version())
}

func TestSessionRefreshClaims(t *testing.T) {
	t.Parallel()

	sessionID := uuid.New()
	refreshKey := uuid.New()
	claims := &SessionRefreshClaims{
		SessionID:      sessionID,
		RefreshKey:     refreshKey,
		RefreshVersion: 5,
	}

	assert.Equal(t, sessionID, claims.SessionID)
	assert.Equal(t, refreshKey, claims.RefreshKey)
	assert.Equal(t, uint64(5), claims.RefreshVersion)
}

func TestNewTenant_Success(t *testing.T) {
	t.Parallel()

	tenant, err := NewTenant("mycompany", "My Company", false)
	require.NoError(t, err)
	require.NotNil(t, tenant)

	assert.NotEqual(t, uuid.Nil, tenant.ID)
	assert.Equal(t, "mycompany", tenant.TextID)
	assert.Equal(t, "My Company", tenant.Name)
	assert.False(t, tenant.IsDemo)
	assert.False(t, tenant.IsActive)
}

func TestNewTenant_InvalidTextID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		textID string
	}{
		{name: "too short", textID: "ab"},
		{name: "too long", textID: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz123456789012"},
		{name: "starts with dash", textID: "-test"},
		{name: "ends with dash", textID: "test-"},
		{name: "blacklisted", textID: "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTenant(tt.textID, "Name", false)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrTenantInvalidTextID)
		})
	}
}

func TestNewTenant_EmptyName(t *testing.T) {
	t.Parallel()

	_, err := NewTenant("mycompany", "  ", false)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTenantInvalidName)
}

func TestTenant_SetTextID(t *testing.T) {
	t.Parallel()

	tenant := &Tenant{}
	err := tenant.SetTextID("  My-Tenant-123  ")
	require.NoError(t, err)
	assert.Equal(t, "my-tenant-123", tenant.TextID)
}

func TestTenant_SetName(t *testing.T) {
	t.Parallel()

	tenant := &Tenant{}
	err := tenant.SetName("  My Company  ")
	require.NoError(t, err)
	assert.Equal(t, "My Company", tenant.Name)
}

func TestTenant_SetName_Empty(t *testing.T) {
	t.Parallel()

	tenant := &Tenant{}
	err := tenant.SetName("  ")
	assert.ErrorIs(t, err, ErrTenantInvalidName)
}

func TestTenant_SetIsActive(t *testing.T) {
	t.Parallel()

	tenant := &Tenant{IsActive: true}
	tenant.SetIsActive(false)
	assert.False(t, tenant.IsActive)
}

func TestTenant_GetDomain(t *testing.T) {
	t.Parallel()

	tenant := &Tenant{TextID: "mycompany"}
	assert.Equal(t, "mycompany.example.com", tenant.GetDomain("example.com"))
}

func TestTenant_GetUrl(t *testing.T) {
	t.Parallel()

	tenant := &Tenant{TextID: "mycompany"}
	assert.Equal(t, "https://mycompany.example.com", tenant.GetUrl("example.com", "https"))
}

func TestTenant_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	tenant := &Tenant{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), tenant.Version())
}
