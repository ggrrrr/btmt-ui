package awss3

import (
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
)

func cfg() awsclient.AwsConfig {
	return awsclient.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	}
}

func TestPut(t *testing.T) {
	// testBucketName := "test-bucket-1"
	// ctx := context.Background()

	tests := []struct {
		name     string
		awsId    string
		prepFunc func(t *testing.T, data *blob.BlobMetadata) *Client
		exptErr  error
		exptBlob *blob.BlobMetadata
	}{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// client := tc.prepFunc(t, tc.exptBlob)
			// resp, err := client.get(ctx, tc.awsId)
		})
	}
}
