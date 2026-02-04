# jmap-service-libs

Shared Go libraries for the jmap-service-* family of repositories.

## Installation

```bash
go get github.com/jarrod-lowe/jmap-service-libs
```

## Available Packages

*No packages available yet. Packages will be added as they are migrated from sibling services.*

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
