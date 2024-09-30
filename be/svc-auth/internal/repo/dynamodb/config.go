package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
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
