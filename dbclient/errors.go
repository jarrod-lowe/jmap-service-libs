package dbclient

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// IsConditionalCheckFailed returns true if the error is a DynamoDB
// ConditionalCheckFailedException.
func IsConditionalCheckFailed(err error) bool {
	var ccf *types.ConditionalCheckFailedException
	return errors.As(err, &ccf)
}

// IsTransactionCanceled returns true if the error is a DynamoDB
// TransactionCanceledException.
func IsTransactionCanceled(err error) bool {
	var tc *types.TransactionCanceledException
	return errors.As(err, &tc)
}

// TransactionCancellationReason represents the reason for a transaction item failure.
type TransactionCancellationReason struct {
	Index int
	Code  string
}

// GetTransactionCancellationReasons extracts the cancellation reasons from a
// TransactionCanceledException. Returns nil if not a TransactionCanceledException.
func GetTransactionCancellationReasons(err error) []TransactionCancellationReason {
	var tc *types.TransactionCanceledException
	if !errors.As(err, &tc) {
		return nil
	}

	reasons := make([]TransactionCancellationReason, len(tc.CancellationReasons))
	for i, r := range tc.CancellationReasons {
		code := ""
		if r.Code != nil {
			code = *r.Code
		}
		reasons[i] = TransactionCancellationReason{
			Index: i,
			Code:  code,
		}
	}
	return reasons
}

// HasConditionalCheckFailure checks if any item in a canceled transaction
// failed due to ConditionalCheckFailed.
func HasConditionalCheckFailure(err error) bool {
	return GetConditionalCheckFailureIndex(err) >= 0
}

// GetConditionalCheckFailureIndex returns the index of the first item that
// failed with ConditionalCheckFailed, or -1 if none.
func GetConditionalCheckFailureIndex(err error) int {
	reasons := GetTransactionCancellationReasons(err)
	for _, r := range reasons {
		if r.Code == "ConditionalCheckFailed" {
			return r.Index
		}
	}
	return -1
}
