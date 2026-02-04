# Tracing Module Design

## Overview

Shared OpenTelemetry/X-Ray tracing module for the jmap-service-* family of repositories. Provides consistent tracer initialization, span creation helpers, and attribute helpers.

## Development Process

**All implementation MUST use the TDD Superpower.**

RED test requirements:
- Tests MUST compile
- Tests MUST run
- Tests MUST NOT panic
- Tests MUST fail (proving the test is actually testing something)

Only after RED tests are verified do we write implementation to make them pass.

## Package Import

```go
import "github.com/jarrod-lowe/jmap-service-libs/tracing"
```

## API

### Initialization

```go
// Init initializes the X-Ray tracer provider and sets up propagators.
// Call once at Lambda cold start before any AWS SDK clients are created.
// Returns the TracerProvider for use with otellambda.InstrumentHandler.
func Init(ctx context.Context) (*sdktrace.TracerProvider, error)

// InitPropagator registers the composite text map propagator.
// Called automatically by Init. Exported for testing.
func InitPropagator()

// Tracer returns a tracer with the given name.
// Convenience wrapper around otel.Tracer().
func Tracer(name string) trace.Tracer
```

### Attribute Helpers

Each returns `attribute.KeyValue` for consistent span annotation:

| Function | Attribute Key | Type |
|----------|--------------|------|
| `RequestID(id string)` | `request_id` | string |
| `AccountID(id string)` | `account_id` | string |
| `BlobID(id string)` | `blob_id` | string |
| `ParentBlobID(id string)` | `parent_blob_id` | string |
| `ContentType(ct string)` | `content_type` | string |
| `Function(name string)` | `function` | string |
| `JMAPMethod(method string)` | `jmap.method` | string |
| `JMAPClientID(id string)` | `jmap.client_id` | string |
| `JMAPCallIndex(idx int)` | `jmap.call_index` | int |

### Span Helpers

```go
// StartHandlerSpan creates a root span for a Lambda handler.
// Caller must defer span.End().
func StartHandlerSpan(ctx context.Context, handlerName string, attrs ...attribute.KeyValue) (context.Context, trace.Span)

// StartMethodSpan creates a span for a JMAP method call with standard attributes.
// Caller must defer span.End().
func StartMethodSpan(ctx context.Context, tracerName, methodName, clientID string, callIndex int) (context.Context, trace.Span)

// StartColdStartSpan creates a span for Lambda cold start initialization.
// Caller must defer span.End().
func StartColdStartSpan(ctx context.Context, functionName string) (context.Context, trace.Span)

// RecordError marks the span as errored and records error details.
func RecordError(span trace.Span, err error)
```

## Usage Examples

### Cold Start Initialization

```go
func main() {
    ctx := context.Background()

    // Initialize tracer before any AWS calls
    tp, err := tracing.Init(ctx)
    if err != nil {
        panic(err)
    }
    otel.SetTracerProvider(tp)

    // Track cold start
    ctx, coldSpan := tracing.StartColdStartSpan(ctx, "email-get")
    defer coldSpan.End()

    // Initialize AWS SDK (now traced)
    cfg, _ := config.LoadDefaultConfig(ctx)
    otelaws.AppendMiddlewares(&cfg.APIOptions)

    // Start handler
    lambda.Start(otellambda.InstrumentHandler(handler, xrayconfig.WithRecommendedOptions(tp)...))
}
```

### Handler Span

```go
func handler(ctx context.Context, req Request) (Response, error) {
    ctx, span := tracing.StartHandlerSpan(ctx, "EmailGetHandler",
        tracing.Function("email-get"),
        tracing.AccountID(req.AccountID),
    )
    defer span.End()

    // Handler logic...
    if err != nil {
        tracing.RecordError(span, err)
        return Response{}, err
    }

    return response, nil
}
```

### Using Tracer Directly

```go
tracer := tracing.Tracer("jmap-blob-client")
ctx, span := tracer.Start(ctx, "blob.Fetch",
    trace.WithAttributes(
        tracing.AccountID(accountID),
        tracing.BlobID(blobID),
    ))
defer span.End()
```

## Dependencies

```
go.opentelemetry.io/otel
go.opentelemetry.io/otel/attribute
go.opentelemetry.io/otel/codes
go.opentelemetry.io/otel/propagation
go.opentelemetry.io/otel/sdk/trace
go.opentelemetry.io/otel/trace
go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig
go.opentelemetry.io/contrib/propagators/aws/xray
```

## Testing Approach

Tests use `tracetest.NewInMemoryExporter()` for verification without real X-Ray.

Each function gets a test that:
1. Sets up in-memory exporter
2. Calls the function
3. Verifies the expected behavior (attribute key/value, span name, error status, etc.)
