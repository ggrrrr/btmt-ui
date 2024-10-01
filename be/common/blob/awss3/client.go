package awss3

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

// https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html
type Client struct {
	s3Client   *s3.Client
	bucketName string
	region     string
}

func NewClient(bucketName string, appCfg awsclient.AwsConfig) (*Client, error) {
	//nolint:staticcheck
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               appCfg.Endpoint,
					SigningRegion:     appCfg.Region,
					HostnameImmutable: true,
				}, nil
			}),
		))

	if err != nil {
		return nil, fmt.Errorf("aws.session: %w", err)
	}
	// s3Client := s3.NewFromConfig(config)
	s3client := s3.NewFromConfig(cfg)
	return &Client{
		s3Client:   s3client,
		bucketName: bucketName,
		region:     cfg.Region,
	}, nil
}

var _ (blob.Fetcher) = (*Client)(nil)
var _ (blob.Uploader) = (*Client)(nil)

func (u *Client) Fetch(ctx context.Context, folder string, blobId string, version string) (*blob.FetchResult, error) {
	var err error
	blobKey, err := blobKey(folder, blobId)
	if err != nil {
		return nil, err
	}

	versions, err := u.list(ctx, blobKey)
	if err != nil {
		return nil, fmt.Errorf("unable to list folder[%s], %w", blobId, err)
	}
	if len(versions) == 0 {
		return nil, blob.NewNotFoundError(blobId, nil)
	}
	lastVersion := versions[len(versions)-1]
	object, err := u.get(ctx, lastVersion)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (u *Client) Upload(ctx context.Context, folder string, metadata *blob.BlobMetadata, reader io.ReadSeeker) error {
	blobKey, err := blobKey(folder, metadata.Id)
	if err != nil {
		return err
	}

	objects, err := u.list(ctx, blobKey)
	if err != nil {
		return fmt.Errorf("unable list objects: %w", err)
	}

	versionId := (1 + len(objects))
	metadata.Version = fmt.Sprintf("%d", versionId)

	awsKey := fmt.Sprintf("%s/%d", blobKey, versionId)

	md := map[string]string{
		"version": metadata.Version,
		"type":    metadata.Type,
		"name":    metadata.Name,
		"owner":   metadata.Owner,
	}

	object := &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(awsKey),
		ContentType: aws.String(metadata.ContentType),
		Metadata:    md,
		Body:        reader,
	}

	_, err = u.s3Client.PutObject(ctx, object)
	metadata.CreatedAt = time.Now()
	logger.InfoCtx(ctx).
		Str("ContentType", *object.ContentType).
		Str("Key", *object.Key).
		Msg("Uploaded")

	return err
}

// result is slice of aws keys
func (c *Client) list(ctx context.Context, blobKey string) ([]string, error) {

	result, err := c.s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(blobKey),
	})
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, v := range result.Contents {
		out = append(out, *v.Key)
	}
	return out, nil
}

func (c *Client) get(ctx context.Context, awsKey string) (*blob.FetchResult, error) {
	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(awsKey),
	})
	if err != nil {
		noSuchKey := &types.NoSuchKey{}
		if errors.As(err, &noSuchKey) {
			return nil, blob.NewNotFoundError(awsKey, err)
		}
		return nil, err
	}

	version := ""
	versionSplit := strings.Split(awsKey, "/")
	if len(versionSplit) > 1 {
		version = versionSplit[len(versionSplit)-1]
	}
	blobId := strings.Join(versionSplit[1:len(versionSplit)-1], "/")
	return &blob.FetchResult{
		Metadata: blob.BlobMetadata{
			Id:          blobId,
			ContentType: *result.ContentType,
			CreatedAt:   timeInLocal(result.LastModified),
			Type:        result.Metadata["type"],
			Name:        result.Metadata["name"],
			Owner:       result.Metadata["owner"],
			Version:     version,
		},
		ReadCloser: result.Body,
	}, nil
}

func (c *Client) delete(ctx context.Context, awsId string) error {
	logger.InfoCtx(ctx).Str("awsId", awsId).Msg("delete")
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    aws.String(awsId),
	})
	return err
}

func (c *Client) deleteBlob(ctx context.Context, folder string, blobId string) error {
	blobKey, err := blobKey(folder, blobId)
	if err != nil {
		return err
	}
	logger.InfoCtx(ctx).Str("blobId", blobId).Msg("deleteBlob")
	result, err := c.list(ctx, blobKey)
	if err != nil {
		return err
	}
	var lastErr error
	for _, v := range result {
		err = c.delete(ctx, v)
		if err != nil {
			lastErr = err
			logger.WarnCtx(ctx).Err(err).Str("aws.key", v).Msg("delete.aws.object")
		}

	}
	return lastErr
}

func blobKey(folder, blobId string) (string, error) {
	ok := blob.NameRegExp.MatchString(folder)
	if !ok {
		return "", &blob.FolderInputError{}
	}
	ok = blob.NameRegExp.MatchString(blobId)
	if !ok {
		return "", &blob.IdInputError{}
	}
	return fmt.Sprintf("%s/%s", folder, blobId), nil
}

func timeInLocal(from *time.Time) time.Time {
	if from == nil {
		return time.Now()
	}
	to := from.Local()
	return to
}
