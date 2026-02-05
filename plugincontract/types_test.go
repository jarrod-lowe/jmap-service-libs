package plugincontract

import (
	"encoding/json"
	"testing"
)

func TestPluginInvocationRequest_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"requestId": "apigw-request-id",
		"callIndex": 0,
		"accountId": "user-123",
		"method": "Email/get",
		"args": {
			"accountId": "user-123",
			"ids": ["email-1", "email-2"],
			"properties": ["id", "subject", "from"]
		},
		"clientId": "c0"
	}`

	var req PluginInvocationRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.RequestID != "apigw-request-id" {
		t.Errorf("RequestID: expected 'apigw-request-id', got %q", req.RequestID)
	}
	if req.CallIndex != 0 {
		t.Errorf("CallIndex: expected 0, got %d", req.CallIndex)
	}
	if req.AccountID != "user-123" {
		t.Errorf("AccountID: expected 'user-123', got %q", req.AccountID)
	}
	if req.Method != "Email/get" {
		t.Errorf("Method: expected 'Email/get', got %q", req.Method)
	}
	if req.ClientID != "c0" {
		t.Errorf("ClientID: expected 'c0', got %q", req.ClientID)
	}

	// Test Args helper methods work on unmarshaled data
	accountID, ok := req.Args.String("accountId")
	if !ok {
		t.Fatal("expected Args.String('accountId') to return true")
	}
	if accountID != "user-123" {
		t.Errorf("Args accountId: expected 'user-123', got %q", accountID)
	}

	ids, ok := req.Args.StringSlice("ids")
	if !ok {
		t.Fatal("expected Args.StringSlice('ids') to return true")
	}
	if len(ids) != 2 || ids[0] != "email-1" || ids[1] != "email-2" {
		t.Errorf("Args ids: expected ['email-1', 'email-2'], got %v", ids)
	}
}

func TestPluginInvocationRequest_JSONMarshal(t *testing.T) {
	req := PluginInvocationRequest{
		RequestID: "req-123",
		CallIndex: 1,
		AccountID: "acc-456",
		Method:    "Mailbox/get",
		Args: Args{
			"accountId": "acc-456",
			"ids":       []any{"mb-1"},
		},
		ClientID: "c1",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back to verify round-trip
	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["requestId"] != "req-123" {
		t.Errorf("requestId: expected 'req-123', got %v", result["requestId"])
	}
	if result["callIndex"].(float64) != 1 {
		t.Errorf("callIndex: expected 1, got %v", result["callIndex"])
	}
	if result["accountId"] != "acc-456" {
		t.Errorf("accountId: expected 'acc-456', got %v", result["accountId"])
	}
	if result["method"] != "Mailbox/get" {
		t.Errorf("method: expected 'Mailbox/get', got %v", result["method"])
	}
	if result["clientId"] != "c1" {
		t.Errorf("clientId: expected 'c1', got %v", result["clientId"])
	}
}

func TestMethodResponse_JSONMarshal(t *testing.T) {
	resp := MethodResponse{
		Name: "Email/get",
		Args: Args{
			"accountId": "user-123",
			"list":      []any{},
			"notFound":  []any{},
		},
		ClientID: "c0",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["name"] != "Email/get" {
		t.Errorf("name: expected 'Email/get', got %v", result["name"])
	}
	if result["clientId"] != "c0" {
		t.Errorf("clientId: expected 'c0', got %v", result["clientId"])
	}

	args, ok := result["args"].(map[string]any)
	if !ok {
		t.Fatal("expected args to be a map")
	}
	if args["accountId"] != "user-123" {
		t.Errorf("args.accountId: expected 'user-123', got %v", args["accountId"])
	}
}

func TestMethodResponse_ErrorResponse(t *testing.T) {
	// Test error response format
	resp := MethodResponse{
		Name: "error",
		Args: Args{
			"type":        "invalidArguments",
			"description": "ids is required",
		},
		ClientID: "c0",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["name"] != "error" {
		t.Errorf("name: expected 'error', got %v", result["name"])
	}

	args, ok := result["args"].(map[string]any)
	if !ok {
		t.Fatal("expected args to be a map")
	}
	if args["type"] != "invalidArguments" {
		t.Errorf("args.type: expected 'invalidArguments', got %v", args["type"])
	}
	if args["description"] != "ids is required" {
		t.Errorf("args.description: expected 'ids is required', got %v", args["description"])
	}
}

func TestPluginInvocationResponse_JSONMarshal(t *testing.T) {
	resp := PluginInvocationResponse{
		MethodResponse: MethodResponse{
			Name: "Email/get",
			Args: Args{
				"accountId": "user-123",
				"list":      []any{},
				"notFound":  []any{},
			},
			ClientID: "c0",
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	methodResp, ok := result["methodResponse"].(map[string]any)
	if !ok {
		t.Fatal("expected methodResponse to be a map")
	}
	if methodResp["name"] != "Email/get" {
		t.Errorf("methodResponse.name: expected 'Email/get', got %v", methodResp["name"])
	}
	if methodResp["clientId"] != "c0" {
		t.Errorf("methodResponse.clientId: expected 'c0', got %v", methodResp["clientId"])
	}
}

func TestPluginInvocationResponse_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"methodResponse": {
			"name": "Email/get",
			"args": {
				"accountId": "user-123",
				"list": [{"id": "email-1"}],
				"notFound": []
			},
			"clientId": "c0"
		}
	}`

	var resp PluginInvocationResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.MethodResponse.Name != "Email/get" {
		t.Errorf("Name: expected 'Email/get', got %q", resp.MethodResponse.Name)
	}
	if resp.MethodResponse.ClientID != "c0" {
		t.Errorf("ClientID: expected 'c0', got %q", resp.MethodResponse.ClientID)
	}

	accountID, ok := resp.MethodResponse.Args.String("accountId")
	if !ok {
		t.Fatal("expected Args.String('accountId') to return true")
	}
	if accountID != "user-123" {
		t.Errorf("Args accountId: expected 'user-123', got %q", accountID)
	}
}
