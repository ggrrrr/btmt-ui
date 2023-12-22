package dynamodb

import (
	"context"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/awsdb"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	email1 string = "ggrrrr1@gmail.com"
)

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

func Test_List(t *testing.T) {
	ctx := context.Background()

	r, err := New(cfg())
	require.NoError(t, err)
	r.List(ctx)

	_ = r.createTableAuth()

	list, err := r.List(ctx)
	assert.NoError(t, err)
	logger.Info().Any("asd", list)
}

func TestSave(t *testing.T) {
	auth1 := ddd.AuthPasswd{
		Email:       email1,
		SystemRoles: []string{"admin"},
	}
	passwd1 := "asd1asd"

	ctx := context.Background()

	r, err := New(cfg())
	require.NoError(t, err)

	_ = r.createTableAuth()

	err = r.Save(ctx, auth1)
	assert.NoError(t, err, err)
	items, err := r.Get(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, "", items[0].Passwd, "expected email")
	err = r.UpdatePassword(ctx, email1, passwd1)
	assert.NoError(t, err, err)
	items, err = r.Get(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, passwd1, items[0].Passwd, "expected email")
	logger.Debug().Any("items", items).Msg("data")

	err = r.UpdateStatus(ctx, email1, ddd.StatusEnabled)
	assert.NoError(t, err, err)
	items, err = r.Get(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, ddd.StatusEnabled, items[0].Status, "expected email")

	authUpdate := ddd.AuthPasswd{
		Email:       auth1.Email,
		Status:      ddd.StatusPending,
		SystemRoles: []string{"noadmin", "shit"},
	}

	err = r.Update(ctx, authUpdate)
	assert.NoError(t, err)
	items, err = r.Get(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, authUpdate.Status, items[0].Status, "expected email")
	assert.Equal(t, authUpdate.SystemRoles, items[0].SystemRoles, "expected email")

}
