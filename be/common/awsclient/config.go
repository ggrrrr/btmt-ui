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
)
