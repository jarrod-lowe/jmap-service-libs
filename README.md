# jmap-service-libs

Shared Go libraries for the jmap-service-* family of repositories.

## Installation

```bash
go get github.com/jarrod-lowe/jmap-service-libs
```

## Available Packages

### tracing

OpenTelemetry/X-Ray instrumentation for JMAP services running on AWS Lambda.

```go
import "github.com/jarrod-lowe/jmap-service-libs/tracing"
```

Features:

- Automatic X-Ray and W3C Trace Context propagation
- Standardized attribute helpers: `RequestID`, `AccountID`, `BlobID`, `ParentBlobID`, `ContentType`, `Function`, `JMAPMethod`, `JMAPClientID`, `JMAPCallIndex`
- Span creation helpers: `StartHandlerSpan`, `StartMethodSpan`, `StartColdStartSpan`
- Error recording with proper status codes: `RecordError`
- Convenience tracer wrapper: `Tracer`

### logging

Structured JSON logging for Lambda environments using `slog`.

```go
import "github.com/jarrod-lowe/jmap-service-libs/logging"

// Simple usage - reads LOG_LEVEL from environment (DEBUG, INFO, WARN, ERROR)
var logger = logging.New()

// Override log level programmatically
var debugLogger = logging.New(logging.WithLevel(slog.LevelDebug))

// Capture output for testing
var buf bytes.Buffer
testLogger := logging.New(logging.WithOutput(&buf))
```

Features:

- JSON output format for CloudWatch Logs
- Environment-based log level via `LOG_LEVEL` (defaults to INFO)
- Option pattern for overriding level or output
- Zero dependencies beyond standard library

## Planned Migrations

The following code patterns have been identified across `jmap-service-core` and `jmap-service-email` as candidates for migration to this shared library.

### High Priority — Identical or Near-Identical Code

These patterns exist in both repositories with minimal variation:

| Package | Description | Source Locations |
| ------- | ----------- | ---------------- |
| ~~`logging`~~ | ~~Structured JSON logging setup with `slog.NewJSONHandler`~~ | **Done** - see `logging` package |
| `awsinit` | AWS SDK config loading with OTel middleware instrumentation | All `main.go` files in both repos |
| `jmaperror` | JMAP protocol error response formatting, standard error type constants (`unknownMethod`, `invalidArguments`, `serverFail`, etc.) | All command handlers in both repos |
| `auth` | Account ID extraction from JWT claims and IAM path parameters, IAM authentication detection, ARN normalization, principal authorization | `jmap-service-core/cmd/*/main.go`, `jmap-service-core/internal/plugin/authorization.go` |
| `dbclient` | DynamoDB client interface definition, key prefix constants (`ACCOUNT#`, `META#`, etc.), conditional check error handling helpers | `jmap-service-core/internal/db/`, `jmap-service-email/internal/dynamo/`, repository files in both repos |
| `plugincontract` | JMAP plugin invocation request/response types (`PluginInvocationRequest`, `PluginInvocationResponse`, `MethodResponse`) | `jmap-service-core/pkg/plugincontract/` — verify `jmap-service-email` uses this |
| `apiresponse` | API Gateway proxy response formatting, HTTP error response helpers | `jmap-service-core/cmd/blob-upload/main.go`, `jmap-service-core/cmd/blob-download/main.go` |

### Medium Priority — Similar Patterns Requiring Abstraction

These patterns are similar but need some generalization:

| Package | Description | Source Locations |
| ------- | ----------- | ---------------- |
| `validation` | Media type (MIME) validation, AWS resource tag validation | `jmap-service-core/internal/bloballocate/allocate.go`, `jmap-service-core/cmd/blob-upload/main.go` |
| `arnutil` | ARN parsing, SQS ARN to Queue URL conversion | `jmap-service-core/cmd/account-init/main.go` |
| `httpclient` | HTTP client wrapper with exponential backoff retry | `jmap-service-email/internal/blob/client.go` |
| `sqspublish` | Generic SQS message publisher with JSON serialization | `jmap-service-email/internal/blobdelete/publisher.go`, `jmap-service-email/internal/mailboxcleanup/publisher.go` |
| `txerror` | DynamoDB `TransactionCanceledException` handling, conditional check failure detection and categorization | Multiple repository files in both repos |

### Future Candidates — Currently in One Repository

These patterns exist in only one repo but are likely needed as more services are added:

| Package | Description | Current Location | Rationale |
| ------- | ----------- | ---------------- | --------- |
| `resultref` | JMAP result reference resolution (RFC 8620 §3.7), JSON pointer evaluation with wildcard support | `jmap-service-core/internal/resultref/` | Any service processing JMAP method calls |
| `plugininvoke` | Lambda-based plugin invocation interface and implementation | `jmap-service-core/internal/plugin/invoker.go` | Core pattern for JMAP method dispatch |
| `pluginregistry` | Plugin metadata registry, method-to-Lambda routing, capability management | `jmap-service-core/internal/plugin/registry.go` | Central to JMAP multi-service architecture |
| `emailparse` | RFC 5322 email parsing, MIME structure extraction, body part handling | `jmap-service-email/internal/email/parser.go` | Any email-related service |
| `headers` | Email header parsing (RFC 2047 decoding, address list parsing, date parsing) | `jmap-service-email/internal/headers/` | Any email-related service |
| `charset` | Character set detection and decoding for email body content | `jmap-service-email/internal/charset/` | Any email-related service |
| `keywords` | JMAP Email keyword validation per RFC 8621 | `jmap-service-email/internal/email/keywords.go` | Any `Email/*` method handler |

### Architectural Patterns (Document Only)

These patterns should be documented as best practices but don't require code extraction:

- **Lambda Dependency Injection** — Handler struct with interface dependencies, initialized in `main()`, enables testability
- **Single-Table Design Keys** — `PK()` and `SK()` methods on domain types with consistent prefix constants
- **Repository Builder Pattern** — Methods returning `[]types.TransactWriteItem` for composable atomic transactions
- **Functional Mocks** — Test mocks with optional function pointers for flexible per-test stubbing
- **SQS/DynamoDB Streams Consumer Pattern** — Standard event loop with batch failure tracking

## Development

### Using Local Changes in Sibling Services

When developing changes to this library alongside a consuming service, use a `replace` directive in the consuming service's `go.mod`:

```go
replace github.com/jarrod-lowe/jmap-service-libs => ../jmap-service-libs
```

Remember to remove the replace directive before committing.

### Development Commands

```bash
make help     # Show all available targets
make deps     # Tidy dependencies
make test     # Run all tests
make lint     # Run golangci-lint
make fmt      # Format code
make clean    # Remove build artifacts
```

## Contributing

When adding a new package:

1. Create a directory at the root level (e.g., `tracing/`)
2. Include comprehensive tests
3. Add package documentation in a `doc.go` file
4. Update this README with the package description
