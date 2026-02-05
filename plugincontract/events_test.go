package plugincontract

import (
	"encoding/json"
	"testing"
)

func TestEventPayload_JSONMarshal(t *testing.T) {
	event := EventPayload{
		EventType:  "account.created",
		OccurredAt: "2025-01-20T10:30:00Z",
		AccountID:  "abc123-def456",
		Data: Args{
			"quotaBytes": float64(10000000),
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["eventType"] != "account.created" {
		t.Errorf("eventType: expected 'account.created', got %v", result["eventType"])
	}
	if result["occurredAt"] != "2025-01-20T10:30:00Z" {
		t.Errorf("occurredAt: expected '2025-01-20T10:30:00Z', got %v", result["occurredAt"])
	}
	if result["accountId"] != "abc123-def456" {
		t.Errorf("accountId: expected 'abc123-def456', got %v", result["accountId"])
	}

	data2, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("expected data to be a map")
	}
	if data2["quotaBytes"].(float64) != 10000000 {
		t.Errorf("data.quotaBytes: expected 10000000, got %v", data2["quotaBytes"])
	}
}

func TestEventPayload_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"eventType": "account.created",
		"occurredAt": "2025-01-20T10:30:00Z",
		"accountId": "abc123-def456",
		"data": {
			"quotaBytes": 10000000
		}
	}`

	var event EventPayload
	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if event.EventType != "account.created" {
		t.Errorf("EventType: expected 'account.created', got %q", event.EventType)
	}
	if event.OccurredAt != "2025-01-20T10:30:00Z" {
		t.Errorf("OccurredAt: expected '2025-01-20T10:30:00Z', got %q", event.OccurredAt)
	}
	if event.AccountID != "abc123-def456" {
		t.Errorf("AccountID: expected 'abc123-def456', got %q", event.AccountID)
	}

	quotaBytes, ok := event.Data.Int("quotaBytes")
	if !ok {
		t.Fatal("expected Data.Int('quotaBytes') to return true")
	}
	if quotaBytes != 10000000 {
		t.Errorf("Data quotaBytes: expected 10000000, got %d", quotaBytes)
	}
}

func TestEventPayload_OmitEmptyData(t *testing.T) {
	event := EventPayload{
		EventType:  "account.deleted",
		OccurredAt: "2025-01-20T11:00:00Z",
		AccountID:  "xyz789",
		Data:       nil,
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if _, exists := result["data"]; exists {
		t.Error("expected data field to be omitted when nil")
	}
}

func TestEventPayload_EmptyDataIsOmitted(t *testing.T) {
	// Go's omitempty also omits empty maps, not just nil
	event := EventPayload{
		EventType:  "account.updated",
		OccurredAt: "2025-01-20T11:00:00Z",
		AccountID:  "xyz789",
		Data:       Args{},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if _, exists := result["data"]; exists {
		t.Error("expected data field to be omitted when empty map (Go omitempty behavior)")
	}
}
