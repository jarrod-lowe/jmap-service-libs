package jmaperror

import (
	"errors"
	"testing"
)

// Test that types implement JMAPError interface
func TestMethodErrorImplementsJMAPError(t *testing.T) {
	var _ JMAPError = &MethodError{}
}

func TestSetErrorImplementsJMAPError(t *testing.T) {
	var _ JMAPError = &SetError{}
}

func TestHTTPProblemImplementsJMAPError(t *testing.T) {
	var _ JMAPError = &HTTPProblem{}
}

// MethodError tests

func TestUnknownMethod(t *testing.T) {
	err := UnknownMethod("method Foo/bar not found")

	if err.Type() != "unknownMethod" {
		t.Errorf("Type() = %q, want %q", err.Type(), "unknownMethod")
	}
	if err.Error() != "unknownMethod: method Foo/bar not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "unknownMethod: method Foo/bar not found")
	}

	m := err.ToMap()
	if m["type"] != "unknownMethod" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "unknownMethod")
	}
	if m["description"] != "method Foo/bar not found" {
		t.Errorf("ToMap()[description] = %v, want %q", m["description"], "method Foo/bar not found")
	}
}

func TestInvalidArguments(t *testing.T) {
	err := InvalidArguments("mailboxId must be provided")

	if err.Type() != "invalidArguments" {
		t.Errorf("Type() = %q, want %q", err.Type(), "invalidArguments")
	}
	if err.Error() != "invalidArguments: mailboxId must be provided" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalidArguments: mailboxId must be provided")
	}

	m := err.ToMap()
	if m["type"] != "invalidArguments" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "invalidArguments")
	}
	if m["description"] != "mailboxId must be provided" {
		t.Errorf("ToMap()[description] = %v, want %q", m["description"], "mailboxId must be provided")
	}
}

func TestServerFail(t *testing.T) {
	underlying := errors.New("database connection failed")
	err := ServerFail("internal error occurred", underlying)

	if err.Type() != "serverFail" {
		t.Errorf("Type() = %q, want %q", err.Type(), "serverFail")
	}
	if err.Error() != "serverFail: internal error occurred" {
		t.Errorf("Error() = %q, want %q", err.Error(), "serverFail: internal error occurred")
	}
	if err.Unwrap() != underlying {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), underlying)
	}

	m := err.ToMap()
	if m["type"] != "serverFail" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "serverFail")
	}
	if m["description"] != "internal error occurred" {
		t.Errorf("ToMap()[description] = %v, want %q", m["description"], "internal error occurred")
	}
}

func TestServerFailWithNilError(t *testing.T) {
	err := ServerFail("internal error occurred", nil)

	if err.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", err.Unwrap())
	}
}

func TestAccountNotFound(t *testing.T) {
	err := AccountNotFound("account abc123 not found")

	if err.Type() != "accountNotFound" {
		t.Errorf("Type() = %q, want %q", err.Type(), "accountNotFound")
	}
	if err.Error() != "accountNotFound: account abc123 not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "accountNotFound: account abc123 not found")
	}

	m := err.ToMap()
	if m["type"] != "accountNotFound" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "accountNotFound")
	}
}

func TestInvalidResultReference(t *testing.T) {
	err := InvalidResultReference("result reference #1 is invalid")

	if err.Type() != "invalidResultReference" {
		t.Errorf("Type() = %q, want %q", err.Type(), "invalidResultReference")
	}
	if err.Error() != "invalidResultReference: result reference #1 is invalid" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalidResultReference: result reference #1 is invalid")
	}

	m := err.ToMap()
	if m["type"] != "invalidResultReference" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "invalidResultReference")
	}
}

