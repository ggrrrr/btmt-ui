package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/ggrrrr/btmt-ui/be/common/awsdb"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	email1 string = "ggrrrr1@gmail.com"
)

func Test_List(t *testing.T) {
	ctx := context.Background()
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	require.NoError(t, err)

	r, err := New(sess, awsdb.DynamodbConfig{
		Prefix: "test",
	})
	require.NoError(t, err)
	r.List(ctx)

	_ = r.createTableAuth()

	list, err := r.List(ctx)
	assert.NoError(t, err)
	logger.Log().Info().Any("asd", list)
}

func TestSave(t *testing.T) {
	auth1 := ddd.AuthPasswd{
		Email:       email1,
		SystemRoles: []string{"admin"},
	}
	passwd1 := "asd1asd"

	ctx := context.Background()
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	require.NoError(t, err)

	// svc := dynamodb.New(sess)

	r, err := New(sess, awsdb.DynamodbConfig{
		Prefix: "test",
	})
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
	logger.Log().Debug().Any("items", items).Msg("data")

	err = r.UpdateStatus(ctx, email1, ddd.StatusEnabled)
	assert.NoError(t, err, err)
	items, err = r.Get(ctx, email1)
	assert.NoError(t, err, err)
	assert.Equal(t, 1, len(items), "expected email")
	assert.Equal(t, ddd.StatusEnabled, items[0].Status, "expected email")

}
