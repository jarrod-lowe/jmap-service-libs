.PHONY: help deps test test-race lint fmt fmt-check fuzz vulncheck mod-check license-check apidiff clean setup setup-repo setup-branch-protection

FUZZ_TIME ?= 30s

help:
	@echo "jmap-service-libs - Shared Go libraries"
	@echo ""
	@echo "Available targets:"
	@echo "  make deps          - Fetch dependencies (go mod tidy)"
	@echo "  make test          - Run tests (go test -v ./...)"
	@echo "  make test-race     - Run tests with race detector"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make fmt           - Format Go code (go fmt ./...)"
	@echo "  make fmt-check     - Check formatting (fails if not gofmt'd)"
	@echo "  make fuzz          - Run fuzz tests (FUZZ_TIME=30s)"
	@echo "  make vulncheck     - Scan dependencies for known CVEs"
	@echo "  make mod-check     - Verify go.mod and go.sum are tidy"
	@echo "  make license-check - Check dependency license compatibility"
	@echo "  make apidiff       - Detect breaking API changes vs last tag"
	@echo "  make clean         - Remove build artifacts"
	@echo ""
	@echo "Repository setup (requires gh CLI and admin access):"
	@echo "  make setup                   - Run all repo setup targets"
	@echo "  make setup-repo              - Configure repo settings and Actions permissions"
	@echo "  make setup-branch-protection - Apply branch protection to main"
	@echo ""

# Fetch and tidy dependencies
deps:
	@echo "Tidying Go module dependencies..."
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -race ./...

# Run linter
# PATH includes ~/go/bin for go-installed tools
lint:
	@PATH="$(HOME)/go/bin:$$PATH"; \
	if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "ERROR: golangci-lint is not installed"; \
		echo "Install it with: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
		exit 1; \
	fi; \
	echo "Running golangci-lint..."; \
	golangci-lint run ./...

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Check formatting (fails if code is not gofmt'd)
fmt-check:
	@echo "Checking formatting..."
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "ERROR: The following files are not formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi; \
	echo "All files are formatted."

# Run fuzz tests in parallel with a time limit
fuzz:
	@echo "Running fuzz tests in parallel ($(FUZZ_TIME) per target)..."
	@TMPDIR=$$(mktemp -d) && \
	trap 'rm -rf "$$TMPDIR"' EXIT && \
	PIDS="" && \
	NAMES="" && \
	for pkg in plugincontract jmaperror; do \
		for fuzz_func in $$(go test -list 'Fuzz.*' ./$$pkg/... 2>/dev/null | grep '^Fuzz'); do \
			echo "  Starting $$pkg/$$fuzz_func..." && \
			go test -fuzz="^$$fuzz_func$$" -fuzztime=$(FUZZ_TIME) -parallel=1 ./$$pkg/... \
				> "$$TMPDIR/$$pkg-$$fuzz_func.out" 2>&1 & \
			PIDS="$$PIDS $$!" && \
			NAMES="$$NAMES $$pkg/$$fuzz_func" && \
			true; \
		done; \
	done && \
	echo "  All fuzz targets launched." && \
	FAILED=0 && \
	set -- $$NAMES && \
	for pid in $$PIDS; do \
		NAME=$$1 && shift && \
		PKG=$${NAME%/*} && \
		FUNC=$${NAME#*/} && \
		if wait $$pid; then \
			echo "  PASS: $$NAME"; \
		else \
			echo "  FAIL: $$NAME" && \
			echo "--- Output from $$NAME ---" && \
			cat "$$TMPDIR/$$PKG-$$FUNC.out" && \
			echo "--- End output ---" && \
			FAILED=$$((FAILED + 1)); \
		fi; \
	done && \
	if [ "$$FAILED" -gt 0 ]; then \
		echo "$$FAILED fuzz target(s) failed." && \
		exit 1; \
	fi && \
	echo "All fuzz targets passed."

# Scan dependencies for known CVEs
vulncheck:
	@PATH="$(HOME)/go/bin:$$PATH"; \
	if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "ERROR: govulncheck is not installed"; \
		echo "Install it with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
		exit 1; \
	fi; \
	echo "Scanning for known vulnerabilities..."; \
	govulncheck ./...

# Verify go.mod and go.sum are tidy
mod-check:
	@echo "Checking go.mod and go.sum are tidy..."
	go mod tidy
	@if ! git diff --exit-code go.mod go.sum; then \
		echo "ERROR: go.mod or go.sum are not tidy. Run 'go mod tidy' and commit the changes."; \
		exit 1; \
	fi
	@echo "go.mod and go.sum are tidy."

# Check dependency license compatibility
license-check:
	@PATH="$(HOME)/go/bin:$$PATH"; \
	if ! command -v go-licenses >/dev/null 2>&1; then \
		echo "ERROR: go-licenses is not installed"; \
		echo "Install it with: go install github.com/google/go-licenses@latest"; \
		exit 1; \
	fi; \
	echo "Checking dependency licenses..."; \
	GOTOOLCHAIN=local go-licenses check --ignore github.com/jarrod-lowe/jmap-service-libs ./... 2>&1

# Detect breaking API changes vs last tag
apidiff:
	@PATH="$(HOME)/go/bin:$$PATH"; \
	if ! command -v apidiff >/dev/null 2>&1; then \
		echo "ERROR: apidiff is not installed"; \
		echo "Install it with: go install golang.org/x/exp/cmd/apidiff@latest"; \
		exit 1; \
	fi; \
	LATEST_TAG=$$(git describe --tags --abbrev=0 2>/dev/null); \
	if [ -z "$$LATEST_TAG" ]; then \
		echo "No previous tags found, skipping apidiff."; \
		exit 0; \
	fi; \
	echo "Comparing API against $$LATEST_TAG..."; \
	MODULE=$$(go list -m); \
	INCOMPATIBLE=0; \
	for pkg in $$(go list ./...); do \
		SUFFIX=$${pkg#$$MODULE}; \
		OLD="$$MODULE$$SUFFIX@$$LATEST_TAG"; \
		echo "Checking $$pkg vs $$OLD..."; \
		RESULT=$$(apidiff "$$OLD" "$$pkg" 2>&1) || true; \
		if echo "$$RESULT" | grep -q "Incompatible changes:"; then \
			echo "$$RESULT"; \
			INCOMPATIBLE=1; \
		fi; \
	done; \
	if [ "$$INCOMPATIBLE" -eq 1 ]; then \
		echo "ERROR: Incompatible API changes detected."; \
		exit 1; \
	fi; \
	echo "No incompatible API changes detected."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage.html
	rm -f *.test
	@echo "Clean complete."

# Repository setup targets (require gh CLI and admin access)
setup: setup-repo setup-branch-protection

setup-repo:
	@echo "Configuring repository settings..."
	gh repo edit --delete-branch-on-merge --enable-auto-merge --enable-wiki=false
	gh api -X PUT repos/{owner}/{repo}/actions/permissions/workflow \
		-f default_workflow_permissions=read \
		-F can_approve_pull_request_reviews=true
	@echo "Repository settings applied."

setup-branch-protection:
	@echo "Applying branch protection to main..."
	gh api -X PUT repos/{owner}/{repo}/branches/main/protection \
		--input .github/branch-protection.json
	gh api -X POST repos/{owner}/{repo}/branches/main/protection/required_signatures
	@echo "Branch protection applied."
