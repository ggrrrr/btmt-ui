package awsdb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

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

func NewClient(cfg AwsConfig) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:   aws.String(cfg.Region),
		Endpoint: aws.String(cfg.Endpoint),
	})

}
