package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

const tableNameName string = "auth-passwd"

type (
	repo struct {
		prefix string
		svc    *awsdynamodb.DynamoDB
	}
)

var _ (ddd.AuthPasswdRepo) = (*repo)(nil)
var _ (ddd.AuthHistoryRepo) = (*repo)(nil)

func (r *repo) ListHistory(ctx context.Context, user string) (authHistory []ddd.AuthHistory, err error) {
	return nil, app.ErrTeapot
}

func (r *repo) DeleteHistory(ctx context.Context, id string) (err error) {
	return
}

func (r *repo) GetHistory(ctx context.Context, id uuid.UUID) (authHistory *ddd.AuthHistory, err error) {
	return
}

func (r *repo) SaveHistory(ctx context.Context, info roles.AuthInfo, method string) (err error) {
	return app.ErrTeapot
}

func New(cfg awsclient.AwsConfig, dbCfg awsclient.DynamodbConfig) (*repo, error) {
	// sess * session.Session,
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(cfg.Region),
		Endpoint: aws.String(cfg.Endpoint),
	})
	if err != nil {
		return nil, err
	}

	svc := awsdynamodb.New(sess)

	logger.Info().
		Str("table", tableNameName).
		Str("prefix", dbCfg.Prefix).Msg("New")
	return &repo{
		prefix: dbCfg.Prefix,
		svc:    svc,
	}, nil
}

func (r *repo) table() string {
	if r.prefix == "" {
		return tableNameName
	}
	return fmt.Sprintf("%s-%s", r.prefix, tableNameName)
}
