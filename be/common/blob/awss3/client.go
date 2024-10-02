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
	s3Clients map[string]*s3Client
}

type s3Client struct {
	s3Client   *s3.Client
	bucketName string
	region     string
}

type awsId struct {
	folder string
	id     string
	ver    string
}

func awsIdFromString(fromStr string) (awsId, error) {
	blobId, err := blob.ParseBlobId(fromStr)
	if err != nil {
		return awsId{}, err
	}
	return awsId{
		folder: blobId.Folder(),
		id:     blobId.Id(),
		ver:    blobId.Version(),
	}, nil

}

// folder/id
func (i awsId) idFolder() string {
	return fmt.Sprintf("%s/%s", i.folder, i.id)
}

// folder/id/ver
func (i awsId) keyVer() string {
	return fmt.Sprintf("%s/%s/%s", i.folder, i.id, i.ver)
}

// folder/id/ver
func (i awsId) String() string {
	return i.keyVer()
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
func (client *Client) Push(ctx context.Context, tenant string, idString string, metadata *blob.BlobMetadata, reader io.ReadSeeker) (*blob.BlockId, error) {
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

	newId, err := put(ctx, c, id, metadata, reader)
	blobId := blob.New(newId.folder, newId.id, newId.ver)
	return &blobId, err
}

func put(ctx context.Context, c *s3Client, id awsId, metadata *blob.BlobMetadata, reader io.ReadSeeker) (awsId, error) {

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
		// logger.ErrorCtx(ctx, err).Any("asd", object).Send()
		// noBucket := &types.NoSuchBucket{}
		// if errors.As(err, &noBucket) {
		// 	return awsId{}, blob.NewNotStoreError(c.bucketName, err)
		// }
		return awsId{}, fmt.Errorf("awss3.pub %w", err)
	}

	metadata.CreatedAt = time.Now()

	logger.InfoCtx(ctx).
		Str("ContentType", *object.ContentType).
		Str("Key", *object.Key).
		Msg("Pushed")

	return id, nil
}

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
		Metadata: blob.BlobMetadata{
			ContentType: *result.ContentType,
			CreatedAt:   timeInLocal(result.LastModified),
			Type:        result.Metadata["type"],
			Name:        result.Metadata["name"],
			Owner:       result.Metadata["owner"],
		},
		ReadCloser: result.Body,
	}, nil
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

func (c *Client) getClient(tenant string) (*s3Client, error) {
	s3C, ok := c.s3Clients[tenant]
	if !ok {
		return nil, blob.NewTenantNotFoundError(tenant)
	}
	return s3C, nil
}

func timeInLocal(from *time.Time) time.Time {
	if from == nil {
		return time.Now()
	}
	to := from.Local()
	return to
}