func TestStateMismatch(t *testing.T) {
	err := StateMismatch("state does not match expected value")

	if err.Type() != "stateMismatch" {
		t.Errorf("Type() = %q, want %q", err.Type(), "stateMismatch")
	}
	if err.Error() != "stateMismatch: state does not match expected value" {
		t.Errorf("Error() = %q, want %q", err.Error(), "stateMismatch: state does not match expected value")
	}

	m := err.ToMap()
	if m["type"] != "stateMismatch" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "stateMismatch")
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("operation not permitted for this user")

	if err.Type() != "forbidden" {
		t.Errorf("Type() = %q, want %q", err.Type(), "forbidden")
	}
	if err.Error() != "forbidden: operation not permitted for this user" {
		t.Errorf("Error() = %q, want %q", err.Error(), "forbidden: operation not permitted for this user")
	}

	m := err.ToMap()
	if m["type"] != "forbidden" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "forbidden")
	}
}

func TestCannotCalculateChanges(t *testing.T) {
	err := CannotCalculateChanges("changes unavailable")

	if err.Type() != "cannotCalculateChanges" {
		t.Errorf("Type() = %q, want %q", err.Type(), "cannotCalculateChanges")
	}
	if err.Error() != "cannotCalculateChanges: changes unavailable" {
		t.Errorf("Error() = %q, want %q", err.Error(), "cannotCalculateChanges: changes unavailable")
	}

	m := err.ToMap()
	if m["type"] != "cannotCalculateChanges" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "cannotCalculateChanges")
	}
}

func TestUnsupportedFilter(t *testing.T) {
	err := UnsupportedFilter("filter not supported")

	if err.Type() != "unsupportedFilter" {
		t.Errorf("Type() = %q, want %q", err.Type(), "unsupportedFilter")
	}
	if err.Error() != "unsupportedFilter: filter not supported" {
		t.Errorf("Error() = %q, want %q", err.Error(), "unsupportedFilter: filter not supported")
	}

	m := err.ToMap()
	if m["type"] != "unsupportedFilter" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "unsupportedFilter")
	}
}

func TestUnsupportedSort(t *testing.T) {
	err := UnsupportedSort("sort not supported")

	if err.Type() != "unsupportedSort" {
		t.Errorf("Type() = %q, want %q", err.Type(), "unsupportedSort")
	}
	if err.Error() != "unsupportedSort: sort not supported" {
		t.Errorf("Error() = %q, want %q", err.Error(), "unsupportedSort: sort not supported")
	}

	m := err.ToMap()
	if m["type"] != "unsupportedSort" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "unsupportedSort")
	}
}

func TestAnchorNotFound(t *testing.T) {
	err := AnchorNotFound("anchor not found")

	if err.Type() != "anchorNotFound" {
		t.Errorf("Type() = %q, want %q", err.Type(), "anchorNotFound")
	}
	if err.Error() != "anchorNotFound: anchor not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "anchorNotFound: anchor not found")
	}

	m := err.ToMap()
	if m["type"] != "anchorNotFound" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "anchorNotFound")
	}
}

func TestMethodErrorUnwrapWithoutWrappedError(t *testing.T) {
	err := InvalidArguments("test")

	if err.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", err.Unwrap())
	}
}

// SetError tests

func TestNotFound(t *testing.T) {
	err := NotFound("object abc123 not found")

	if err.Type() != "notFound" {
		t.Errorf("Type() = %q, want %q", err.Type(), "notFound")
	}
	if err.Error() != "notFound: object abc123 not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "notFound: object abc123 not found")
	}

	m := err.ToMap()
	if m["type"] != "notFound" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "notFound")
	}
	if m["description"] != "object abc123 not found" {
		t.Errorf("ToMap()[description] = %v, want %q", m["description"], "object abc123 not found")
	}
	if _, ok := m["properties"]; ok {
		t.Errorf("ToMap() should not contain properties key")
	}
}

