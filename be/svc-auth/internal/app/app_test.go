package app

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	authpb "github.com/ggrrrr/btmt-ui/be/svc-auth/authpb/v1"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/mem"
)

var mockUser string = "mockuser"

func TestCheckPassword(t *testing.T) {
	var ok bool
	t1 := "newpass"
	t2 := "newpass"
	t11, err := HashPassword(t1)
	assert.NoError(t, err)
	t21, err := HashPassword(t2)
	assert.NoError(t, err)

	ok = checkPasswordHash(t1, t11)
	assert.True(t, ok)

	ok = checkPasswordHash(t2, t21)
	assert.True(t, ok)

	ok = checkPasswordHash(t2, "newpasS")
	assert.True(t, !ok)
}

type testCase struct {
	test string
	prep func(*testing.T)
}

func TestLogin(t *testing.T) {

	ctx := context.Background()
	admin := roles.CreateSystemAdminUser(roles.SystemRealm, "test", app.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	// store, err := dynamodb.New(cfg())
	store, err := mem.New()
	require.NoError(t, err)

	testApp, err := New(WithTokenTTL(10*time.Minute, 2*time.Hour), WithAuthRepo(store), WithHistoryRepo(store), WithTokenSigner(token.NewSignerMock()))
	require.NoError(t, err)

	authItem := ddd.AuthPasswd{
		Subject:     "test@asd",
		Passwd:      "newpass",
		Status:      ddd.StatusEnabled,
		SystemRoles: []string{"admin"},
	}

	authItemLocked := ddd.AuthPasswd{
		Subject: "test@asdlocked",
		Passwd:  "newpass",
		Status:  ddd.StatusDisable,
	}

	// if this fail it will attempt to create table.
	// We will ignore first error coz it will.
	_ = testApp.UserCreate(ctx, authItem)
	err = testApp.UserCreate(ctx, authItem)
	require.NoError(t, err)

	err = testApp.UserCreate(ctx, authItemLocked)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "Login ok",
			prep: func(t *testing.T) {
				loginToken, err := testApp.LoginPasswd(ctx, authItem.Subject, authItem.Passwd)
				require.NoError(t, err)
				testAuthToken(t,
					ddd.LoginToken{
						Subject:      authItem.Subject,
						AccessToken:  ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(10 * time.Minute)},
						RefreshToken: ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(2 * time.Hour)},
					},
					loginToken,
				)
				h, err := store.ListHistory(ctx, "")
				require.NoError(t, err)
				assert.Equalf(t, 1, len(h), "history not updated")
				assert.Equal(t, loginToken.ID, h[0].ID)

				refreshCtx := roles.CtxWithAuthInfo(ctx, roles.AuthInfo{
					ID:      loginToken.ID,
					Subject: loginToken.Subject,
					Realm:   "localhost",
					Roles:   []string{authpb.AuthSvc_TokenRefresh_FullMethodName},
				})

				newLoginToken, err := testApp.TokenRefresh(refreshCtx)
				require.NoError(t, err)
				require.NotNil(t, newLoginToken)
				testAuthToken(t,
					ddd.LoginToken{
						Subject:     authItem.Subject,
						AccessToken: ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(10 * time.Minute)},
					},
					newLoginToken,
				)

			},
		},
		{
			test: "wrong pass",
			prep: func(t *testing.T) {
				_, err = testApp.LoginPasswd(ctx, authItem.Subject, "authItem.Passwd")
				assert.ErrorIs(t, err, ErrAuthUserPassword)
			},
		},
		{
			test: "EmailNotFound",
			prep: func(t *testing.T) {
				_, err = testApp.LoginPasswd(ctx, "authItem.Email", "authItem.Passwd")
				assert.ErrorIs(t, err, ErrAuthEmailNotFound)
			},
		},
		{
			test: "empty email",
			prep: func(t *testing.T) {
				_, err = testApp.LoginPasswd(ctx, "", "authItem.Passwd")
				assert.ErrorIs(t, err, ErrAuthEmailEmpty)
			},
		},
		{
			test: "account locked",
			prep: func(t *testing.T) {
				_, err = testApp.LoginPasswd(ctx, authItemLocked.Subject, "authItem.Passwd")
				assert.ErrorIs(t, err, ErrAuthEmailLocked)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.prep)
	}

}

