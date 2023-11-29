package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (r *repo) List(ctx context.Context) (out []ddd.AuthPasswd, err error) {
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#E": aws.String("email"),
			"#S": aws.String("status"),
			"#T": aws.String("created_at"),
			"#L": aws.String("system_roles"),
		},
		ProjectionExpression: aws.String("#E, #S, #T, #L"),
		TableName:            aws.String(r.table()),
	}
	defer func() {
		if err != nil {
			logger.Error(err).Msg("ops")
		}
	}()

	result, err := r.svc.Scan(input)
	if err != nil {
		r.errorIsNotFound(err)
		return nil, fmt.Errorf("unable to list items %v", err)
	}
	for _, item := range result.Items {
		auth := ddd.AuthPasswd{}

		err = dynamodbattribute.UnmarshalMap(item, &auth)
		if err != nil {
			logger.Error(err).Any("item", item).Msg("UnmarshalMap")
			continue
		}
		out = append(out, auth)
	}
	logger.DebugCtx(ctx).Msg("List")
	return
}

func (r *repo) Save(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err).Str("email", auth.Email).Err(err).Msg("ops")
		}
	}()
	auth.CreatedAt = time.Now()
	if err := r.saveItem(ctx, r.table(), auth); err != nil {
		r.errorIsNotFound(err)
		return fmt.Errorf("unable to save(%s.%s): %v", r.table(), auth.Email, err)
	}
	logger.DebugCtx(ctx).
		Str("email", auth.Email).
		Msg("Save")
	return nil
}

func (r *repo) Get(ctx context.Context, email string) (out []ddd.AuthPasswd, err error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.table()),
	}
	input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
		":val": {
			S: aws.String(email),
		},
	}
	input.KeyConditionExpression = aws.String("email = :val")
	result, err := r.svc.Query(input)
	if err != nil {
		r.errorIsNotFound(err)
		logger.Error(err).Msg("Get.Query")
		return
	}
	for _, item := range result.Items {
		auth := ddd.AuthPasswd{}

		err = dynamodbattribute.UnmarshalMap(item, &auth)
		if err != nil {
			logger.ErrorCtx(ctx, err).Any("item", item).Msg("UnmarshalMap")
			continue
		}
		out = append(out, auth)
	}
	logger.DebugCtx(ctx).
		Str("email", email).
		Msg("Get")
	return
}

func (r *repo) UpdatePassword(ctx context.Context, email string, passwd string) error {
	logger.DebugCtx(ctx).Str("email", email).Msg("UpdatePassword")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":passwd": {
				S: aws.String(string(passwd)),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(string(email)),
			},
		},
		UpdateExpression: aws.String("set #passwd = :passwd"),
		ExpressionAttributeNames: map[string]*string{
			"#passwd": aws.String("passwd"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return err
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Str("email", email).
		Msg("UpdatePassword")
	return err
}

func (r *repo) EnableEmail(ctx context.Context, email string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":enabled": {
				BOOL: aws.Bool(true),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(string(email)),
			},
		},
		UpdateExpression: aws.String("set #enabled = :enabled"),
		ExpressionAttributeNames: map[string]*string{
			"#enabled": aws.String("enabled"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return err
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Str("email", email).
		Msg("EnableEmail")
	return err
}

func (r *repo) UpdateStatus(ctx context.Context, email string, status ddd.StatusType) error {
	logger.DebugCtx(ctx).Str("email", email).Str("status", string(status)).Msg("UpdateStatus")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(string(status)),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(string(email)),
			},
		},
		UpdateExpression: aws.String("set #status = :status"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return err
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Str("email", email).
		Msg("UpdateStatus")
	return err
}

func (r *repo) errorIsNotFound(err error) {
	if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
		r.createTableAuth()
	}
}

func (r *repo) createTableAuth() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("email"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("email"), KeyType: aws.String("HASH")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(r.table()),
	}

	if _, err := r.svc.CreateTable(input); err != nil {
		logger.Error(err).Msg("createTableAuth")
		return err
	}
	logger.Warn().Str("table", r.table()).Msg("createTableAuth")
	return nil
}

func (r *repo) saveItem(ctx context.Context, table string, item interface{}) (err error) {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}
	if _, err := r.svc.PutItem(input); err != nil {
		return err
	}
	return nil
}
