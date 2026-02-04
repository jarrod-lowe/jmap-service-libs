package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestRequestID(t *testing.T) {
	attr := RequestID("test-request-123")

	if attr.Key != "request_id" {
		t.Errorf("expected key 'request_id', got %q", attr.Key)
	}
	if attr.Value.AsString() != "test-request-123" {
		t.Errorf("expected value 'test-request-123', got %q", attr.Value.AsString())
	}
}

func TestAccountID(t *testing.T) {
	attr := AccountID("user-456")

	if attr.Key != "account_id" {
		t.Errorf("expected key 'account_id', got %q", attr.Key)
	}
	if attr.Value.AsString() != "user-456" {
		t.Errorf("expected value 'user-456', got %q", attr.Value.AsString())
	}
}

func TestBlobID(t *testing.T) {
	attr := BlobID("blob-789")

	if attr.Key != "blob_id" {
		t.Errorf("expected key 'blob_id', got %q", attr.Key)
	}
	if attr.Value.AsString() != "blob-789" {
		t.Errorf("expected value 'blob-789', got %q", attr.Value.AsString())
	}
}

func TestParentBlobID(t *testing.T) {
	attr := ParentBlobID("parent-blob-123")

	if attr.Key != "parent_blob_id" {
		t.Errorf("expected key 'parent_blob_id', got %q", attr.Key)
	}
	if attr.Value.AsString() != "parent-blob-123" {
		t.Errorf("expected value 'parent-blob-123', got %q", attr.Value.AsString())
	}
}

func TestContentType(t *testing.T) {
	attr := ContentType("application/json")

	if attr.Key != "content_type" {
		t.Errorf("expected key 'content_type', got %q", attr.Key)
	}
	if attr.Value.AsString() != "application/json" {
		t.Errorf("expected value 'application/json', got %q", attr.Value.AsString())
	}
}

func TestFunction(t *testing.T) {
	attr := Function("blob-upload")

	if attr.Key != "function" {
		t.Errorf("expected key 'function', got %q", attr.Key)
	}
	if attr.Value.AsString() != "blob-upload" {
		t.Errorf("expected value 'blob-upload', got %q", attr.Value.AsString())
	}
}

func TestJMAPMethod(t *testing.T) {
	attr := JMAPMethod("Email/get")

	if attr.Key != "jmap.method" {
		t.Errorf("expected key 'jmap.method', got %q", attr.Key)
	}
	if attr.Value.AsString() != "Email/get" {
		t.Errorf("expected value 'Email/get', got %q", attr.Value.AsString())
	}
}

func TestJMAPClientID(t *testing.T) {
	attr := JMAPClientID("c0")

	if attr.Key != "jmap.client_id" {
		t.Errorf("expected key 'jmap.client_id', got %q", attr.Key)
	}
	if attr.Value.AsString() != "c0" {
		t.Errorf("expected value 'c0', got %q", attr.Value.AsString())
	}
}

func TestJMAPCallIndex(t *testing.T) {
	attr := JMAPCallIndex(2)

	if attr.Key != "jmap.call_index" {
		t.Errorf("expected key 'jmap.call_index', got %q", attr.Key)
	}
	if attr.Value.AsInt64() != 2 {
		t.Errorf("expected value 2, got %d", attr.Value.AsInt64())
	}
}

func TestTracer(t *testing.T) {
	tracer := Tracer("test-tracer")

	if tracer == nil {
		t.Fatal("expected non-nil tracer")
	}

	// Verify we can create a span with the tracer
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	// Get a new tracer after setting the provider
	tracer = Tracer("test-tracer")
	ctx, span := tracer.Start(context.Background(), "test-span")
	span.End()

	_ = tp.ForceFlush(ctx)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "test-span" {
		t.Errorf("expected span name 'test-span', got %q", spans[0].Name)
	}
}

func TestStartHandlerSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	ctx := context.Background()

	ctx, span := StartHandlerSpan(ctx, "TestHandler",
		RequestID("req-123"),
		AccountID("acct-456"),
	)
	span.End()

	_ = tp.ForceFlush(ctx)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	s := spans[0]
	if s.Name != "TestHandler" {
		t.Errorf("expected span name 'TestHandler', got %q", s.Name)
	}

	attrMap := make(map[string]string)
	for _, attr := range s.Attributes {
		attrMap[string(attr.Key)] = attr.Value.AsString()
	}

	if attrMap["request_id"] != "req-123" {
		t.Errorf("expected request_id 'req-123', got %q", attrMap["request_id"])
	}
	if attrMap["account_id"] != "acct-456" {
		t.Errorf("expected account_id 'acct-456', got %q", attrMap["account_id"])
	}
}

func TestStartMethodSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	ctx := context.Background()

	ctx, span := StartMethodSpan(ctx, "jmap-api", "Email/get", "c0", 1)
	span.End()

	_ = tp.ForceFlush(ctx)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	s := spans[0]
	if s.Name != "JMAP Method" {
		t.Errorf("expected span name 'JMAP Method', got %q", s.Name)
	}

	attrMap := make(map[attribute.Key]attribute.Value)
	for _, attr := range s.Attributes {
		attrMap[attr.Key] = attr.Value
	}

	if attrMap["jmap.method"].AsString() != "Email/get" {
		t.Errorf("expected jmap.method 'Email/get', got %q", attrMap["jmap.method"].AsString())
	}
	if attrMap["jmap.client_id"].AsString() != "c0" {
		t.Errorf("expected jmap.client_id 'c0', got %q", attrMap["jmap.client_id"].AsString())
	}
	if attrMap["jmap.call_index"].AsInt64() != 1 {
		t.Errorf("expected jmap.call_index 1, got %d", attrMap["jmap.call_index"].AsInt64())
	}
}

func TestStartColdStartSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	ctx := context.Background()

	ctx, span := StartColdStartSpan(ctx, "test-function")
	span.End()

	_ = tp.ForceFlush(ctx)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	s := spans[0]
	if s.Name != "ColdStart" {
		t.Errorf("expected span name 'ColdStart', got %q", s.Name)
	}

	attrMap := make(map[string]string)
	for _, attr := range s.Attributes {
		attrMap[string(attr.Key)] = attr.Value.AsString()
	}
	if attrMap["function"] != "test-function" {
		t.Errorf("expected function 'test-function', got %q", attrMap["function"])
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

func TestRecordError(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	ctx := context.Background()
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "TestSpan")

	testErr := &testError{message: "something went wrong"}
	RecordError(span, testErr)
	span.End()

	_ = tp.ForceFlush(ctx)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	s := spans[0]

	if len(s.Events) == 0 {
		t.Error("expected at least one event (error), got none")
	}

	if s.Status.Code != codes.Error {
		t.Errorf("expected error status code %d, got %d", codes.Error, s.Status.Code)
	}
}

func TestInitSetsPropagator(t *testing.T) {
	// Reset propagator to default before test
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator())

	InitPropagator()

	propagator := otel.GetTextMapPropagator()

	carrier := propagation.MapCarrier{}

	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	propagator.Inject(ctx, carrier)

	// Verify X-Amzn-Trace-Id header is set (X-Ray propagator injects this)
	xrayHeader := carrier.Get("X-Amzn-Trace-Id")
	if xrayHeader == "" {
		t.Error("expected X-Amzn-Trace-Id header to be set after Init, got empty string")
	}

	// Verify traceparent header is also set (W3C TraceContext propagator)
	traceparent := carrier.Get("traceparent")
	if traceparent == "" {
		t.Error("expected traceparent header to be set after Init, got empty string")
	}
}
