package jmaperror

import (
	"errors"
	"testing"
)

// FuzzMethodError verifies that MethodError methods never panic on arbitrary input.
func FuzzMethodError(f *testing.F) {
	f.Add("unknownMethod", "test description")
	f.Add("", "")
	f.Add("serverFail", "internal error occurred")

	f.Fuzz(func(t *testing.T, errType, description string) {
		err := &MethodError{
			ErrType:     errType,
			Description: description,
			Err:         errors.New(description),
		}
		_ = err.Error()
		_ = err.Type()
		_ = err.ToMap()
		_ = err.Unwrap()

		// Also test with nil wrapped error
		err2 := &MethodError{
			ErrType:     errType,
			Description: description,
		}
		_ = err2.Error()
		_ = err2.Type()
		_ = err2.ToMap()
		_ = err2.Unwrap()
	})
}

// FuzzMethodErrorConstructors verifies that all constructor functions never panic.
func FuzzMethodErrorConstructors(f *testing.F) {
	f.Add("test description")
	f.Add("")
	f.Add("a]b[c{d}e")

	f.Fuzz(func(t *testing.T, description string) {
		constructors := []func(string) *MethodError{
			UnknownMethod,
			InvalidArguments,
			AccountNotFound,
			InvalidResultReference,
			StateMismatch,
			Forbidden,
			CannotCalculateChanges,
			UnsupportedFilter,
			UnsupportedSort,
			AnchorNotFound,
		}
		for _, ctor := range constructors {
			err := ctor(description)
			_ = err.Error()
			_ = err.Type()
			_ = err.ToMap()
		}

		// ServerFail takes an extra error arg
		sfErr := ServerFail(description, errors.New(description))
		_ = sfErr.Error()
		_ = sfErr.Type()
		_ = sfErr.ToMap()
		_ = sfErr.Unwrap()

		sfNil := ServerFail(description, nil)
		_ = sfNil.Error()
		_ = sfNil.ToMap()
		_ = sfNil.Unwrap()
	})
}

// FuzzSetError verifies that SetError methods never panic on arbitrary input.
func FuzzSetError(f *testing.F) {
	f.Add("notFound", "object not found", "prop1")
	f.Add("", "", "")

	f.Fuzz(func(t *testing.T, errType, description, prop string) {
		// With properties
		err := &SetError{
			ErrType:     errType,
			Description: description,
			Properties:  []string{prop},
		}
		_ = err.Error()
		_ = err.Type()
		_ = err.ToMap()

		// Without properties
		err2 := &SetError{
			ErrType:     errType,
			Description: description,
		}
		_ = err2.Error()
		_ = err2.Type()
		_ = err2.ToMap()
	})
}

// FuzzSetErrorConstructors verifies that all SetError constructor functions never panic.
func FuzzSetErrorConstructors(f *testing.F) {
	f.Add("test description", "prop1")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, description, prop string) {
		simpleConstructors := []func(string) *SetError{
			NotFound,
			TooLarge,
			OverQuota,
			TooManyPending,
			BlobNotFound,
			InvalidMailboxId,
			InvalidEmail,
			SetForbidden,
			InvalidPatch,
			MailboxHasEmail,
			SetServerFail,
		}
		for _, ctor := range simpleConstructors {
			err := ctor(description)
			_ = err.Error()
			_ = err.Type()
			_ = err.ToMap()
		}

		// InvalidProperties takes a properties slice
		ipErr := InvalidProperties(description, []string{prop})
		_ = ipErr.Error()
		_ = ipErr.Type()
		_ = ipErr.ToMap()

		ipNilErr := InvalidProperties(description, nil)
		_ = ipNilErr.ToMap()
	})
}

// FuzzHTTPProblem verifies that HTTPProblem methods never panic on arbitrary input.
func FuzzHTTPProblem(f *testing.F) {
	f.Add("urn:ietf:params:jmap:error:notJSON", "Not JSON", "bad json", 400, "")
	f.Add("", "", "", 0, "maxSize")

	f.Fuzz(func(t *testing.T, problemType, title, detail string, status int, limit string) {
		err := &HTTPProblem{
			ProblemType: problemType,
			Title:       title,
			Detail:      detail,
			Status:      status,
			Limit:       limit,
		}
		_ = err.Error()
		_ = err.Type()
		_ = err.ToMap()
	})
}

// FuzzHTTPProblemConstructors verifies that all HTTPProblem constructor functions never panic.
func FuzzHTTPProblemConstructors(f *testing.F) {
	f.Add("test detail", "limitName")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, detail, limitName string) {
		simpleConstructors := []func(string) *HTTPProblem{
			UnknownCapability,
			NotJSON,
			NotRequest,
		}
		for _, ctor := range simpleConstructors {
			err := ctor(detail)
			_ = err.Error()
			_ = err.Type()
			_ = err.ToMap()
		}

		// Limit takes two args
		lErr := Limit(limitName, detail)
		_ = lErr.Error()
		_ = lErr.Type()
		_ = lErr.ToMap()
	})
}
