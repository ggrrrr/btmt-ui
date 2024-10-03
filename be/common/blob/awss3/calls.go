package awss3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

// result is slice of aws keys
func list(ctx context.Context, c s3Client, id awsId) ([]awsId, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "awss3.list", nil, logger.AttributeString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	logger.DebugCtx(ctx).
		Str("Key", id.idFolder()).
		Msg("awss3.list")

	result, err := c.s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(id.idFolder()),
	})
	if err != nil {
		noBucket := &types.NoSuchBucket{}
		if errors.As(err, &noBucket) {
			return nil, app.SystemError("store not found", nil)
		}
		return nil, fmt.Errorf("awss3.ListObjects[%s]: %w", id.idFolder(), err)
	}

	out := make([]awsId, 0, len(result.Contents))
	for _, v := range result.Contents {
		newKey := *v.Key
		idParts := strings.Split(newKey, "/")

		item := awsId{
			folder: idParts[0],
			id:     idParts[1],
			ver:    idParts[2],
		}
		out = append(out, item)

	}
	return out, nil
}

func head(ctx context.Context, c s3Client, id awsId) (*blob.BlobInfo, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "awss3.head", nil, logger.AttributeString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	if id.ver == "" {
		err = fmt.Errorf("awss3.head: ver is empty")
		return nil, err
	}
	logger.DebugCtx(ctx).
		Str("Key", id.keyVer()).
		Msg("awss3.head")

	result, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(id.keyVer()),
	})
	if err != nil {
		err = fmt.Errorf("awss3.HeadObject[%s] %w", id.keyVer(), err)
		return nil, err
	}

	return &blob.BlobInfo{
		ContentType:   *result.ContentType,
		ContentLength: *result.ContentLength,
		CreatedAt:     timeInLocal(result.LastModified),
		Type:          result.Metadata["type"],
		Name:          result.Metadata["name"],
		Owner:         result.Metadata["owner"],
	}, nil
}

func get(ctx context.Context, c s3Client, id awsId) (*blob.FetchResult, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "awss3.get", nil, logger.AttributeString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	if id.ver == "" {
		err = fmt.Errorf("awss3.get: ver is empty")
		return nil, err
	}
	logger.DebugCtx(ctx).
		Str("Key", id.keyVer()).
		Msg("awss3.get")

	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(id.keyVer()),
	})
	if err != nil {
		// noBucket := &types.NoSuchBucket{}
		// if errors.As(err, &noBucket) {
		// return nil, blob.NewNotStoreError(c.bucketName, err)
		// }
		// noSuchKey := &types.NoSuchKey{}
		// if errors.As(err, &noSuchKey) {
		// 	return nil, blob.NewNotFoundError(id.keyVer(), err)
		// }
		err = fmt.Errorf("awss3.GetObject[%s]: %w", id.String(), err)
		return nil, err
	}

	return &blob.FetchResult{
		Id: blob.NewBlobId(id.folder, id.id, id.ver),
		Info: blob.BlobInfo{
			ContentType:   *result.ContentType,
			ContentLength: *result.ContentLength,
			CreatedAt:     timeInLocal(result.LastModified),
			Type:          result.Metadata["type"],
			Name:          result.Metadata["name"],
			Owner:         result.Metadata["owner"],
		},
		ReadCloser: result.Body,
	}, nil
}

func put(ctx context.Context, c s3Client, id awsId, metadata *blob.BlobInfo, reader io.ReadSeeker) (awsId, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "awss3.put", nil, logger.AttributeString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	if id.ver == "" {
		err = fmt.Errorf("awss3.put: ver is empty")
		return awsId{}, err
	}

	logger.DebugCtx(ctx).
		Str("Key", id.keyVer()).
		Msg("awss3.put")

	md := map[string]string{
		"type":  metadata.Type,
		"name":  metadata.Name,
		"owner": metadata.Owner,
	}

	object := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(id.keyVer()),
		ContentType: aws.String(metadata.ContentType),
		Metadata:    md,
		Body:        reader,
	}

	_, err = c.s3Client.PutObject(ctx, object)
	if err != nil {
		return awsId{}, fmt.Errorf("awss3.PutObject[%s]: %w", id.keyVer(), err)
	}

	metadata.CreatedAt = time.Now()

	return id, nil
}

func delete(ctx context.Context, c s3Client, id awsId) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "awss3.delete", nil, logger.AttributeString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	logger.DebugCtx(ctx).
		Str("key", id.String()).
		Msg("awss3.delete")

	_, err = c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    aws.String(id.keyVer()),
	})
	return err
}

func deleteAll(ctx context.Context, c s3Client, id awsId) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "awss3.deleteAll", nil, logger.AttributeString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	logger.InfoCtx(ctx).
		Str("key", id.idFolder()).
		Msg("awss3.deleteAll")

	result, err := list(ctx, c, id)
	if err != nil {
		return err
	}
	var lastErr error
	for _, v := range result {
		err = delete(ctx, c, v)
		if err != nil {
			lastErr = err
			logger.WarnCtx(ctx).Err(err).Str("aws.key[]", v.String()).Msg("awss3.deleteAll")
		}

	}
	return lastErr
}
