package awss3

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.common.blob.awss3"

// https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html
type (
	realmKey string

	Client struct {
		otelTracer tracer.OTelTracer
		s3Clients  map[realmKey]s3Client
	}
)
type s3Client struct {
	otelTracer tracer.OTelTracer
	s3Client   *s3.Client
	bucketName string
	region     string
}

var _ (blob.Fetcher) = (*Client)(nil)
var _ (blob.Pusher) = (*Client)(nil)

func NewClient(bucketName string, appCfg awsclient.Config) (*Client, error) {
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
					URL:               appCfg.EndpointURL,
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
		return nil, app.SystemError("aws bucket not found", err)
	}
	log.Log().Info("bucket",
		log.WithString("bucket", bucketName),
		log.WithAny("info", res),
	)

	return &Client{
		otelTracer: tracer.Tracer(otelScope),
		s3Clients: map[realmKey]s3Client{
			realmKey("localhost"): {
				otelTracer: tracer.Tracer(otelScope),
				s3Client:   s3client,
				bucketName: bucketName,
				region:     cfg.Region,
			},
		},
	}, nil
}

func (client *Client) Fetch(ctx context.Context, tenant string, blobId blob.BlobId) (blob.BlobReader, error) {
	var err error
	ctx, span := client.otelTracer.SpanWithAttributes(ctx, "awss3.Fetch", slog.String("id", blobId.String()))
	defer func() {
		span.End(err)
	}()

	c, err := client.getClient(tenant)
	if err != nil {
		return blob.BlobReader{}, err
	}

	id := awsIdFromString(blobId)

	log.Log().InfoCtx(ctx, "awss3.Fetch",
		log.WithString("id", id.String()),
	)

	versions, err := list(ctx, c, id)
	if err != nil {
		return blob.BlobReader{}, err
	}
	if len(versions) == 0 {
		return blob.BlobReader{}, app.ItemNotFoundError("blob", id.pathId())
	}

	lastVersion := versions[len(versions)-1]
	if id.ver != "" {
		if lastVersion.ver != id.ver {
			return blob.BlobReader{}, app.ItemNotFoundError("blob.version", id.String())
		}
	}

	object, err := get(ctx, c, lastVersion)
	if err != nil {
		return blob.BlobReader{}, err
	}
	return object, nil
}

func (client *Client) Head(ctx context.Context, tenant string, blobId blob.BlobId) (blob.BlobMD, error) {
	var err error
	ctx, span := client.otelTracer.SpanWithAttributes(ctx, "awss3.Head", slog.String("id", blobId.String()))
	defer func() {
		span.End(err)
	}()

	c, err := client.getClient(tenant)
	if err != nil {
		return blob.BlobMD{}, err
	}

	id := awsIdFromString(blobId)
	if err != nil {
		return blob.BlobMD{}, err
	}

	log.Log().InfoCtx(ctx, "awss3.Head",
		log.WithString("id", id.String()),
	)

	versions, err := list(ctx, c, id)
	if err != nil {
		return blob.BlobMD{}, err
	}
	if len(versions) == 0 {
		return blob.BlobMD{}, app.ItemNotFoundError("blob", id.pathId())
	}

	lastVersion := versions[len(versions)-1]
	// TODO check all version in the list
	if id.ver != "" {
		if lastVersion.ver != id.ver {
			return blob.BlobMD{}, app.ItemNotFoundError("blob version", id.keyVer())
		}
	}

	md, err := head(ctx, c, lastVersion)
	if err != nil {
		return blob.BlobMD{}, err
	}
	return md, nil
}

