package awsclient

type (
	AwsConfig struct {
		Region   string
		Endpoint string
	}

	DynamodbConfig struct {
		Database string
		Prefix   string
	}

	//
	S3Client struct {
		Region     string
		BucketName string
	}

	S3Config struct {
		Tenants map[string]S3Client
	}
)
