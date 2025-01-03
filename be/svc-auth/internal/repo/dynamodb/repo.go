package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (r *repo) ListPasswd(ctx context.Context, filter app.FilterFactory) (out []ddd.AuthPasswd, err error) {
	ctx, span := logger.Span(ctx, "List", nil)
	defer func() {
		span.End(err)
	}()

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#E": aws.String("subject"),
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

func (r *repo) SavePasswd(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := logger.Span(ctx, "Save", nil)
	defer func() {
		span.End(err)
	}()

	defer func() {
		if err != nil {
			logger.Error(err).Str("subject", auth.Subject).Err(err).Msg("ops")
		}
	}()
	auth.CreatedAt = time.Now()
	if err := r.saveItem(ctx, r.table(), auth); err != nil {
		r.errorIsNotFound(err)
		return fmt.Errorf("unable to save(%s.%s): %v", r.table(), auth.Subject, err)
	}
	logger.DebugCtx(ctx).
		Str("subject", auth.Subject).
		Msg("Save")
	return nil
}

func (r *repo) GetPasswd(ctx context.Context, subject string) (out []ddd.AuthPasswd, err error) {
	ctx, span := logger.Span(ctx, "Get", nil)
	defer func() {
		span.End(err)
	}()

	input := &dynamodb.QueryInput{
		TableName: aws.String(r.table()),
	}
	input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
		":val": {
			S: aws.String(subject),
		},
	}
	input.KeyConditionExpression = aws.String("subject = :val")
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
		Str("subject", subject).
		Msg("Get")
	return
}

func (r *repo) UpdatePassword(ctx context.Context, subject string, passwd string) (err error) {
	ctx, span := logger.Span(ctx, "UpdatePassword", nil)
	defer func() {
		span.End(err)
	}()

	logger.DebugCtx(ctx).Str("subject", subject).Msg("UpdatePassword")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":passwd": {
				S: aws.String(string(passwd)),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"subject": {
				S: aws.String(string(subject)),
			},
		},
		UpdateExpression: aws.String("set #passwd = :passwd"),
		ExpressionAttributeNames: map[string]*string{
			"#passwd": aws.String("passwd"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Str("subject", subject).
		Msg("UpdatePassword")
	return err
}

func (r *repo) EnableEmail(ctx context.Context, subject string) (err error) {
	ctx, span := logger.Span(ctx, "EnableSubject", nil)
	defer func() {
		span.End(err)
	}()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":enabled": {
				BOOL: aws.Bool(true),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"subject": {
				S: aws.String(string(subject)),
			},
		},
		UpdateExpression: aws.String("set #enabled = :enabled"),
		ExpressionAttributeNames: map[string]*string{
			"#enabled": aws.String("enabled"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Str("subject", subject).
		Msg("EnableSubject")
	return
}

func (r *repo) UpdateStatus(ctx context.Context, subject string, status ddd.StatusType) (err error) {
	ctx, span := logger.Span(ctx, "UpdateStatus", nil)
	defer func() {
		span.End(err)
	}()

	logger.DebugCtx(ctx).Str("subject", subject).Str("status", string(status)).Msg("UpdateStatus")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(string(status)),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"subject": {
				S: aws.String(string(subject)),
			},
		},
		UpdateExpression: aws.String("set #status = :status"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Str("subject", subject).
		Msg("UpdateStatus")
	return
}

func toAwsMap(src map[string][]string) map[string]*dynamodb.AttributeValue {
	out := map[string]*dynamodb.AttributeValue{}
	for k, v := range src {
		out[k] = &dynamodb.AttributeValue{SS: aws.StringSlice(v)}

	}
	return out
}

func (r *repo) Update(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := logger.Span(ctx, "Update", nil)
	defer func() {
		span.End(err)
	}()

	logger.DebugCtx(ctx).Str("subject", auth.Subject).Msg("Update")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(string(auth.Status)),
			},
			":system_roles": {
				SS: aws.StringSlice(auth.SystemRoles),
			},
			":tenant_roles": {
				M: toAwsMap(auth.RealmRoles),
			},
		},
		TableName: aws.String(r.table()),
		Key: map[string]*dynamodb.AttributeValue{
			"subject": {
				S: aws.String(auth.Subject),
			},
		},
		UpdateExpression: aws.String("set #status = :status, #system_roles = :system_roles, #tenant_roles = :tenant_roles"),
		ExpressionAttributeNames: map[string]*string{
			"#status":       aws.String("status"),
			"#system_roles": aws.String("system_roles"),
			"#tenant_roles": aws.String("tenant_roles"),
		},
	}
	res, err := r.svc.UpdateItem(input)
	if err != nil {
		return
	}
	logger.DebugCtx(ctx).
		Any("res", res).
		Msg("Update")
	return nil
}

func (r *repo) errorIsNotFound(err error) {
	if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
		err := r.createTableAuth()
		if err != nil {
			logger.Error(err).Msg("createTableAuth")
		}
	}
}

func (r *repo) createTableAuth() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("subject"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("subject"), KeyType: aws.String("HASH")},
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
	logger.DebugCtx(ctx).Any("asdasd", av).Msg("asdasd")
	if _, err := r.svc.PutItem(input); err != nil {
		logger.ErrorCtx(ctx, err).Msg("asdasd")
		return err
	}
	return nil
}
