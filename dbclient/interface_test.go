package dbclient_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jarrod-lowe/jmap-service-libs/dbclient"
)

// TestDynamoDBClientInterface verifies that *dynamodb.Client satisfies DynamoDBClient.
func TestDynamoDBClientInterface(t *testing.T) {
	// This test verifies at compile time that *dynamodb.Client implements DynamoDBClient.
	// We use a nil pointer since we only need to verify interface satisfaction.
	var client *dynamodb.Client
	var _ dbclient.DynamoDBClient = client
}
