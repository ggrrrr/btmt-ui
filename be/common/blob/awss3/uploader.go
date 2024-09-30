package awss3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

var _ (blob.Uploader) = (*Client)(nil)

func (u *Client) Upload(ctx context.Context, metadata *blob.BlobMetadata, reader io.ReadSeeker) error {

	objects, err := u.list(ctx, metadata.Id)
	if err != nil {
		return fmt.Errorf("unable list objects: %w", err)
	}

	versionId := (1 + len(objects))
	metadata.Version = fmt.Sprintf("%d", versionId)

	awsKey := fmt.Sprintf("%s/%d", metadata.Id, versionId)

	md := map[string]string{
		"version": metadata.Version,
		"type":    metadata.Type,
		"name":    metadata.Name,
	}

	object := &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(awsKey),
		ContentType: aws.String(metadata.ContentType),
		Metadata:    md,
		// Tagging: ,
		Body: reader,
	}

	_, err = u.s3Client.PutObject(ctx, object)
	logger.InfoCtx(ctx).
		Str("ContentType", *object.ContentType).
		Str("Key", *object.Key).
		Msg("Uploaded")

	return err
}
