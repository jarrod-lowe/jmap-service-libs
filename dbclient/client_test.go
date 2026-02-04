package dbclient_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jarrod-lowe/jmap-service-libs/dbclient"
)

// TestNewClient verifies that NewClient returns a value satisfying DynamoDBClient.
func TestNewClient(t *testing.T) {
	cfg := aws.Config{}
	client := dbclient.NewClient(cfg)

	// Verify the returned value implements DynamoDBClient.
	var _ dbclient.DynamoDBClient = client
	if client == nil {
		t.Error("NewClient returned nil")
	}
}
