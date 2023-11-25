package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/awsdb"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/dynamodb"
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

func cfg() awsdb.AwsConfig {
	return awsdb.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
		Database: awsdb.DynamodbConfig{
			Database: "",
			Prefix:   "test",
		},
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	admin := roles.CreateAdminUser("test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	store, err := dynamodb.New(cfg())
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
	_ = testApp.CreateAuth(ctx, authItem)
	err = testApp.CreateAuth(ctx, authItem)
	require.NoError(t, err)

	err = testApp.CreateAuth(ctx, authItemLocked)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "Login ok",
			prep: func(t *testing.T) {
				jwt, err := testApp.LoginPasswd(ctx, authItem.Email, authItem.Passwd)
				assert.NoError(t, err)
				assert.Equal(t, jwt.Payload(), AuthToken("ok"))
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
	admin := roles.CreateAdminUser("test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)
	ctxNoEmail := roles.CtxWithAuthInfo(ctx, roles.AuthInfo{})

	store, err := dynamodb.New(cfg())
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
	_ = testApp.CreateAuth(ctx, authItem)
	err = testApp.CreateAuth(ctx, authItem)
	require.NoError(t, err)

	err = testApp.CreateAuth(ctx, authItemLocked)
	require.NoError(t, err)

	tests := []testCase{
		{
			test: "Validate empty auth info",
			prep: func(tt *testing.T) {
				err := testApp.Validate(ctxNoEmail)
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
				err := testApp.Validate(testCtx)
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
				err := testApp.Validate(testCtx)
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
				err := testApp.Validate(testCtx)
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
	admin := roles.CreateAdminUser("test", roles.Device{})
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

	_ = testApp.CreateAuth(ctx, authItem)
	err = testApp.CreateAuth(ctx, authItem)
	require.NoError(t, err)

	jwt, err := testApp.LoginPasswd(ctx, authItem.Email, authItem.Passwd)
	assert.NoError(t, err)
	assert.Equal(t, jwt.Payload(), AuthToken("ok"))

	err = testApp.UpdatePasswd(ctx, authItem.Email, "authItem.Passwd", "newpass")
	assert.ErrorIs(t, err, ErrAuthBadPassword)

	err = testApp.UpdatePasswd(ctx, authItem.Email, authItem.Passwd, "newpass")
	assert.NoError(t, err)

	jwt, err = testApp.LoginPasswd(ctx, authItem.Email, "newpass")
	assert.NoError(t, err)
	assert.Equal(t, jwt.Payload(), AuthToken("ok"))
}
