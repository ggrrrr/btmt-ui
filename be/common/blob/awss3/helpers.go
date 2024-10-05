package awss3

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
)

type awsId struct {
	path string
	id   string
	ver  string
}

func timeInLocal(from *time.Time) time.Time {
	if from == nil {
		return time.Now()
	}
	to := from.Local()
	return to
}

func awsIdFromString(fromStr string) (awsId, error) {
	blobId, err := blob.ParseBlobId(fromStr)
	if err != nil {
		return awsId{}, err
	}
	return awsId{
		path: blobId.Path(),
		id:   blobId.Name(),
		ver:  blobId.Version(),
	}, nil

}

// folder/id
func (i awsId) idPath() string {
	return fmt.Sprintf("%s/%s", i.path, i.id)
}

// folder/id/ver
func (i awsId) keyVer() string {
	if i.ver == "" {
		return fmt.Sprintf("%s/%s", i.path, i.id)
	}
	return fmt.Sprintf("%s/%s/%s", i.path, i.id, i.ver)
}

// folder/id/ver
func (i awsId) String() string {
	return i.keyVer()
}

func fromAwsHeadToBlobInfo(result *s3.HeadObjectOutput) blob.BlobInfo {
	out := blob.BlobInfo{
		ContentType:   *result.ContentType,
		ContentLength: *result.ContentLength,
		CreatedAt:     timeInLocal(result.LastModified),
		Type:          result.Metadata["type"],
		Name:          result.Metadata["name"],
		Owner:         result.Metadata["owner"],
	}
	width, widtErr := strconv.ParseInt(result.Metadata["image_width"], 10, 64)
	height, heightErr := strconv.ParseInt(result.Metadata["image_height"], 10, 64)
	if widtErr == nil && heightErr == nil {
		out.ImageInfo = blob.ImageInfo{
			Width:  width,
			Height: height,
		}
	}

	return out
}

func fromAwsGetToBlobInfo(result *s3.GetObjectOutput) blob.BlobInfo {
	out := blob.BlobInfo{
		ContentType:   *result.ContentType,
		ContentLength: *result.ContentLength,
		CreatedAt:     timeInLocal(result.LastModified),
		Type:          result.Metadata["type"],
		Name:          result.Metadata["name"],
		Owner:         result.Metadata["owner"],
	}
	width, widtErr := strconv.ParseInt(result.Metadata["image_width"], 10, 64)
	height, heightErr := strconv.ParseInt(result.Metadata["image_height"], 10, 64)
	if widtErr == nil && heightErr == nil {
		out.ImageInfo = blob.ImageInfo{
			Width:  width,
			Height: height,
		}
	}

	return out
}

func fromBlobInfoToAwsObject(c s3Client, id awsId, info blob.BlobInfo) *s3.PutObjectInput {
	md := map[string]string{
		"type":  info.Type,
		"name":  info.Name,
		"owner": info.Owner,
	}
	if info.ImageInfo.Height > 0 && info.ImageInfo.Width > 0 {
		md["image_width"] = fmt.Sprintf("%d", info.ImageInfo.Width)
		md["image_height"] = fmt.Sprintf("%d", info.ImageInfo.Height)
	}

	return &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(id.keyVer()),
		ContentType: aws.String(info.ContentType),
		Metadata:    md,
	}
}
