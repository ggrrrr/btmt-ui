package awss3

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func (c *Client) list(ctx context.Context, blobId string) ([]string, error) {
	result, err := c.s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(blobId),
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

	pipe := &blobCopier{
		reader: result.Body,
	}

	version := ""
	versionSplit := strings.Split(awsKey, "/")
	blobId := versionSplit[0]
	if len(versionSplit) == 2 {
		version = versionSplit[1]
	}

	return &blob.FetchResult{
		Metadata: blob.BlobMetadata{
			ContentType: *result.ContentType,
			CreatedAt:   *result.LastModified,
			Type:        result.Metadata["type"],
			Name:        result.Metadata["name"],
			Version:     version,
			Id:          blobId,
		},
		Writer: pipe,
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

func (c *Client) deleteBlob(ctx context.Context, blobId string) error {
	logger.InfoCtx(ctx).Str("blobId", blobId).Msg("deleteBlob")
	result, err := c.list(ctx, blobId)
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
