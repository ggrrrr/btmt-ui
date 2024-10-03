package awss3

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testBucket string = "test-bucket-1"

func TestListPushFetchHead(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "ok Push/Fetch",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = testBucket
				testClient := &Client{
					s3Clients: map[string]*s3Client{
						"localhost": s3c,
					},
				}

				data := "my data"
				md := &blob.BlobInfo{
					Type:          "some_type",
					ContentType:   "text/plain",
					Name:          "some-name",
					Owner:         "user-1",
					ContentLength: int64(len(data)),
				}

				newID, err := testClient.Push(ctx, "localhost", "folder-1/id-1", md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)
				defer deleteAll(ctx, s3c, awsId{folder: "folder-1", id: "id-1"})
				verID := blob.NewBlobId("folder-1", "id-1", "1")
				assert.Equal(t, &verID, newID)

				blobTest, err := testClient.Fetch(ctx, "localhost", "folder-1/id-1")
				require.NoError(t, err)
				blob.TestBlobInfoData(t, blobTest, md, data, 1000)

				blobInfo, err := testClient.Head(ctx, "localhost", "folder-1/id-1:1")
				require.NoError(t, err)
				blob.TestBlobInfo(t, blobInfo, md, 1000)

			},
		}, {
			name: "ok tenent not found",
			testFunc: func(t *testing.T) {
				testClient := &Client{
					s3Clients: map[string]*s3Client{},
				}
				_, err := testClient.Push(ctx, "notfound", "asd/asd", &blob.BlobInfo{}, nil)
				require.Error(t, err)
				// tenantNotFound := &blob.TenantNotFoundError{}
				// assert.ErrorAs(t, err, &tenantNotFound)

				_, err = testClient.Fetch(ctx, "notfound", "asd/asd")
				require.Error(t, err)
				// tenantNotFound := &blob.TenantNotFoundError{}
				// assert.ErrorAs(t, err, &tenantNotFound)

				_, err = testClient.Head(ctx, "notfound", "asd/asd")
				require.Error(t, err)
				// tenantNotFound := &blob.TenantNotFoundError{}
				// assert.ErrorAs(t, err, &tenantNotFound)
			},
		},
		{
			name: "ok Push parsing id err",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = testBucket

				testClient := &Client{
					s3Clients: map[string]*s3Client{
						"localhost": s3c,
					},
				}
				_, err = testClient.Push(ctx, "localhost", "123asd/asd", &blob.BlobInfo{}, nil)
				require.Error(t, err)
				// tenantNotFound := &blob.BlobIdInputError{}
				// assert.ErrorAsf(t, err, &tenantNotFound, "%#v, %+v", err, err)

				_, err = testClient.Fetch(ctx, "localhost", "123ad/asd")
				require.Error(t, err)
				// assert.ErrorAsf(t, err, &tenantNotFound, "%#v, %+v", err, err)

				_, err = testClient.Head(ctx, "localhost", "123ad/asd")
				require.Error(t, err)
				// assert.ErrorAsf(t, err, &tenantNotFound, "%#v, %+v", err, err)

			},
		},
		{
			name: "ok Fetch not found",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = testBucket

				testClient := &Client{
					s3Clients: map[string]*s3Client{
						"localhost": s3c,
					},
				}
				_, err = testClient.Fetch(ctx, "localhost", "ad/asd")
				require.Error(t, err)
				// notFound := &blob.NotFoundError{}
				// assert.ErrorAsf(t, err, &notFound, "%#v, %+v", err, err)

			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestListGetPutDelete(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "ok list 0",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				blockId := awsId{folder: "not", id: "no", ver: "asd"}
				result, err := list(ctx, s3c, blockId)
				require.NoError(t, err)
				assert.Equal(t, 0, len(result))
			},
		},
		{
			name: "ok list bucket not found",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = "notfound"
				blockId := awsId{folder: "not", id: "no", ver: "asd"}
				_, err = list(ctx, s3c, blockId)
				require.Error(t, err)
				// storErr := &blob.StoreNotFoundError{}
				// assert.ErrorAs(t, err, &storErr)
			},
		},
		{
			name: "ok get bucket not found",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = "notfound"
				blockId := awsId{folder: "not", id: "no", ver: "asd"}
				_, err = get(ctx, s3c, blockId)
				require.Error(t, err)
			},
		},
		{
			name: "ok get put head id.ver empty",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = testBucket
				blockId := awsId{folder: "not", id: "no", ver: ""}
				_, err = get(ctx, s3c, blockId)
				require.Error(t, err)

				_, err = head(ctx, s3c, blockId)
				require.Error(t, err)

				_, err = put(ctx, s3c, blockId, &blob.BlobInfo{}, nil)
				require.Error(t, err)
			},
		},
		{
			name: "ok put list 2",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = testBucket

				data := "mydata"

				idV1 := awsId{folder: "folder-1", id: "id-1", ver: "1"}
				idV2 := awsId{folder: "folder-1", id: "id-1", ver: "2"}
				blobInfo := &blob.BlobInfo{
					Type:          "some_type",
					ContentType:   "text/plain",
					Name:          "some-name",
					Owner:         "user1",
					ContentLength: int64(len(data)),
				}
				_, err = put(ctx, s3c, idV1, blobInfo, bytes.NewReader([]byte(data)))
				require.NoError(t, err)

				defer func() {
					err := deleteAll(ctx, s3c, idV1)
					assert.NoError(t, err)
				}()

				_, err = put(ctx, s3c, idV2, blobInfo, bytes.NewReader([]byte(data)))
				require.NoError(t, err)

				list2, err := list(ctx, s3c, idV1)
				require.NoError(t, err)
				assert.Equal(t, 2, len(list2))
				assert.Equal(t, idV1, list2[0])
				assert.Equal(t, idV2, list2[1])

			},
		},

		{
			name: "ok push 1 list 1 get 1",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg().Endpoint, cfg().Region)
				require.NoError(t, err)
				s3c.bucketName = testBucket

				data := "mydata"

				id := awsId{folder: "folder-1", id: "id-1", ver: "1"}
				blobInfo := &blob.BlobInfo{
					Type:          "some_type",
					ContentType:   "text/plain",
					Name:          "some-name",
					Owner:         "user1",
					ContentLength: int64(len(data)),
				}
				newId, err := put(ctx, s3c, id, blobInfo, bytes.NewReader([]byte(data)))
				require.NoError(t, err)

				defer func() {
					err := deleteAll(ctx, s3c, id)
					assert.NoError(t, err)
					listRes, err := list(ctx, s3c, id)
					require.NoError(t, err)
					assert.True(t, len(listRes) == 0)
				}()

				assert.Equal(t, id, newId)

				listRes, err := list(ctx, s3c, id)
				require.NoError(t, err)
				assert.Equal(t, 1, len(listRes))

				blobResult, err := get(ctx, s3c, id)
				require.NoError(t, err)
				blob.TestBlobInfoData(t, blobResult, blobInfo, data, 1000)
				assert.Equal(t, id.folder, blobResult.Id.Folder())
				assert.Equal(t, id.id, blobResult.Id.Id())
				assert.Equal(t, id.ver, blobResult.Id.Version())

				headMd, err := head(ctx, s3c, id)
				require.NoError(t, err)
				blob.TestBlobInfo(t, headMd, blobInfo, 1000)

			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestInLocal(t *testing.T) {

	utcTime := time.Date(2009, 9, 10, 13, 15, 59, 0, time.UTC)

	localTime := time.Now()

	timePrint("utcTime", utcTime)
	timePrint("utc local", utcTime.Local())
	timePrint("now", localTime)
	timePrint("now local", localTime.Local())

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "from nil",
			testFunc: func(t *testing.T) {
				resultTime := timeInLocal(nil)
				assert.WithinDuration(t, localTime, resultTime, 100*time.Millisecond)
			},
		},
		{
			name: "from UTC",
			testFunc: func(t *testing.T) {
				resultTime := timeInLocal(&utcTime)
				assert.WithinDuration(t, utcTime.Local(), resultTime, 100*time.Millisecond)
			},
		},
		{
			name: "from Local",
			testFunc: func(t *testing.T) {
				resultTime := timeInLocal(&localTime)
				assert.WithinDuration(t, localTime, resultTime, 100*time.Millisecond)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestNameRegExp(t *testing.T) {
	tests := []struct {
		fromStr   string
		excpMatch bool
	}{
		{
			fromStr:   "asd-123",
			excpMatch: true,
		},
		{
			fromStr:   "AaAsd123",
			excpMatch: true,
		},
		{
			fromStr:   "asd123",
			excpMatch: true,
		},
		{
			fromStr:   "1asd123",
			excpMatch: false,
		},
		{
			fromStr:   "/1asd123",
			excpMatch: false,
		},
		{
			fromStr:   " /1asd123",
			excpMatch: false,
		},
		{
			fromStr:   " /1asd123",
			excpMatch: false,
		},
		{
			fromStr:   "1asd123 ",
			excpMatch: false,
		},
		{
			fromStr:   "1as d123",
			excpMatch: false,
		},
		{
			fromStr:   "1asd123",
			excpMatch: false,
		},
		{
			fromStr:   `1asd\123`,
			excpMatch: false,
		},
		{
			fromStr:   `1asd/123`,
			excpMatch: false,
		},
		{
			fromStr:   "",
			excpMatch: false,
		},
		{
			fromStr:   " ",
			excpMatch: false,
		},
	}

	for _, tc := range tests {
		matched := blob.NameRegExp.MatchString(tc.fromStr)
		assert.Equalf(t, tc.excpMatch, matched, "expected/actual: %v/%v for:%s", tc.excpMatch, matched, tc.fromStr)
	}
}

func cfg() awsclient.AwsConfig {
	return awsclient.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	}
}

func timePrint(name string, t time.Time) {
	fmt.Printf("%15s: [%10s] [%5s] time: %+v \n", name, t.Location(), t.Format("-0700"), t)
}
