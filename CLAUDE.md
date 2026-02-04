# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a shared Go library (`github.com/jarrod-lowe/jmap-service-libs`) for the jmap-service-* family of repositories. It provides common functionality used across multiple JMAP services.

## Build Commands

```bash
make deps     # Tidy dependencies (go mod tidy)
make test     # Run all tests (go test -v ./...)
make lint     # Run golangci-lint (must be installed)
make fmt      # Format code (go fmt ./...)
```

Run a single test:

```bash
go test -v -run TestName ./package/...
```

## Architecture

Packages are organized at the root level for clean import paths:

```go
import "github.com/jarrod-lowe/jmap-service-libs/tracing"
import "github.com/jarrod-lowe/jmap-service-libs/dbclient"
```

## Development with Sibling Services

When testing changes locally with a consuming service, add a replace directive to the consuming service's go.mod:

```go
replace github.com/jarrod-lowe/jmap-service-libs => ../jmap-service-libs
```

Remove the replace directive before committing.
