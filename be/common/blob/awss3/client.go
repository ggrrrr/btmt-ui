package awss3

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

// https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html
type Client struct {
	s3Clients map[string]*s3Client
}

type s3Client struct {
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
		return nil, fmt.Errorf("aws.config: %w", err)
	}
	// s3Client := s3.NewFromConfig(config)
	s3client := s3.NewFromConfig(cfg)

	res, err := s3client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, blob.NewNotStoreError(bucketName, err)
	}
	logger.Info().
		Any("bucket", bucketName).
		Any("info", res).
		Msg("bucket")

	return &Client{
		s3Clients: map[string]*s3Client{
			"localhost": {
				s3Client:   s3client,
				bucketName: bucketName,
				region:     cfg.Region,
			},
		},
	}, nil
}

var _ (blob.Fetcher) = (*Client)(nil)
var _ (blob.Pusher) = (*Client)(nil)

func (client *Client) Fetch(ctx context.Context, tenant string, idString string) (*blob.FetchResult, error) {
	var err error

	c, err := client.getClient(tenant)
	if err != nil {
		return nil, err
	}

	id, err := awsIdFromString(idString)
	if err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).Str("id", id.String()).Msg("Fetch")
	versions, err := list(ctx, c, id)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, blob.NewNotFoundError(id.idFolder(), nil)
	}
	lastVersion := versions[len(versions)-1]
	if id.ver != "" {
		if lastVersion.ver != id.ver {
			return nil, blob.NewNotFoundError(id.String(), fmt.Errorf("version not found"))
		}
	}

	object, err := get(ctx, c, id)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (client *Client) Head(ctx context.Context, tenant string, idString string) (*blob.BlobInfo, error) {
	var err error

	c, err := client.getClient(tenant)
	if err != nil {
		return nil, err
	}

	id, err := awsIdFromString(idString)
	if err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).Str("id", id.String()).Msg("Fetch")
	versions, err := list(ctx, c, id)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, blob.NewNotFoundError(id.idFolder(), nil)
	}
	lastVersion := versions[len(versions)-1]
	if id.ver != "" {
		if lastVersion.ver != id.ver {
			return nil, blob.NewNotFoundError(id.String(), fmt.Errorf("version not found"))
		}
	}
	md, err := head(ctx, c, id)
	if err != nil {
		return nil, err
	}
	return md, nil
}

func (client *Client) Push(ctx context.Context, tenant string, idString string, blobInfo *blob.BlobInfo, reader io.ReadSeeker) (*blob.BlobId, error) {
	c, err := client.getClient(tenant)
	if err != nil {
		return nil, err
	}

	id, err := awsIdFromString(idString)
	if err != nil {
		return nil, err
	}

	objects, err := list(ctx, c, id)
	if err != nil {
		return nil, err
	}

	// TODO: add error if we need to overwrite
	newVer := fmt.Sprintf("%d", (1 + len(objects)))

	id.ver = newVer

	newId, err := put(ctx, c, id, blobInfo, reader)
	blobId := blob.NewBlobId(newId.folder, newId.id, newId.ver)
	return &blobId, err
}

func (c *Client) getClient(tenant string) (*s3Client, error) {
	s3C, ok := c.s3Clients[tenant]
	if !ok {
		return nil, blob.NewTenantNotFoundError(tenant)
	}
	return s3C, nil
}

func createS3Client(endpoint, region string) (*s3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithHTTPClient(
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			}),
		))
	if err != nil {
		return nil, fmt.Errorf("aws.config: %w", err)
	}
	return &s3Client{
		s3Client: s3.NewFromConfig(cfg),
	}, nil
}
