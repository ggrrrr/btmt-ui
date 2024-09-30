package awss3

import (
	"bytes"
	"context"
	"testing"
	"time"

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
				res, err := client.Fetch(ctx, "not-fond-2", "")
				require.Error(t, err)
				notFount := &blob.NotFoundError{}
				require.ErrorAs(t, err, &notFount)
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
				}
				err = client.Upload(ctx, data, bytes.NewReader([]byte(testData)))
				require.NoError(t, err)
				assert.Equal(t, data.Version, "1")

				defer func() {
					err := client.deleteBlob(ctx, "upload-1")
					assert.NoError(t, err)
				}()

				testData1 := "testData1"
				data1 := &blob.BlobMetadata{
					Type:        "me",
					ContentType: "no",
					Name:        "name",
					Version:     "",
					Id:          "upload-1",
				}
				err = client.Upload(ctx, data1, bytes.NewReader([]byte(testData1)))
				require.NoError(t, err)
				assert.Equal(t, data1.Version, "2")

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
			name: "ok no list",
			prepFunc: func(t *testing.T) {
				client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)
				res, err := client.list(ctx, "not-fond-2")
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
				err = client.Upload(ctx, &blob.BlobMetadata{
					Type:        "me",
					ContentType: "no",
					Name:        "name",
					Version:     "",
					Id:          "id-1",
				}, bytes.NewReader([]byte(testData)))
				require.NoError(t, err)

				res, err := client.list(ctx, "id-1")
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, len(res), 1)

				err = client.deleteBlob(ctx, "id-1")
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

func TestGet(t *testing.T) {
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
			name:  "upload get ok",
			awsId: "some-id-1/1",
			prepFunc: func(t *testing.T, data *blob.BlobMetadata) *Client {
				s3Client, err := NewClient(testBucketName, cfg())
				require.NoError(t, err)

				testData := "testData"
				err = s3Client.Upload(ctx, data, bytes.NewReader([]byte(testData)))
				require.NoError(t, err)

				return s3Client
			},
			exptErr: nil,
			exptBlob: &blob.BlobMetadata{
				ContentType: "text/plain",
				Name:        "blob-name",
				Type:        "my-type",
				Id:          "some-id-1",
				Version:     "1",
				CreatedAt:   time.Now(),
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
					client.delete(ctx, tc.awsId)
				}()

			}
		})

	}

}
