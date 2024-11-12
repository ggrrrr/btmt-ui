package awss3

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
)

var testBucket string = "test-bucket-1"

func TestList2(t *testing.T) {
	ctx := context.TODO()

	folderName := "root/folder-1"
	fileFolderId, _ := blob.ParseBlobId(folderName)

	file1v1, _ := blob.ParseBlobId(fmt.Sprintf("%s/file-1:1", folderName))
	file1v2, _ := blob.ParseBlobId(fmt.Sprintf("%s/file-1:2", folderName))
	file2v1, _ := blob.ParseBlobId(fmt.Sprintf("%s/file-2:1", folderName))
	file2v2, _ := blob.ParseBlobId(fmt.Sprintf("%s/file-2:2", folderName))
	file2v3, _ := blob.ParseBlobId(fmt.Sprintf("%s/file-2:3", folderName))

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "ok Push/List",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = testBucket
				testClient := &Client{
					s3Clients: map[realmKey]s3Client{
						"localhost": s3c,
					},
				}

				data := "my data asd"
				md := blob.BlobMD{
					Type:        "some_type",
					ContentType: "text/plain",
					Name:        "some-name",
					Owner:       "user-1",
					CreatedAt:   time.Now(),
					ImageInfo: blob.MDImageInfo{
						Width:  10,
						Height: 123,
					},
					ContentLength: int64(len(data)),
				}

				delayDuration := time.Duration(900 * time.Millisecond)

				_, err = testClient.Push(ctx, "localhost", file1v1, md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)
				time.Sleep(delayDuration)
				_, err = testClient.Push(ctx, "localhost", file1v2, md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)
				time.Sleep(delayDuration)
				_, err = testClient.Push(ctx, "localhost", file2v1, md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)
				time.Sleep(delayDuration)
				_, err = testClient.Push(ctx, "localhost", file2v2, md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)
				time.Sleep(delayDuration)
				_, err = testClient.Push(ctx, "localhost", file2v3, md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)

				defer func() {
					_ = deleteAll(ctx, s3c, awsId{
						path: folderName,
					})
				}()

				// file2v1 := fmt.Sprintf("%s", folderName)
				list2, err := testClient.List(ctx, "localhost", fileFolderId)
				require.NoError(t, err)

				fmt.Printf("\n\ntestClient\n")
				for _, ll := range list2 {
					fmt.Printf("\t%+v: %v \n", ll.Id, ll.Blob.MD.CreatedAt)
					for _, vv := range ll.Versions {
						fmt.Printf("\t\t : %+v %v\n", vv.Id, vv.MD.CreatedAt)
					}
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestCallsListPushFetchHead(t *testing.T) {
	ctx := context.TODO()
	var folder1v1 = blob.BlobId{}
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "ok Push/Fetch",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = testBucket
				testClient := &Client{
					s3Clients: map[realmKey]s3Client{
						"localhost": s3c,
					},
				}

				data := "my data"
				md := blob.BlobMD{
					Type:        "some_type",
					ContentType: "text/plain",
					Name:        "some-name",
					Owner:       "user-1",
					CreatedAt:   time.Now(),
					ImageInfo: blob.MDImageInfo{
						Width:  10,
						Height: 123,
					},
					ContentLength: int64(len(data)),
				}

				folder1v1, _ = blob.ParseBlobId("folder-1/id-1")

				newID, err := testClient.Push(ctx, "localhost", folder1v1, md, bytes.NewReader([]byte(data)))
				require.NoError(t, err)
				defer func() {
					_ = deleteAll(ctx, s3c, awsId{path: "folder-1", id: "id-1"})
				}()
				verID, err := blob.NewBlobId("folder-1", "id-1", "1")
				require.NoError(t, err)
				assert.Equal(t, verID, newID)

				folder1v1, _ = blob.ParseBlobId("folder-1/id-1")
				blobTest, err := testClient.Fetch(ctx, "localhost", folder1v1)
				require.NoError(t, err)
				blob.TestBlobInfoData(t, blobTest, md, data, 1000)

				folder1v1, _ = blob.ParseBlobId("folder-1/id-1:1")
				blobInfo, err := testClient.Head(ctx, "localhost", folder1v1)
				require.NoError(t, err)
				blob.TestBlobInfo("head", t, blobInfo, md, 1000)

			},
		},
		{
			name: "ok tenent not found",
			testFunc: func(t *testing.T) {
				testClient := &Client{
					s3Clients: map[realmKey]s3Client{},
				}

				folder1v1, _ = blob.ParseBlobId("asdasdasd/asd")

				_, err := testClient.Push(ctx, "notfound", folder1v1, blob.BlobMD{}, nil)
				require.Error(t, err)

				_, err = testClient.Fetch(ctx, "notfound", folder1v1)
				require.Error(t, err)

				_, err = testClient.Head(ctx, "notfound", folder1v1)
				require.Error(t, err)

				_, err = testClient.List(ctx, "notfound", folder1v1)
				require.Error(t, err)
			},
		},
		{
			name: "ok Fetch not found",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = testBucket

				testClient := &Client{
					s3Clients: map[realmKey]s3Client{
						"localhost": s3c,
					},
				}
				folder1v1, _ = blob.ParseBlobId("ad/asd")

				_, err = testClient.Fetch(ctx, "localhost", folder1v1)
				require.Error(t, err)
				// notFound := &blob.NotFoundError{}
				// assert.ErrorAsf(t, err, &notFound, "%#v, %+v", err, err)

			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
		// return
	}
}

