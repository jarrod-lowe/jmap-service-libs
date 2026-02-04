package awsinit

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestWithHTTPHandler_SetsHTTPHandlerMode(t *testing.T) {
	cfg := &config{}

	opt := WithHTTPHandler("test-api")
	opt(cfg)

	if !cfg.httpHandler {
		t.Error("expected httpHandler to be true")
	}
	if cfg.functionName != "test-api" {
		t.Errorf("expected functionName 'test-api', got %q", cfg.functionName)
	}
}

func TestWithFunctionName_SetsFunctionName(t *testing.T) {
	cfg := &config{}

	opt := WithFunctionName("my-function")
	opt(cfg)

	if cfg.functionName != "my-function" {
		t.Errorf("expected functionName 'my-function', got %q", cfg.functionName)
	}
}

func TestInit_ReturnsResultWithTracerProvider(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{Region: "us-east-1"}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result == nil {
		t.Fatal("Init() returned nil result")
	}
	if result.TracerProvider != mockTP {
		t.Error("expected TracerProvider to be the mock provider")
	}
}

func TestInit_ReturnsResultWithAWSConfig(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{Region: "eu-west-1"}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result == nil {
		t.Fatal("Init() returned nil result")
	}
	if result.Config.Region != "eu-west-1" {
		t.Errorf("expected region 'eu-west-1', got %q", result.Config.Region)
	}
}

func TestInit_ReturnsContext(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result == nil {
		t.Fatal("Init() returned nil result")
	}
	if result.Ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestInit_WithHTTPHandler_CreatesColdStartSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		WithHTTPHandler("test-api"),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result == nil {
		t.Fatal("Init() returned nil result")
	}

	// Verify span was created by checking context contains a span
	span := trace.SpanFromContext(result.Ctx)
	if !span.SpanContext().IsValid() {
		t.Error("expected valid span in context for HTTP handler")
	}

	// Verify coldStartSpan is stored for cleanup
	if result.coldStartSpan == nil {
		t.Error("expected coldStartSpan to be set")
	}
}

func TestInit_WithoutHTTPHandler_NoColdStartSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if result == nil {
		t.Fatal("Init() returned nil result")
	}

	if result.coldStartSpan != nil {
		t.Error("expected no coldStartSpan for event-driven handler")
	}
}

func TestInit_TracingError_ReturnsTracingInitError(t *testing.T) {
	tracingErr := errors.New("xray connection failed")

	_, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return nil, tracingErr
		}),
	)

	if err == nil {
		t.Fatal("expected error")
	}

	var target *TracingInitError
	if !errors.As(err, &target) {
		t.Errorf("expected TracingInitError, got %T", err)
	}
	if !errors.Is(err, tracingErr) {
		t.Error("expected underlying error to be wrapped")
	}
}

func TestInit_AWSConfigError_ReturnsAWSConfigError(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	awsErr := errors.New("credentials not found")

	_, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, awsErr
		}),
	)

	if err == nil {
		t.Fatal("expected error")
	}

	var target *AWSConfigError
	if !errors.As(err, &target) {
		t.Errorf("expected AWSConfigError, got %T", err)
	}
	if !errors.Is(err, awsErr) {
		t.Error("expected underlying error to be wrapped")
	}
}

func TestResult_Cleanup_EndsColdStartSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		WithHTTPHandler("test-api"),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Call Cleanup to end the span
	result.Cleanup()

	// Force flush to export spans
	_ = mockTP.ForceFlush(context.Background())

	// Verify span was ended
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Name != "ColdStart" {
		t.Errorf("expected span name 'ColdStart', got %q", spans[0].Name)
	}
}

func TestResult_Cleanup_SafeWhenNoSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	result, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, nil
		}),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should not panic even when there's no cold start span
	result.Cleanup()
}

func TestResult_Start_WrapsHandlerAndStoresStartFunc(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	mockTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	var startCalled bool
	mockStarter := func(handler any, opts ...any) {
		startCalled = true
	}

	result, err := Init(context.Background(),
		withTracingInit(func(_ context.Context) (*sdktrace.TracerProvider, error) {
			return mockTP, nil
		}),
		withAWSConfigLoader(func(_ context.Context) (aws.Config, error) {
			return aws.Config{}, nil
		}),
		withLambdaStarter(mockStarter),
	)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	handler := func() {}
	result.Start(handler)

	if !startCalled {
		t.Error("expected lambda.Start to be called")
	}
}