func TestInvalidProperties(t *testing.T) {
	err := InvalidProperties("invalid property values", []string{"name", "email"})

	if err.Type() != "invalidProperties" {
		t.Errorf("Type() = %q, want %q", err.Type(), "invalidProperties")
	}
	if err.Error() != "invalidProperties: invalid property values" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalidProperties: invalid property values")
	}

	m := err.ToMap()
	if m["type"] != "invalidProperties" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "invalidProperties")
	}
	if m["description"] != "invalid property values" {
		t.Errorf("ToMap()[description] = %v, want %q", m["description"], "invalid property values")
	}

	props, ok := m["properties"].([]string)
	if !ok {
		t.Fatalf("ToMap()[properties] is not []string, got %T", m["properties"])
	}
	if len(props) != 2 || props[0] != "name" || props[1] != "email" {
		t.Errorf("ToMap()[properties] = %v, want [name email]", props)
	}
}

func TestInvalidPropertiesEmptyList(t *testing.T) {
	err := InvalidProperties("invalid property values", []string{})

	m := err.ToMap()
	props, ok := m["properties"].([]string)
	if !ok {
		t.Fatalf("ToMap()[properties] is not []string, got %T", m["properties"])
	}
	if len(props) != 0 {
		t.Errorf("ToMap()[properties] = %v, want empty slice", props)
	}
}

func TestTooLarge(t *testing.T) {
	err := TooLarge("object exceeds maximum size")

	if err.Type() != "tooLarge" {
		t.Errorf("Type() = %q, want %q", err.Type(), "tooLarge")
	}
	if err.Error() != "tooLarge: object exceeds maximum size" {
		t.Errorf("Error() = %q, want %q", err.Error(), "tooLarge: object exceeds maximum size")
	}

	m := err.ToMap()
	if m["type"] != "tooLarge" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "tooLarge")
	}
}

func TestOverQuota(t *testing.T) {
	err := OverQuota("account storage quota exceeded")

	if err.Type() != "overQuota" {
		t.Errorf("Type() = %q, want %q", err.Type(), "overQuota")
	}
	if err.Error() != "overQuota: account storage quota exceeded" {
		t.Errorf("Error() = %q, want %q", err.Error(), "overQuota: account storage quota exceeded")
	}

	m := err.ToMap()
	if m["type"] != "overQuota" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "overQuota")
	}
}

func TestTooManyPending(t *testing.T) {
	err := TooManyPending("too many pending submissions")

	if err.Type() != "tooManyPending" {
		t.Errorf("Type() = %q, want %q", err.Type(), "tooManyPending")
	}
	if err.Error() != "tooManyPending: too many pending submissions" {
		t.Errorf("Error() = %q, want %q", err.Error(), "tooManyPending: too many pending submissions")
	}

	m := err.ToMap()
	if m["type"] != "tooManyPending" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "tooManyPending")
	}
}

func TestBlobNotFound(t *testing.T) {
	err := BlobNotFound("blob xyz789 not found")

	if err.Type() != "blobNotFound" {
		t.Errorf("Type() = %q, want %q", err.Type(), "blobNotFound")
	}
	if err.Error() != "blobNotFound: blob xyz789 not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "blobNotFound: blob xyz789 not found")
	}

	m := err.ToMap()
	if m["type"] != "blobNotFound" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "blobNotFound")
	}
}

func TestInvalidMailboxId(t *testing.T) {
	err := InvalidMailboxId("mailbox ID is not valid")

	if err.Type() != "invalidMailboxId" {
		t.Errorf("Type() = %q, want %q", err.Type(), "invalidMailboxId")
	}
	if err.Error() != "invalidMailboxId: mailbox ID is not valid" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalidMailboxId: mailbox ID is not valid")
	}

	m := err.ToMap()
	if m["type"] != "invalidMailboxId" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "invalidMailboxId")
	}
}

