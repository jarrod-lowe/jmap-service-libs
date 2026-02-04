// Package dbclient provides shared DynamoDB client utilities for JMAP services.
//
// This package consolidates common DynamoDB patterns used across jmap-service-core
// and jmap-service-email, including:
//
//   - A unified [DynamoDBClient] interface for DynamoDB operations
//   - Client creation via [NewClient] that integrates with awsinit
//   - Common key constants and helpers ([AttrPK], [AttrSK], [AccountPK], [UserPK])
//   - Error handling helpers for [ConditionalCheckFailedException] and
//     [TransactionCanceledException]
//
// # Usage with awsinit
//
// The typical usage pattern integrates with the awsinit package:
//
//	result, err := awsinit.Init(ctx, awsinit.WithHTTPHandler("my-handler"))
//	if err != nil {
//	    return err
//	}
//	ddb := dbclient.NewClient(result.Config)
//	repo := myrepo.New(ddb, os.Getenv("TABLE_NAME"))
//
// # Key Conventions
//
// JMAP services use a single-table design with common prefixes:
//
//	pk: "ACCOUNT#<accountId>" or "USER#<userId>"
//	sk: "META#" or domain-specific (e.g., "MAILBOX#<mailboxId>")
//
// # Error Handling
//
// The package provides helpers for common DynamoDB error scenarios:
//
//	if dbclient.IsConditionalCheckFailed(err) {
//	    return ErrNotFound
//	}
//
//	if idx := dbclient.GetConditionalCheckFailureIndex(err); idx >= 0 {
//	    // Handle specific item failure in transaction
//	}
package dbclient