func (client *Client) List(ctx context.Context, tenant string, blobId blob.BlobId) ([]blob.ListResult, error) {
	var err error
	ctx, span := client.otelTracer.SpanWithAttributes(ctx, "awss3.ListDir", slog.String("id", blobId.String()))
	defer func() {
		span.End(err)
	}()

	c, err := client.getClient(tenant)
	if err != nil {
		return nil, err
	}

	id := awsIdFromString(blobId)
	if err != nil {
		return nil, err
	}

	log.Log().InfoCtx(ctx, "awss3.ListDir",
		log.WithString("id", id.String()),
	)

	blobs, err := list(ctx, c, id)
	if err != nil {
		return nil, err
	}
	if len(blobs) == 0 {
		return nil, app.ItemNotFoundError("blob", id.pathId())
	}

	sort.Slice(blobs, func(i, j int) bool {
		return blobs[i].createdAt.After(blobs[j].createdAt)
	})

	tempMap := map[string]*blob.ListResult{}

	// TODO consider multi threading here, to speedup long lists
	for _, awsId := range blobs {
		blobMD, err := head(ctx, c, awsId)
		if err != nil {
			log.Log().ErrorCtx(ctx, err, "awss3.head",
				log.WithString("id", id.String()),
			)

			continue
		}

		blobId, err := blob.NewBlobId(awsId.path, awsId.id, awsId.ver)
		if err != nil {
			return nil, fmt.Errorf("blob.NewBlobId[%s] %w", awsId.String(), err)
		}

		currentBlob := blob.Blob{
			Id: blobId,
			MD: blobMD,
		}

		currentResult := &blob.ListResult{
			Blob: currentBlob,
		}

		lastResult, ok := tempMap[awsId.id]
		if ok {
			lastResult.Versions = append(lastResult.Versions, currentBlob)
			continue
		}

		tempMap[awsId.id] = currentResult

	}

	out := []blob.ListResult{}
	for k := range tempMap {
		out = append(out, *tempMap[k])
	}
	return out, nil
}

func (client *Client) Push(ctx context.Context, tenant string, blobId blob.BlobId, blobInfo blob.BlobMD, reader io.ReadSeeker) (blob.BlobId, error) {
	var err error
	ctx, span := client.otelTracer.SpanWithAttributes(ctx, "awss3.Push", slog.String("id", blobId.String()))
	defer func() {
		span.End(err)
	}()

	if reader == nil {
		err = fmt.Errorf("reader is nil")
		return blob.BlobId{}, err
	}

	c, err := client.getClient(tenant)
	if err != nil {
		return blob.BlobId{}, err
	}

	id := awsIdFromString(blobId)
	if err != nil {
		return blob.BlobId{}, err
	}

	log.Log().InfoCtx(ctx, "awss3.Push",
		log.WithString("id", id.String()),
	)

	objects, err := list(ctx, c, id)
	if err != nil {
		return blob.BlobId{}, err
	}

	// TODO: add error if we need to overwrite
	newVer := fmt.Sprintf("%d", (1 + len(objects)))

	id.ver = newVer

	newId, err := put(ctx, c, id, blobInfo, reader)
	if err != nil {
		return blob.BlobId{}, err
	}

	newBlobId, err := blob.NewBlobId(newId.path, newId.id, newId.ver)
	if err != nil {
		err = fmt.Errorf("awss3.NewBlobId[%s]: %w", id.String(), err)
		return blob.BlobId{}, err
	}

	return newBlobId, err
}

func (c *Client) getClient(realm string) (s3Client, error) {
	s3C, ok := c.s3Clients[realmKey(realm)]
	if !ok {
		return s3Client{}, app.SystemError("tenant not configured", nil)
	}
	return s3C, nil
}

func createS3Client(awsCfg awsclient.Config, appCfg awsclient.S3Client) (s3Client, error) {
	clientCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(appCfg.Region),
		config.WithHTTPClient(
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               awsCfg.EndpointURL,
					SigningRegion:     appCfg.Region,
					HostnameImmutable: true,
				}, nil
			}),
		))
	if err != nil {
		return s3Client{}, fmt.Errorf("aws.config: %w", err)
	}
	return s3Client{
		otelTracer: tracer.Tracer(otelScope),
		s3Client:   s3.NewFromConfig(clientCfg),
	}, nil
}