func TestInvalidEmail(t *testing.T) {
	err := InvalidEmail("email content is malformed")

	if err.Type() != "invalidEmail" {
		t.Errorf("Type() = %q, want %q", err.Type(), "invalidEmail")
	}
	if err.Error() != "invalidEmail: email content is malformed" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalidEmail: email content is malformed")
	}

	m := err.ToMap()
	if m["type"] != "invalidEmail" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "invalidEmail")
	}
}

func TestSetForbidden(t *testing.T) {
	err := SetForbidden("operation not allowed")

	if err.Type() != "forbidden" {
		t.Errorf("Type() = %q, want %q", err.Type(), "forbidden")
	}
	if err.Error() != "forbidden: operation not allowed" {
		t.Errorf("Error() = %q, want %q", err.Error(), "forbidden: operation not allowed")
	}

	m := err.ToMap()
	if m["type"] != "forbidden" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "forbidden")
	}
}

func TestInvalidPatch(t *testing.T) {
	err := InvalidPatch("invalid JSON Pointer")

	if err.Type() != "invalidPatch" {
		t.Errorf("Type() = %q, want %q", err.Type(), "invalidPatch")
	}
	if err.Error() != "invalidPatch: invalid JSON Pointer" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalidPatch: invalid JSON Pointer")
	}

	m := err.ToMap()
	if m["type"] != "invalidPatch" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "invalidPatch")
	}
}

func TestMailboxHasEmail(t *testing.T) {
	err := MailboxHasEmail("mailbox contains emails")

	if err.Type() != "mailboxHasEmail" {
		t.Errorf("Type() = %q, want %q", err.Type(), "mailboxHasEmail")
	}
	if err.Error() != "mailboxHasEmail: mailbox contains emails" {
		t.Errorf("Error() = %q, want %q", err.Error(), "mailboxHasEmail: mailbox contains emails")
	}

	m := err.ToMap()
	if m["type"] != "mailboxHasEmail" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "mailboxHasEmail")
	}
}

func TestSetServerFail(t *testing.T) {
	err := SetServerFail("internal server error")

	if err.Type() != "serverFail" {
		t.Errorf("Type() = %q, want %q", err.Type(), "serverFail")
	}
	if err.Error() != "serverFail: internal server error" {
		t.Errorf("Error() = %q, want %q", err.Error(), "serverFail: internal server error")
	}

	m := err.ToMap()
	if m["type"] != "serverFail" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "serverFail")
	}
}

// HTTPProblem tests

func TestUnknownCapability(t *testing.T) {
	err := UnknownCapability("capability urn:example:foo not supported")

	if err.Type() != "urn:ietf:params:jmap:error:unknownCapability" {
		t.Errorf("Type() = %q, want %q", err.Type(), "urn:ietf:params:jmap:error:unknownCapability")
	}
	if err.Error() != "unknownCapability: capability urn:example:foo not supported" {
		t.Errorf("Error() = %q, want %q", err.Error(), "unknownCapability: capability urn:example:foo not supported")
	}

	m := err.ToMap()
	if m["type"] != "urn:ietf:params:jmap:error:unknownCapability" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "urn:ietf:params:jmap:error:unknownCapability")
	}
	if m["title"] != "Unknown Capability" {
		t.Errorf("ToMap()[title] = %v, want %q", m["title"], "Unknown Capability")
	}
	if m["detail"] != "capability urn:example:foo not supported" {
		t.Errorf("ToMap()[detail] = %v, want %q", m["detail"], "capability urn:example:foo not supported")
	}
	if m["status"] != 400 {
		t.Errorf("ToMap()[status] = %v, want %d", m["status"], 400)
	}
	if _, ok := m["limit"]; ok {
		t.Errorf("ToMap() should not contain limit key")
	}
}

