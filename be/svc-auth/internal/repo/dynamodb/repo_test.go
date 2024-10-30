package dynamodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

var (
	email1 string = "ggrrrr1@gmail.com"
)

func cfg() (awsclient.AwsConfig, awsclient.DynamodbConfig) {
	return awsclient.AwsConfig{
			Region:   "us-east-1",
			Endpoint: "http://localhost:4566",
		},
		awsclient.DynamodbConfig{
			Database: "testdb",
			Prefix:   "test1",
		}
}

func Test_List(t *testing.T) {
	ctx := context.Background()

	r, err := New(cfg())
	require.NoError(t, err)

	_ = r.createTableAuth()

	_, err = r.ListPasswd(ctx, nil)
	require.NoError(t, err)

	list, err := r.ListPasswd(ctx, nil)
	assert.NoError(t, err)
	logger.Info().Any("asda", list)
}

func TestSave(t *testing.T) {
	t.Skipf("aws dynombe needs to be moved to github.com/aws/aws-sdk-go-v2")
	auth1 := ddd.AuthPasswd{
		Subject:     email1,
		Status:      ddd.StatusDisable,
		RealmRoles:  map[string][]string{"local": {"role1"}},
		SystemRoles: []string{"admin"},
	}
	passwd1 := "asd1asd"

	ctx := context.Background()

	r, err := New(cfg())
	r.prefix = "test"
	require.NoError(t, err)

	_ = r.createTableAuth()

	err = r.SavePasswd(ctx, auth1)
	assert.NoError(t, err, err)
	items, err := r.GetPasswd(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	auth1.CreatedAt = items[0].CreatedAt
	assert.Equal(t, auth1, items[0], "expected email")

	err = r.UpdatePassword(ctx, email1, passwd1)
	assert.NoError(t, err, err)
	items, err = r.GetPasswd(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, passwd1, items[0].Passwd, "expected email")
	logger.Debug().Any("items", items).Msg("data")

	err = r.UpdateStatus(ctx, email1, ddd.StatusEnabled)
	assert.NoError(t, err, err)
	items, err = r.GetPasswd(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, ddd.StatusEnabled, items[0].Status, "expected email")

	authUpdate := ddd.AuthPasswd{
		Subject:     auth1.Subject,
		Status:      ddd.StatusPending,
		SystemRoles: []string{"noadmin"},
		RealmRoles:  map[string][]string{"localhost": {"roles"}},
	}

	err = r.Update(ctx, authUpdate)
	assert.NoError(t, err)
	items, err = r.GetPasswd(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	authUpdate.CreatedAt = items[0].CreatedAt
	authUpdate.Passwd = items[0].Passwd
	assert.Equal(t, authUpdate, items[0], "expected email")

}
