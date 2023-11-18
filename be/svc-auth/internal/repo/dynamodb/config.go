package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/ggrrrr/btmt-ui/be/common/awsdb"
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

func New(sess *session.Session, cfg awsdb.DynamodbConfig) (*repo, error) {

	svc := awsdynamodb.New(sess)

	logger.Info().
		Str("table", tableNameName).
		Str("prefix", cfg.Prefix).Msg("New")
	return &repo{
		prefix: cfg.Prefix,
		svc:    svc,
	}, nil
}

func (r *repo) table() string {
	if r.prefix == "" {
		return tableNameName
	}
	return fmt.Sprintf("%s-%s", r.prefix, tableNameName)
}
