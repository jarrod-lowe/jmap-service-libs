package dbclient_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jarrod-lowe/jmap-service-libs/dbclient"
)

func TestIsConditionalCheckFailed(t *testing.T) {
	t.Run("returns true for ConditionalCheckFailedException", func(t *testing.T) {
		err := &types.ConditionalCheckFailedException{}
		if !dbclient.IsConditionalCheckFailed(err) {
			t.Error("expected true for ConditionalCheckFailedException")
		}
	})

	t.Run("returns true for wrapped ConditionalCheckFailedException", func(t *testing.T) {
		inner := &types.ConditionalCheckFailedException{}
		err := errors.Join(errors.New("outer"), inner)
		if !dbclient.IsConditionalCheckFailed(err) {
			t.Error("expected true for wrapped ConditionalCheckFailedException")
		}
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("some error")
		if dbclient.IsConditionalCheckFailed(err) {
			t.Error("expected false for non-CCF error")
		}
	})

	t.Run("returns false for nil", func(t *testing.T) {
		if dbclient.IsConditionalCheckFailed(nil) {
			t.Error("expected false for nil")
		}
	})
}

func TestIsTransactionCanceled(t *testing.T) {
	t.Run("returns true for TransactionCanceledException", func(t *testing.T) {
		err := &types.TransactionCanceledException{}
		if !dbclient.IsTransactionCanceled(err) {
			t.Error("expected true for TransactionCanceledException")
		}
	})

	t.Run("returns true for wrapped TransactionCanceledException", func(t *testing.T) {
		inner := &types.TransactionCanceledException{}
		err := errors.Join(errors.New("outer"), inner)
		if !dbclient.IsTransactionCanceled(err) {
			t.Error("expected true for wrapped TransactionCanceledException")
		}
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("some error")
		if dbclient.IsTransactionCanceled(err) {
			t.Error("expected false for non-TC error")
		}
	})

	t.Run("returns false for nil", func(t *testing.T) {
		if dbclient.IsTransactionCanceled(nil) {
			t.Error("expected false for nil")
		}
	})
}

func TestGetTransactionCancellationReasons(t *testing.T) {
	t.Run("extracts reasons from TransactionCanceledException", func(t *testing.T) {
		err := &types.TransactionCanceledException{
			CancellationReasons: []types.CancellationReason{
				{Code: strPtr("None")},
				{Code: strPtr("ConditionalCheckFailed")},
				{Code: strPtr("None")},
			},
		}
		reasons := dbclient.GetTransactionCancellationReasons(err)
		if len(reasons) != 3 {
			t.Fatalf("got %d reasons, want 3", len(reasons))
		}
		if reasons[0].Index != 0 || reasons[0].Code != "None" {
			t.Errorf("reasons[0] = %+v, want {Index:0, Code:None}", reasons[0])
		}
		if reasons[1].Index != 1 || reasons[1].Code != "ConditionalCheckFailed" {
			t.Errorf("reasons[1] = %+v, want {Index:1, Code:ConditionalCheckFailed}", reasons[1])
		}
		if reasons[2].Index != 2 || reasons[2].Code != "None" {
			t.Errorf("reasons[2] = %+v, want {Index:2, Code:None}", reasons[2])
		}
	})

	t.Run("handles nil Code pointer", func(t *testing.T) {
		err := &types.TransactionCanceledException{
			CancellationReasons: []types.CancellationReason{
				{Code: nil},
			},
		}
		reasons := dbclient.GetTransactionCancellationReasons(err)
		if len(reasons) != 1 {
			t.Fatalf("got %d reasons, want 1", len(reasons))
		}
		if reasons[0].Code != "" {
			t.Errorf("reasons[0].Code = %q, want empty string", reasons[0].Code)
		}
	})

	t.Run("returns nil for non-TransactionCanceledException", func(t *testing.T) {
		err := errors.New("some error")
		reasons := dbclient.GetTransactionCancellationReasons(err)
		if reasons != nil {
			t.Errorf("expected nil, got %+v", reasons)
		}
	})

	t.Run("returns nil for nil error", func(t *testing.T) {
		reasons := dbclient.GetTransactionCancellationReasons(nil)
		if reasons != nil {
			t.Errorf("expected nil, got %+v", reasons)
		}
	})
}

func TestHasConditionalCheckFailure(t *testing.T) {
	t.Run("returns true when CCF in reasons", func(t *testing.T) {
		err := &types.TransactionCanceledException{
			CancellationReasons: []types.CancellationReason{
				{Code: strPtr("None")},
				{Code: strPtr("ConditionalCheckFailed")},
			},
		}
		if !dbclient.HasConditionalCheckFailure(err) {
			t.Error("expected true when CCF in reasons")
		}
	})

	t.Run("returns false when no CCF in reasons", func(t *testing.T) {
		err := &types.TransactionCanceledException{
			CancellationReasons: []types.CancellationReason{
				{Code: strPtr("None")},
				{Code: strPtr("ValidationError")},
			},
		}
		if dbclient.HasConditionalCheckFailure(err) {
			t.Error("expected false when no CCF in reasons")
		}
	})

	t.Run("returns false for non-TC error", func(t *testing.T) {
		err := errors.New("some error")
		if dbclient.HasConditionalCheckFailure(err) {
			t.Error("expected false for non-TC error")
		}
	})

	t.Run("returns false for nil", func(t *testing.T) {
		if dbclient.HasConditionalCheckFailure(nil) {
			t.Error("expected false for nil")
		}
	})
}

func TestGetConditionalCheckFailureIndex(t *testing.T) {
	t.Run("returns index of first CCF", func(t *testing.T) {
		err := &types.TransactionCanceledException{
			CancellationReasons: []types.CancellationReason{
				{Code: strPtr("None")},
				{Code: strPtr("ConditionalCheckFailed")},
				{Code: strPtr("ConditionalCheckFailed")},
			},
		}
		idx := dbclient.GetConditionalCheckFailureIndex(err)
		if idx != 1 {
			t.Errorf("got %d, want 1", idx)
		}
	})

	t.Run("returns -1 when no CCF", func(t *testing.T) {
		err := &types.TransactionCanceledException{
			CancellationReasons: []types.CancellationReason{
				{Code: strPtr("None")},
				{Code: strPtr("ValidationError")},
			},
		}
		idx := dbclient.GetConditionalCheckFailureIndex(err)
		if idx != -1 {
			t.Errorf("got %d, want -1", idx)
		}
	})

	t.Run("returns -1 for non-TC error", func(t *testing.T) {
		err := errors.New("some error")
		idx := dbclient.GetConditionalCheckFailureIndex(err)
		if idx != -1 {
			t.Errorf("got %d, want -1", idx)
		}
	})

	t.Run("returns -1 for nil", func(t *testing.T) {
		idx := dbclient.GetConditionalCheckFailureIndex(nil)
		if idx != -1 {
			t.Errorf("got %d, want -1", idx)
		}
	})
}

func strPtr(s string) *string {
	return &s
}
