package awss3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

// result is slice of aws keys
func list(ctx context.Context, c s3Client, id awsId) ([]awsId, error) {
	var err error
	ctx, span := c.otelTracer.SpanWithAttributes(ctx, "awss3.list", slog.String("id", id.String()))
	defer func() {
		span.End(err)
	}()

	log.Log().DebugCtx(ctx, "awss3.lit", log.WithString("key", id.pathId()))

	result, err := c.s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(id.pathId()),
	})
	if err != nil {
		noBucket := &types.NoSuchBucket{}
		if errors.As(err, &noBucket) {
			return nil, app.SystemError("store not found", nil)
		}
		return nil, fmt.Errorf("awss3.ListObjects[%s]: %w", id.pathId(), err)
	}

	out := make([]awsId, 0, len(result.Contents))
	for _, v := range result.Contents {
		newKey := *v.Key
		idParts := strings.Split(newKey, "/")
		partsLen := len(idParts)
		item := awsId{
			createdAt: *v.LastModified,
		}
		switch partsLen {
		case 0, 1, 2:
			log.Log().ErrorCtx(ctx, fmt.Errorf("aws key len"), "awss3.lit ignored",
				log.WithString("key", newKey),
				log.WithInt("len", partsLen),
				log.WithString("id", id.String()),
			)
			continue
		default:
			item.path = strings.Join(idParts[:partsLen-2], "/")
			item.id = idParts[partsLen-2]
			item.ver = idParts[partsLen-1]
		}
		out = append(out, item)
	}
	// sort.Slice(out, func(i, j int) bool {
	// fromI := out[i]
	// fromJ := out[j]
	// return fromI.createdAt.After(fromJ.createdAt)
	// })

	return out, nil
}

func head(ctx context.Context, c s3Client, id awsId) (blob.BlobMD, error) {
	var err error
	ctx, span := c.otelTracer.SpanWithAttributes(ctx, "awss3.head", slog.String("id", id.String()))
	defer func() {
		span.End(err)
	}()

	if id.ver == "" {
		err = fmt.Errorf("awss3.head: ver is empty")
		return blob.BlobMD{}, err
	}
	log.Log().DebugCtx(ctx, "awss3.head",
		log.WithString("id", id.String()),
	)

	result, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(id.keyVer()),
	})
	if err != nil {
		err = fmt.Errorf("awss3.HeadObject[%s] %w", id.keyVer(), err)
		return blob.BlobMD{}, err
	}

	return fromAwsHeadToBlobInfo(result), nil
}

func get(ctx context.Context, c s3Client, id awsId) (blob.BlobReader, error) {
	var err error
	ctx, span := c.otelTracer.SpanWithAttributes(ctx, "awss3.get", slog.String("id", id.String()))
	defer func() {
		span.End(err)
	}()

	if id.ver == "" {
		err = fmt.Errorf("awss3.get: ver is empty")
		return blob.BlobReader{}, err
	}
	log.Log().DebugCtx(ctx, "awss3.get",
		log.WithString("id", id.String()),
	)

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
		return blob.BlobReader{}, err
	}

	blobMD := fromAwsGetToBlobInfo(result)

	blobId, err := blob.NewBlobId(id.path, id.id, id.ver)
	if err != nil {
		err = fmt.Errorf("awss3.NewBlobId[%s]: %w", id.String(), err)
		return blob.BlobReader{}, err
	}

	return blob.BlobReader{

		Blob: blob.Blob{
			Id: blobId,
			MD: blobMD,
		},
		ReadCloser: result.Body,
	}, nil
}

func put(ctx context.Context, c s3Client, id awsId, metadata blob.BlobMD, reader io.ReadSeeker) (awsId, error) {
	var err error
	ctx, span := c.otelTracer.SpanWithAttributes(ctx, "awss3.put", slog.String("id", id.String()))
	defer func() {
		span.End(err)
	}()

	if id.ver == "" {
		err = fmt.Errorf("awss3.put: ver is empty")
		return awsId{}, err
	}

	log.Log().DebugCtx(ctx, "awss3.put",
		log.WithString("id", id.String()),
	)

	object := fromBlobInfoToAwsObject(c, id, metadata)
	object.Body = reader

	_, err = c.s3Client.PutObject(ctx, object)
	if err != nil {
		return awsId{}, fmt.Errorf("awss3.PutObject[%s]: %w", id.keyVer(), err)
	}

	return id, nil
}

func delete(ctx context.Context, c s3Client, id awsId) error {
	var err error
	ctx, span := c.otelTracer.SpanWithAttributes(ctx, "awss3.delete", slog.String("id", id.String()))
	defer func() {
		span.End(err)
	}()

	log.Log().DebugCtx(ctx, "awss3.delete",
		log.WithString("id", id.String()),
	)

	_, err = c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    aws.String(id.keyVer()),
	})
	return err
}

func deleteAll(ctx context.Context, c s3Client, id awsId) error {
	var err error
	ctx, span := c.otelTracer.SpanWithAttributes(ctx, "awss3.deleteAll", slog.String("id", id.String()))
	defer func() {
		span.End(err)
	}()

	log.Log().DebugCtx(ctx, "awss3.deleteAll",
		log.WithString("id", id.String()),
	)

	result, err := list(ctx, c, id)
	if err != nil {
		return err
	}
	var lastErr error
	for _, v := range result {
		err = delete(ctx, c, v)
		if err != nil {
			lastErr = err
			log.Log().WarnCtx(ctx, err, "awss3.deleteAll",
				log.WithString("id.key", v.String()),
			)
		}

	}
	return lastErr
}
