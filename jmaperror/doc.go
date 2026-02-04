// Package jmaperror provides type-safe JMAP error handling per RFC 8620.
//
// This package defines three error types for different levels of the JMAP protocol:
//
//   - MethodError: Method-level failures returned in methodResponses as ["error", {...}, "clientId"]
//   - SetError: Per-object failures in Foo/set operations (notCreated, notUpdated, notDestroyed)
//   - HTTPProblem: Request-level failures returned as application/problem+json per RFC 7807
//
// # MethodError Example
//
//	err := jmaperror.InvalidArguments("mailboxId must be provided")
//	// Returns: {"type": "invalidArguments", "description": "mailboxId must be provided"}
//
// # SetError Example
//
//	err := jmaperror.InvalidProperties("invalid property values", []string{"name", "email"})
//	// Returns: {"type": "invalidProperties", "description": "...", "properties": ["name", "email"]}
//
// # HTTPProblem Example
//
//	err := jmaperror.NotJSON("request body is not valid JSON")
//	// Returns: {"type": "urn:ietf:params:jmap:error:notJSON", "title": "Not JSON", ...}
package jmaperror
