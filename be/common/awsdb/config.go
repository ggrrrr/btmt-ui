package awsdb

type (
	AwsConfig struct {
		Region   string
		Endpoint string
		Database DynamodbConfig
	}

	DynamodbConfig struct {
		Database string
		Prefix   string
	}
)
