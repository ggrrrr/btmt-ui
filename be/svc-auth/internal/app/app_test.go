package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/dynamodb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/mem"
)

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

func cfg() (awsclient.AwsConfig, awsclient.DynamodbConfig) {
	return awsclient.AwsConfig{
			Region:   "us-east-1",
			Endpoint: "http://localhost:4566",
		},
		awsclient.DynamodbConfig{
			Database: "",
			Prefix:   "test",
		}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	admin := roles.CreateSystemAdminUser(roles.SystemRealm, "test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	// store, err := dynamodb.New(cfg())
	store, err := mem.New()
	require.NoError(t, err)

	testApp, err := New(WithAuthRepo(store), WithTokenSigner(token.NewSignerMock()))
	require.NoError(t, err)

	authItem := ddd.AuthPasswd{
		Email:       "test@asd",
		Passwd:      "newpass",
		Status:      ddd.StatusEnabled,
		SystemRoles: []string{"admin"},
	}

	authItemLocked := ddd.AuthPasswd{
		Email:  "test@asdlocked",
		Passwd: "newpass",
		Status: ddd.StatusDisable,
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
				jwt, err := testApp.LoginPasswd(ctx, authItem.Email, authItem.Passwd)
				assert.NoError(t, err)
				testAuthToken(t, jwt, ddd.AuthToken{Token: "ok", ExpiresAt: time.Now().Add(1 * time.Hour)})

			},
		},
		{
			test: "wrong pass",
			prep: func(t *testing.T) {
				_, err = testApp.LoginPasswd(ctx, authItem.Email, "authItem.Passwd")
				assert.ErrorIs(t, err, ErrAuthBadPassword)
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
				_, err = testApp.LoginPasswd(ctx, authItemLocked.Email, "authItem.Passwd")
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
	admin := roles.CreateSystemAdminUser(roles.SystemRealm, "test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)
	ctxNoEmail := roles.CtxWithAuthInfo(ctx, roles.AuthInfo{})

	// store, err := dynamodb.New(cfg())
	store, err := mem.New()
	require.NoError(t, err)

	testApp, err := New(WithAuthRepo(store), WithTokenSigner(token.NewSignerMock()))
	require.NoError(t, err)

	authItem := ddd.AuthPasswd{
		Email:       "test@asd",
		Passwd:      "newpass",
		Status:      ddd.StatusEnabled,
		SystemRoles: []string{"admin"},
	}

	authItemLocked := ddd.AuthPasswd{
		Email:  "test@asdlocked",
		Passwd: "newpass",
		Status: ddd.StatusDisable,
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
					User: "asdasdasd",
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
					User: authItemLocked.Email,
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
					User: authItem.Email,
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
	admin := roles.CreateSystemAdminUser(roles.SystemRealm, "test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	store, err := dynamodb.New(cfg())
	require.NoError(t, err)

	testApp, err := New(WithAuthRepo(store), WithTokenSigner(token.NewSignerMock()))
	if err != nil {
		t.Fatal(err)
	}

	authItem := ddd.AuthPasswd{
		Email:       "test@asd",
		Passwd:      "newpass",
		Status:      ddd.StatusEnabled,
		SystemRoles: []string{"admin"},
	}

	_ = testApp.UserCreate(ctx, authItem)
	err = testApp.UserCreate(ctx, authItem)
	require.NoError(t, err)

	jwt, err := testApp.LoginPasswd(ctx, authItem.Email, authItem.Passwd)
	assert.NoError(t, err)
	testAuthToken(t, jwt, ddd.AuthToken{Token: "ok", ExpiresAt: time.Now().Add(1 * time.Hour)})

	err = testApp.UserChangePasswd(ctx, authItem.Email, "authItem.Passwd", "newpass")
	assert.ErrorIs(t, err, ErrAuthBadPassword)

	err = testApp.UserChangePasswd(ctx, authItem.Email, authItem.Passwd, "newpass")
	assert.NoError(t, err)

	jwt, err = testApp.LoginPasswd(ctx, authItem.Email, "newpass")
	assert.NoError(t, err)
	testAuthToken(t, jwt, ddd.AuthToken{Token: "ok", ExpiresAt: time.Now().Add(1 * time.Hour)})
}

func testAuthToken(t *testing.T, expected ddd.AuthToken, actual ddd.AuthToken) bool {
	assert.WithinDuration(t, expected.ExpiresAt, actual.ExpiresAt, 1*time.Second)
	assert.Equal(t, expected.Token, actual.Token)
	return true
}