func TestNotJSON(t *testing.T) {
	err := NotJSON("request body is not valid JSON")

	if err.Type() != "urn:ietf:params:jmap:error:notJSON" {
		t.Errorf("Type() = %q, want %q", err.Type(), "urn:ietf:params:jmap:error:notJSON")
	}
	if err.Error() != "notJSON: request body is not valid JSON" {
		t.Errorf("Error() = %q, want %q", err.Error(), "notJSON: request body is not valid JSON")
	}

	m := err.ToMap()
	if m["type"] != "urn:ietf:params:jmap:error:notJSON" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "urn:ietf:params:jmap:error:notJSON")
	}
	if m["title"] != "Not JSON" {
		t.Errorf("ToMap()[title] = %v, want %q", m["title"], "Not JSON")
	}
	if m["status"] != 400 {
		t.Errorf("ToMap()[status] = %v, want %d", m["status"], 400)
	}
}

func TestNotRequest(t *testing.T) {
	err := NotRequest("missing methodCalls field")

	if err.Type() != "urn:ietf:params:jmap:error:notRequest" {
		t.Errorf("Type() = %q, want %q", err.Type(), "urn:ietf:params:jmap:error:notRequest")
	}
	if err.Error() != "notRequest: missing methodCalls field" {
		t.Errorf("Error() = %q, want %q", err.Error(), "notRequest: missing methodCalls field")
	}

	m := err.ToMap()
	if m["type"] != "urn:ietf:params:jmap:error:notRequest" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "urn:ietf:params:jmap:error:notRequest")
	}
	if m["title"] != "Not Request" {
		t.Errorf("ToMap()[title] = %v, want %q", m["title"], "Not Request")
	}
	if m["status"] != 400 {
		t.Errorf("ToMap()[status] = %v, want %d", m["status"], 400)
	}
}

func TestLimit(t *testing.T) {
	err := Limit("maxSizeRequest", "request exceeds maximum size of 10MB")

	if err.Type() != "urn:ietf:params:jmap:error:limit" {
		t.Errorf("Type() = %q, want %q", err.Type(), "urn:ietf:params:jmap:error:limit")
	}
	if err.Error() != "limit: request exceeds maximum size of 10MB" {
		t.Errorf("Error() = %q, want %q", err.Error(), "limit: request exceeds maximum size of 10MB")
	}

	m := err.ToMap()
	if m["type"] != "urn:ietf:params:jmap:error:limit" {
		t.Errorf("ToMap()[type] = %v, want %q", m["type"], "urn:ietf:params:jmap:error:limit")
	}
	if m["title"] != "Limit Exceeded" {
		t.Errorf("ToMap()[title] = %v, want %q", m["title"], "Limit Exceeded")
	}
	if m["detail"] != "request exceeds maximum size of 10MB" {
		t.Errorf("ToMap()[detail] = %v, want %q", m["detail"], "request exceeds maximum size of 10MB")
	}
	if m["status"] != 400 {
		t.Errorf("ToMap()[status] = %v, want %d", m["status"], 400)
	}
	if m["limit"] != "maxSizeRequest" {
		t.Errorf("ToMap()[limit] = %v, want %q", m["limit"], "maxSizeRequest")
	}
}

func TestLimitEmptyLimitName(t *testing.T) {
	err := Limit("", "some limit exceeded")

	m := err.ToMap()
	if _, ok := m["limit"]; ok {
		t.Errorf("ToMap() should not contain limit key when limitName is empty")
	}
}

// Test that errors work with errors.Is/As

func TestMethodErrorWorksWithErrorsAs(t *testing.T) {
	underlying := errors.New("database error")
	err := ServerFail("server failure", underlying)

	var target *MethodError
	if !errors.As(err, &target) {
		t.Errorf("errors.As should match *MethodError")
	}
	if target.Type() != "serverFail" {
		t.Errorf("Type() = %q, want %q", target.Type(), "serverFail")
	}
}

func TestMethodErrorWorksWithErrorsIs(t *testing.T) {
	underlying := errors.New("database error")
	err := ServerFail("server failure", underlying)

	if !errors.Is(err, underlying) {
		t.Errorf("errors.Is should find wrapped error")
	}
}
