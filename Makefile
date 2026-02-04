.PHONY: help deps test lint fmt clean

help:
	@echo "jmap-service-libs - Shared Go libraries"
	@echo ""
	@echo "Available targets:"
	@echo "  make deps    - Fetch dependencies (go mod tidy)"
	@echo "  make test    - Run tests (go test -v ./...)"
	@echo "  make lint    - Run golangci-lint"
	@echo "  make fmt     - Format Go code (go fmt ./...)"
	@echo "  make clean   - Remove build artifacts"
	@echo ""

# Fetch and tidy dependencies
deps:
	@echo "Tidying Go module dependencies..."
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
# PATH includes ~/go/bin for go-installed tools
lint:
	@PATH="$(HOME)/go/bin:$$PATH"; \
	if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "ERROR: golangci-lint is not installed"; \
		echo "Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi; \
	echo "Running golangci-lint..."; \
	golangci-lint run ./...

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage.html
	rm -f *.test
	@echo "Clean complete."
