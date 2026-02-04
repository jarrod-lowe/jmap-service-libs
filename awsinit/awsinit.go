package awsinit

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/jarrod-lowe/jmap-service-libs/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Option configures Init behavior.
type Option func(*config)

type config struct {
	httpHandler     bool
	functionName    string
	tracingInit     func(context.Context) (*sdktrace.TracerProvider, error)
	awsConfigLoader func(context.Context) (aws.Config, error)
	lambdaStarter   func(handler any, opts ...any)
}

// WithHTTPHandler configures Init for HTTP handlers (API Gateway).
// Creates a cold start span that should be ended with Result.Cleanup().
func WithHTTPHandler(functionName string) Option {
	return func(c *config) {
		c.httpHandler = true
		c.functionName = functionName
	}
}

// WithFunctionName sets the function name for tracing.
// Use this for event-driven handlers that want custom function names.
func WithFunctionName(name string) Option {
	return func(c *config) {
		c.functionName = name
	}
}

// Internal testing options
func withTracingInit(fn func(context.Context) (*sdktrace.TracerProvider, error)) Option {
	return func(c *config) {
		c.tracingInit = fn
	}
}

func withAWSConfigLoader(fn func(context.Context) (aws.Config, error)) Option {
	return func(c *config) {
		c.awsConfigLoader = fn
	}
}

func withLambdaStarter(fn func(handler any, opts ...any)) Option {
	return func(c *config) {
		c.lambdaStarter = fn
	}
}

// Result contains the initialized resources from Init.
type Result struct {
	Ctx            context.Context
	Config         aws.Config
	TracerProvider *sdktrace.TracerProvider
	coldStartSpan  trace.Span
	lambdaStarter  func(handler any, opts ...any)
}

// Cleanup ends the cold start span if one was created.
// Safe to call even when no span exists.
func (r *Result) Cleanup() {
	if r.coldStartSpan != nil {
		r.coldStartSpan.End()
	}
}

// Start wraps the handler with OTel instrumentation and starts the Lambda.
func (r *Result) Start(handler any) {
	wrapped := otellambda.InstrumentHandler(handler,
		otellambda.WithTracerProvider(r.TracerProvider),
		otellambda.WithFlusher(r.TracerProvider),
	)
	r.lambdaStarter(wrapped)
}

func defaultAWSConfigLoader(ctx context.Context) (aws.Config, error) {
	return awsconfig.LoadDefaultConfig(ctx)
}

func defaultLambdaStarter(handler any, _ ...any) {
	lambda.Start(handler)
}

// Init initializes tracing and AWS config for Lambda handlers.
func Init(ctx context.Context, opts ...Option) (*Result, error) {
	cfg := &config{
		tracingInit:     tracing.Init,
		awsConfigLoader: defaultAWSConfigLoader,
		lambdaStarter:   defaultLambdaStarter,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Initialize tracing
	tp, err := cfg.tracingInit(ctx)
	if err != nil {
		return nil, &TracingInitError{Err: err}
	}
	otel.SetTracerProvider(tp)

	// Create cold start span if HTTP handler
	var coldStartSpan trace.Span
	if cfg.httpHandler && cfg.functionName != "" {
		ctx, coldStartSpan = tracing.StartColdStartSpan(ctx, cfg.functionName)
	}

	// Load AWS config
	awsCfg, err := cfg.awsConfigLoader(ctx)
	if err != nil {
		if coldStartSpan != nil {
			coldStartSpan.End()
		}
		return nil, &AWSConfigError{Err: err}
	}
	otelaws.AppendMiddlewares(&awsCfg.APIOptions)

	return &Result{
		Ctx:            ctx,
		Config:         awsCfg,
		TracerProvider: tp,
		coldStartSpan:  coldStartSpan,
		lambdaStarter:  cfg.lambdaStarter,
	}, nil
}
