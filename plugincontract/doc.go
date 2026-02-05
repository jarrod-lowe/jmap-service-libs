// Package plugincontract defines the contract types for JMAP plugins.
//
// These types are used for communication between the JMAP core service and plugins.
// Plugins receive PluginInvocationRequest payloads and return PluginInvocationResponse
// payloads. System events are delivered to plugins as EventPayload messages via SQS.
//
// # Request/Response Types
//
// The core service invokes plugins by sending a PluginInvocationRequest:
//
//	type PluginInvocationRequest struct {
//	    RequestID string // API Gateway request ID for correlation
//	    CallIndex int    // Position in methodCalls array
//	    AccountID string // Authenticated account ID
//	    Method    string // JMAP method name (e.g., "Email/get")
//	    Args      Args   // Method arguments
//	    ClientID  string // Client-provided call identifier
//	}
//
// Plugins respond with a PluginInvocationResponse containing a MethodResponse:
//
//	type PluginInvocationResponse struct {
//	    MethodResponse MethodResponse
//	}
//
//	type MethodResponse struct {
//	    Name     string // Method name or "error"
//	    Args     Args   // Response data or error details
//	    ClientID string // Echo back the client ID
//	}
//
// # Args Helper Methods
//
// The Args type provides type-safe accessor methods for extracting values from
// method arguments or response data. These handle JSON's type system where all
// numbers unmarshal as float64:
//
//	func handler(ctx context.Context, req plugincontract.PluginInvocationRequest) (plugincontract.PluginInvocationResponse, error) {
//	    accountID, _ := req.Args.String("accountId")
//	    ids, _ := req.Args.StringSlice("ids")
//	    limit := req.Args.IntOr("limit", 100)
//
//	    // ... process request ...
//
//	    return plugincontract.PluginInvocationResponse{
//	        MethodResponse: plugincontract.MethodResponse{
//	            Name:     req.Method,
//	            Args:     plugincontract.Args{"accountId": accountID, "list": results},
//	            ClientID: req.ClientID,
//	        },
//	    }, nil
//	}
//
// # Event Payloads
//
// System events (such as account.created) are delivered to plugins via SQS:
//
//	type EventPayload struct {
//	    EventType  string // Event type identifier (e.g., "account.created")
//	    OccurredAt string // ISO 8601 timestamp
//	    AccountID  string // Related account ID
//	    Data       Args   // Event-specific data (optional)
//	}
package plugincontract
