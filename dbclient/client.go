package dbclient

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewClient creates a DynamoDB client from an AWS config.
// The config should already have OTel middleware configured (e.g., from awsinit.Init).
func NewClient(cfg aws.Config) DynamoDBClient {
	return dynamodb.NewFromConfig(cfg)
}
