package app

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	t11, err := hashPassword(t1)
	assert.NoError(t, err)
	t21, err := hashPassword(t2)
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
	admin := roles.CreateAdminUser("test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	require.NoError(t, err)

	store, err := dynamodb.New(sess, awsdb.DynamodbConfig{Prefix: "test"})
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

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	admin := roles.CreateAdminUser("test", roles.Device{})
	ctx = roles.CtxWithAuthInfo(ctx, admin)

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	if err != nil {
		t.Fatal(err)
	}

	store, err := dynamodb.New(sess, awsdb.DynamodbConfig{Prefix: "test"})
	if err != nil {
		t.Fatal(err)
	}

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
