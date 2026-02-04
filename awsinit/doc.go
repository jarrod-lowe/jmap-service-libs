// Package awsinit provides Lambda initialization boilerplate for AWS Lambda handlers
// with OpenTelemetry tracing integration.
//
// It encapsulates the common pattern of initializing tracing, loading AWS config,
// and starting the Lambda handler with instrumentation.
//
// Before (typical handler with ~15-20 lines of boilerplate):
//
//	func main() {
//	    ctx := context.Background()
//	    tp, err := tracing.Init(ctx)
//	    if err != nil {
//	        panic(err)
//	    }
//	    otel.SetTracerProvider(tp)
//	    ctx, coldStartSpan := tracing.StartColdStartSpan(ctx, "jmap-api")
//	    awsCfg, err := config.LoadDefaultConfig(ctx)
//	    if err != nil {
//	        coldStartSpan.End()
//	        panic(err)
//	    }
//	    otelaws.AppendMiddlewares(&awsCfg.APIOptions)
//	    ddb := dynamodb.NewFromConfig(awsCfg)
//	    handler := NewHandler(ddb)
//	    coldStartSpan.End()
//	    wrapped := otellambda.InstrumentHandler(handler.Handle, ...)
//	    lambda.Start(wrapped)
//	}
//
// After (with awsinit ~8 lines):
//
//	func main() {
//	    result, err := awsinit.Init(context.Background(),
//	        awsinit.WithHTTPHandler("jmap-api"),
//	    )
//	    if err != nil {
//	        panic(err)
//	    }
//	    defer result.Cleanup()
//
//	    ddb := dynamodb.NewFromConfig(result.Config)
//	    handler := NewHandler(ddb)
//	    result.Start(handler.Handle)
//	}
//
// For event-driven handlers (SQS, SNS, etc.) that don't need a cold start span:
//
//	func main() {
//	    result, err := awsinit.Init(context.Background())
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    sqsClient := sqs.NewFromConfig(result.Config)
//	    handler := NewHandler(sqsClient)
//	    result.Start(handler.Handle)
//	}
//
// The Result provides:
//   - Ctx: Context that may contain a cold start span (for HTTP handlers)
//   - Config: AWS config with OTel instrumentation middleware attached
//   - TracerProvider: The initialized TracerProvider for use with custom tracing
//   - Start(): Wraps handler with otellambda instrumentation and starts Lambda
//   - Cleanup(): Ends cold start span if one was created (safe to call always)
package awsinit
