package awss3

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetcher(t *testing.T) {
	testBucketName := "test-bucket-1"
	ctx := context.Background()

	tests := []struct {
		name     string
		prepFunc func(t *testing.T)
	}{
		{
			name: "not found",
			prepFunc: func(t *testing.T) {
				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				res, err := client.Fetch(ctx, "folder", "not-fond-2", "")
				require.Error(t, err)
				notFount := &blob.NotFoundError{}
				require.ErrorAs(t, err, &notFount)
				require.Nil(t, res)
			},
		},
		{
			name: "blobKey folder error",
			prepFunc: func(t *testing.T) {
				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				res, err := client.Fetch(ctx, "1folder", "not-fond-2", "")
				require.Error(t, err)
				notFount := blob.FolderInputError{}
				require.ErrorIs(t, err, &notFount)
				require.Nil(t, res)
			},
		},
		{
			name: "blobKey id error",
			prepFunc: func(t *testing.T) {
				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				res, err := client.Fetch(ctx, "folder", " not-fond-2", "")
				require.Error(t, err)
				notFount := blob.IdInputError{}
				require.ErrorIs(t, err, &notFount)
				require.Nil(t, res)
			},
		},
		{
			name: "not upload ver1",
			prepFunc: func(t *testing.T) {
				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)

				testData := "testData"
				data := &blob.BlobMetadata{
					Type:        "me",
					ContentType: "no",
					Name:        "name",
					Version:     "",
					Id:          "upload-1",
					Owner:       "some user",
				}
				err = client.Upload(ctx, "folder", data, bytes.NewReader([]byte(testData)))
				require.NoError(t, err)
				assert.Equal(t, data.Version, "1")

				defer func() {
					err := client.deleteBlob(ctx, "folder", "upload-1")
					assert.NoError(t, err)
				}()

				testData1 := "testData1"
				data1 := &blob.BlobMetadata{
					Owner:       "owner",
					Type:        "me",
					ContentType: "no",
					Name:        "name",
					Version:     "",
					Id:          "upload-1",
					CreatedAt:   time.Now(),
				}
				err = client.Upload(ctx, "folder", data1, bytes.NewReader([]byte(testData1)))
				require.NoError(t, err)
				assert.Equal(t, data1.Version, "2")

				res, err := client.Fetch(ctx, "folder", "upload-1", "")
				require.NoError(t, err)
				require.NotNil(t, res)

				blob.TestBlob(t, res, data1, testData1, 100)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFunc(t)
		})
	}
}

func TestList(t *testing.T) {
	testBucketName := "test-bucket-1"
	ctx := context.Background()

	tests := []struct {
		name     string
		prepFunc func(t *testing.T)
	}{
		{
			name: "err list",
			prepFunc: func(t *testing.T) {
				client, err := NewClient("testBucketName", cfg())
				require.NoError(t, err)
				res, err := client.list(ctx, "folder/not-fond-2")
				assert.Error(t, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "ok no list",
			prepFunc: func(t *testing.T) {
				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				res, err := client.list(ctx, "folder/not-fond-2")
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, len(res), 0)
			},
		},
		{
			name: "ok list with 1rec",
			prepFunc: func(t *testing.T) {

				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				testData := "testData"
				err = client.Upload(ctx, "folder", &blob.BlobMetadata{
					Type:        "me",
					ContentType: "no",
					Name:        "name",
					Version:     "",
					Id:          "id-1",
				}, bytes.NewReader([]byte(testData)))
				require.NoError(t, err)

				res, err := client.list(ctx, "folder/id-1")
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, len(res), 1)

				err = client.deleteBlob(ctx, "folder", "id-1")
				require.NoError(t, err)

			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFunc(t)
		})
	}
}

func TestUploadGet(t *testing.T) {
	testBucketName := "test-bucket-1"
	ctx := context.Background()

	tests := []struct {
		name     string
		awsId    string
		prepFunc func(t *testing.T, data *blob.BlobMetadata) *Client
		exptErr  error
		exptBlob *blob.BlobMetadata
	}{
		{
			name:  "not found key",
			awsId: "not-found",
			prepFunc: func(t *testing.T, _ *blob.BlobMetadata) *Client {
				s3Client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				return s3Client
			},
			exptErr:  &blob.NotFoundError{},
			exptBlob: nil,
		},
		{
			name:  "not bucket",
			awsId: "not-found",
			prepFunc: func(t *testing.T, _ *blob.BlobMetadata) *Client {
				s3Client, err := NewClient("testBucketName", cfg())
				require.NoError(t, err)
				return s3Client
			},
			// exptErr:  &blob.StoreNotFoundError{},
			exptErr:  &blob.NotFoundError{},
			exptBlob: nil,
		},
		{
			name:  "upload get okasdasdasd",
			awsId: "folder/some-id-1/1",
			prepFunc: func(t *testing.T, data *blob.BlobMetadata) *Client {
				s3Client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)

				testData := "testData"
				err = s3Client.Upload(ctx, "folder", data, bytes.NewReader([]byte(testData)))
				require.NoError(t, err)

				return s3Client
			},
			exptErr: nil,
			exptBlob: &blob.BlobMetadata{
				Id:          "some-id-1",
				ContentType: "text/plain",
				Name:        "blob-name",
				Type:        "my-type",
				Version:     "1",
				// CreatedAt:   time.Now(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := tc.prepFunc(t, tc.exptBlob)
			resp, err := client.get(ctx, tc.awsId)
			if tc.exptErr != nil {
				// fmt.Printf("---- %T %v\n", tc.exptErr, tc.exptErr)
				// fmt.Printf("---- %T %v\n", err, err)

				require.Error(t, err)
				// assert.ErrorAs(t, err, &tc.exptErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				// fmt.Printf("-------: %+v\n", resp)
				blob.TestBlob(t, resp, tc.exptBlob, "testData", 100)

				defer func() {
					asd := strings.Split(tc.awsId, "/")
					client.delete(ctx, tc.awsId)
					client.deleteBlob(ctx, "folder", asd[1])
				}()

			}
		})

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

func TestBlobKey(t *testing.T) {
	tests := []struct {
		name   string
		folder string
		id     string
		blobId string
		err    error
	}{
		{
			name:   "ok",
			folder: "folder-1",
			id:     "some-id",
			blobId: "folder-1/some-id",
			err:    nil,
		},
		{
			name:   "bad folder",
			folder: "1folder-1",
			id:     "some-id",
			blobId: "",
			err:    &blob.FolderInputError{},
		},
		{
			name:   "bad id",
			folder: "folder-1",
			id:     " some-id",
			blobId: "",
			err:    &blob.IdInputError{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blobKey, err := blobKey(tc.folder, tc.id)
			if tc.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.Equalf(t, tc.blobId, blobKey, "from folder:%s id:%s -> result: %s != %s", tc.folder, tc.id, tc.id, blobKey)
			}
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
