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

## Development Process

All code changes **MUST** use the Test Driven Design (TDD) Superpower. Note in any plans that RED tests **MUST** compile, run, not panic, and fail - if they do not run, they do not prove we have testing.

All code must use suitable standard libraries wherever possible, in preference to writing things ourselves. All code must be clean and simple. All code must detect errors, and propagate them outward as distinct Error types - we use proper go Error methods for working with error types. Do not write fallbacks unless specifically required by the protocol in question - if a thing does not work, it errors. Use types and interfaces. If we are only going to use part of an interface from an external library, write our own sub-interface for interacting with it.
