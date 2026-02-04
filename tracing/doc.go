// Package tracing provides OpenTelemetry instrumentation for JMAP services
// running on AWS Lambda with X-Ray.
//
// This package simplifies distributed tracing by providing:
//   - Automatic X-Ray and W3C Trace Context propagation setup
//   - Standardized attribute helpers for JMAP-specific context
//   - Span creation helpers for common Lambda patterns
//
// # Initialization
//
// Call Init during Lambda cold start to set up the tracer provider:
//
//	func init() {
//	    ctx := context.Background()
//	    tp, err := tracing.Init(ctx)
//	    if err != nil {
//	        log.Fatalf("failed to initialize tracing: %v", err)
//	    }
//	    // Use tp with otellambda.InstrumentHandler
//	}
//
// # Cold Start Spans
//
// Track initialization time with a cold start span:
//
//	func init() {
//	    ctx := context.Background()
//	    ctx, span := tracing.StartColdStartSpan(ctx, "my-function")
//	    defer span.End()
//
//	    // Perform initialization...
//	}
//
// # Handler Spans
//
// Create spans for Lambda handlers with attributes:
//
//	func handleRequest(ctx context.Context, event Event) error {
//	    ctx, span := tracing.StartHandlerSpan(ctx, "HandleEvent",
//	        tracing.RequestID(event.RequestID),
//	        tracing.AccountID(event.AccountID),
//	    )
//	    defer span.End()
//
//	    // Handle the request...
//	    return nil
//	}
//
// # JMAP Method Spans
//
// Track individual JMAP method calls within a request:
//
//	func processMethod(ctx context.Context, method string, clientID string, index int) {
//	    ctx, span := tracing.StartMethodSpan(ctx, "jmap-api", method, clientID, index)
//	    defer span.End()
//
//	    // Process the JMAP method...
//	}
//
// # Error Recording
//
// Record errors on spans with proper status:
//
//	if err := doWork(); err != nil {
//	    tracing.RecordError(span, err)
//	    return err
//	}
//
// # Using the Tracer Directly
//
// For custom spans, use the Tracer convenience function:
//
//	tracer := tracing.Tracer("my-component")
//	ctx, span := tracer.Start(ctx, "custom-operation")
//	defer span.End()
//
// For more information about OpenTelemetry concepts, see:
// https://opentelemetry.io/docs/concepts/
package tracing