func TestValidate(t *testing.T) {
	ctx := context.Background()
	admin := roles.CreateSystemAdminUser(roles.SystemRealm, "test", app.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)
	ctxNoEmail := roles.CtxWithAuthInfo(ctx, roles.AuthInfo{})

	// store, err := dynamodb.New(cfg())
	store, err := mem.New()
	require.NoError(t, err)

	testApp, err := New(
		WithAuthRepo(store),
		WithTokenSigner(token.NewSignerMock()),
		WithTokenTTL(1, 2),
	)
	require.NoError(t, err)

	authItem := ddd.AuthPasswd{
		Subject:     "test@asd",
		Passwd:      "newpass",
		Status:      ddd.StatusEnabled,
		SystemRoles: []string{"admin"},
	}

	authItemLocked := ddd.AuthPasswd{
		Subject: "test@asdlocked",
		Passwd:  "newpass",
		Status:  ddd.StatusDisable,
	}

	// if this fail it will attempt to create table.
	// We will ignore first error coz it will.
	_ = testApp.UserCreate(ctx, authItem)
	err = testApp.UserCreate(ctx, authItem)
	require.NoError(t, err)

	err = testApp.UserCreate(ctx, authItemLocked)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "Validate empty auth info",
			prep: func(tt *testing.T) {
				err := testApp.TokenValidate(ctxNoEmail)
				assert.ErrorIs(tt, err, app.ErrAuthUnauthenticated)
			},
		},
		{
			test: "Validate ErrAuthEmailNotFound",
			prep: func(tt *testing.T) {
				authInfoNotFound := roles.AuthInfo{
					Subject: "asdasdasd",
				}
				testCtx := roles.CtxWithAuthInfo(ctx, authInfoNotFound)
				err := testApp.TokenValidate(testCtx)
				assert.ErrorIs(tt, err, ErrAuthEmailNotFound)
			},
		},
		{
			test: "Validate locked",
			prep: func(tt *testing.T) {
				authInfoNotFound := roles.AuthInfo{
					Subject: authItemLocked.Subject,
				}
				testCtx := roles.CtxWithAuthInfo(ctx, authInfoNotFound)
				err := testApp.TokenValidate(testCtx)
				assert.ErrorIs(tt, err, ErrAuthEmailLocked)
			},
		},
		{
			test: "Validate ok",
			prep: func(tt *testing.T) {
				authInfoNotFound := roles.AuthInfo{
					Subject: authItem.Subject,
				}
				testCtx := roles.CtxWithAuthInfo(ctx, authInfoNotFound)
				err := testApp.TokenValidate(testCtx)
				assert.NoError(tt, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.prep)
	}

}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	admin := roles.CreateSystemAdminUser(roles.SystemRealm, "test", app.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	store, err := mem.New()
	require.NoError(t, err)

	testApp, err := New(WithTokenTTL(1*time.Minute, 2*time.Hour), WithAuthRepo(store), WithHistoryRepo(store), WithTokenSigner(token.NewSignerMock()))
	if err != nil {
		t.Fatal(err)
	}

	authItem := ddd.AuthPasswd{
		Subject:     "test@asd",
		Passwd:      "newpass",
		Status:      ddd.StatusEnabled,
		SystemRoles: []string{"admin"},
	}

	_ = testApp.UserCreate(ctx, authItem)
	err = testApp.UserCreate(ctx, authItem)
	require.NoError(t, err)

	jwt, err := testApp.LoginPasswd(ctx, authItem.Subject, authItem.Passwd)
	assert.NoError(t, err)
	testAuthToken(t,
		ddd.LoginToken{
			Subject:      authItem.Subject,
			AccessToken:  ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(1 * time.Minute)},
			RefreshToken: ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(2 * time.Hour)},
		},
		jwt,
	)

	err = testApp.UserChangePasswd(ctx, authItem.Subject, "authItem.Passwd", "newpass")
	assert.ErrorIs(t, err, ErrAuthUserPassword)

	err = testApp.UserChangePasswd(ctx, authItem.Subject, authItem.Passwd, "newpass")
	assert.NoError(t, err)

	jwt, err = testApp.LoginPasswd(ctx, authItem.Subject, "newpass")
	assert.NoError(t, err)
	testAuthToken(t,
		ddd.LoginToken{
			Subject:      authItem.Subject,
			AccessToken:  ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(1 * time.Minute)},
			RefreshToken: ddd.AuthToken{Value: mockUser, ExpiresAt: time.Now().Add(2 * time.Hour)},
		},
		jwt,
	)
}

func testAuthToken(t *testing.T, expected ddd.LoginToken, actual ddd.LoginToken) bool {
	assert.Truef(t, actual.ID != uuid.Nil, "id must not be nil %+v", actual)
	assert.WithinDurationf(t, expected.AccessToken.ExpiresAt, actual.AccessToken.ExpiresAt, 1*time.Second, "token expr time")
	assert.Equalf(t, expected.AccessToken.Value, actual.AccessToken.Value, "token")

	if actual.RefreshToken.Value != "" {
		assert.WithinDurationf(t, expected.RefreshToken.ExpiresAt, actual.RefreshToken.ExpiresAt, 1*time.Second, "RefreshToken expr time")
		assert.Equalf(t, expected.RefreshToken.Value, actual.RefreshToken.Value, "RefreshToken")
	}
	return true
}
