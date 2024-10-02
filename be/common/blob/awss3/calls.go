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
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

// result is slice of aws keys
func list(ctx context.Context, c *s3Client, id awsId) ([]awsId, error) {

	result, err := c.s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(id.idFolder()),
	})
	if err != nil {
		noBucket := &types.NoSuchBucket{}
		if errors.As(err, &noBucket) {
			return nil, blob.NewNotStoreError(c.bucketName, err)
		}
		return nil, fmt.Errorf("aws.ListObjects %w", err)
	}
	out := []awsId{}
	for _, v := range result.Contents {

		ver := strings.Replace(*v.Key, fmt.Sprintf("%s/%s/", id.folder, id.id), "", 1)
		out = append(out, awsId{
			folder: id.folder,
			id:     id.id,
			ver:    ver,
		})
	}
	return out, nil
}

func head(ctx context.Context, c *s3Client, id awsId) (*blob.BlobInfo, error) {
	result, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(id.keyVer()),
	})
	if err != nil {
		return nil, fmt.Errorf("awss3.HeadObject %w", err)
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

func get(ctx context.Context, c *s3Client, id awsId) (*blob.FetchResult, error) {
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
		return nil, fmt.Errorf("aws.GetObject[%s] %w", id.String(), err)
	}

	return &blob.FetchResult{
		Id: blob.New(id.folder, id.id, id.ver),
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

func put(ctx context.Context, c *s3Client, id awsId, metadata *blob.BlobInfo, reader io.ReadSeeker) (awsId, error) {

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

	_, err := c.s3Client.PutObject(ctx, object)
	if err != nil {
		return awsId{}, fmt.Errorf("awss3.pub %w", err)
	}

	metadata.CreatedAt = time.Now()

	logger.InfoCtx(ctx).
		Str("ContentType", *object.ContentType).
		Str("Key", *object.Key).
		Msg("Pushed")

	return id, nil
}

func delete(ctx context.Context, c *s3Client, id awsId) error {
	logger.InfoCtx(ctx).Str("awsId", id.String()).Msg("delete")
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    aws.String(id.keyVer()),
	})
	return err
}

func deleteAll(ctx context.Context, c *s3Client, id awsId) error {
	logger.InfoCtx(ctx).Str("blobId", id.idFolder()).Msg("deleteBlob")
	result, err := list(ctx, c, id)
	if err != nil {
		return err
	}
	var lastErr error
	for _, v := range result {
		err = delete(ctx, c, v)
		if err != nil {
			lastErr = err
			logger.WarnCtx(ctx).Err(err).Str("aws.key", v.String()).Msg("delete.aws.object")
		}

	}
	return lastErr
}