func TestGetPutDelete(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "ok list 0",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				blockId := awsId{path: "not", id: "no", ver: "asd"}
				result, err := list(ctx, s3c, blockId)
				require.NoError(t, err)
				assert.Equal(t, 0, len(result))
			},
		},
		{
			name: "ok list bucket not found",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = "notfound"
				blockId := awsId{path: "not", id: "no", ver: "asd"}
				_, err = list(ctx, s3c, blockId)
				require.Error(t, err)
				// storErr := &blob.StoreNotFoundError{}
				// assert.ErrorAs(t, err, &storErr)
			},
		},
		{
			name: "ok get bucket not found",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = "notfound"
				blockId := awsId{path: "not", id: "no", ver: "asd"}
				_, err = get(ctx, s3c, blockId)
				require.Error(t, err)
			},
		},
		{
			name: "ok get put head id.ver empty",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = testBucket
				blockId := awsId{path: "not", id: "no", ver: ""}
				_, err = get(ctx, s3c, blockId)
				require.Error(t, err)

				_, err = head(ctx, s3c, blockId)
				require.Error(t, err)

				_, err = put(ctx, s3c, blockId, blob.BlobMD{}, nil)
				require.Error(t, err)
			},
		},
		{
			name: "ok push 1 list 1 get 1",
			testFunc: func(t *testing.T) {
				s3c, err := createS3Client(cfg())
				require.NoError(t, err)
				s3c.bucketName = testBucket

				data := "mydata"

				id := awsId{path: "folder-1", id: "id-1", ver: "1"}
				blobInfo := blob.BlobMD{
					Type:          "some_type",
					ContentType:   "text/plain",
					Name:          "some-name",
					Owner:         "user1",
					ContentLength: int64(len(data)),
					CreatedAt:     time.Now(),
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
				assert.Equal(t, id.path, blobResult.Blob.Id.Path())
				assert.Equal(t, id.id, blobResult.Blob.Id.Id())
				assert.Equal(t, id.ver, blobResult.Blob.Id.Version())

				headMd, err := head(ctx, s3c, id)
				require.NoError(t, err)
				blob.TestBlobInfo("head", t, headMd, blobInfo, 1000)

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

func cfg() (awsclient.AwsConfig, awsclient.S3Client) {
	return awsclient.AwsConfig{
			Region:   "us-east-1",
			Endpoint: "http://localhost:4566",
		},
		awsclient.S3Client{
			Region:     "us-east-1",
			BucketName: testBucket,
		}
}

func timePrint(name string, t time.Time) {
	fmt.Printf("%15s: [%10s] [%5s] time: %+v \n", name, t.Location(), t.Format("-0700"), t)
}
