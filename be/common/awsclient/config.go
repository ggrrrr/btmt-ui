package awsclient

type (
	Config struct {
		Region      string `env:"AWS_REGION"`
		EndpointURL string `env:"AWS_ENDPOINT_URL"`
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
