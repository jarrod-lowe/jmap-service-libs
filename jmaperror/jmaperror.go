package jmaperror

// JMAPError is the common interface for all JMAP errors.
type JMAPError interface {
	error
	Type() string
	ToMap() map[string]any
}

// MethodError represents method-level failures returned in methodResponses.
type MethodError struct {
	ErrType     string
	Description string
	Err         error
}

func (e *MethodError) Error() string {
	return e.ErrType + ": " + e.Description
}

func (e *MethodError) Unwrap() error {
	return e.Err
}

func (e *MethodError) Type() string {
	return e.ErrType
}

func (e *MethodError) ToMap() map[string]any {
	return map[string]any{
		"type":        e.ErrType,
		"description": e.Description,
	}
}

// UnknownMethod creates a MethodError for unknown method calls.
func UnknownMethod(description string) *MethodError {
	return &MethodError{
		ErrType:     "unknownMethod",
		Description: description,
	}
}

// InvalidArguments creates a MethodError for invalid method arguments.
func InvalidArguments(description string) *MethodError {
	return &MethodError{
		ErrType:     "invalidArguments",
		Description: description,
	}
}

// ServerFail creates a MethodError for server failures, wrapping the underlying error.
func ServerFail(description string, err error) *MethodError {
	return &MethodError{
		ErrType:     "serverFail",
		Description: description,
		Err:         err,
	}
}

// AccountNotFound creates a MethodError when the account is not found.
func AccountNotFound(description string) *MethodError {
	return &MethodError{
		ErrType:     "accountNotFound",
		Description: description,
	}
}

// InvalidResultReference creates a MethodError for invalid result references.
func InvalidResultReference(description string) *MethodError {
	return &MethodError{
		ErrType:     "invalidResultReference",
		Description: description,
	}
}

// StateMismatch creates a MethodError when state doesn't match.
func StateMismatch(description string) *MethodError {
	return &MethodError{
		ErrType:     "stateMismatch",
		Description: description,
	}
}

// Forbidden creates a MethodError when the operation is forbidden.
func Forbidden(description string) *MethodError {
	return &MethodError{
		ErrType:     "forbidden",
		Description: description,
	}
}

// CannotCalculateChanges creates a MethodError when changes cannot be calculated.
func CannotCalculateChanges(description string) *MethodError {
	return &MethodError{
		ErrType:     "cannotCalculateChanges",
		Description: description,
	}
}

// UnsupportedFilter creates a MethodError when the filter is not supported.
func UnsupportedFilter(description string) *MethodError {
	return &MethodError{
		ErrType:     "unsupportedFilter",
		Description: description,
	}
}

// UnsupportedSort creates a MethodError when the sort is not supported.
func UnsupportedSort(description string) *MethodError {
	return &MethodError{
		ErrType:     "unsupportedSort",
		Description: description,
	}
}

// AnchorNotFound creates a MethodError when the anchor is not found.
func AnchorNotFound(description string) *MethodError {
	return &MethodError{
		ErrType:     "anchorNotFound",
		Description: description,
	}
}

// SetError represents per-object failures in Foo/set operations.
type SetError struct {
	ErrType     string
	Description string
	Properties  []string
}

func (e *SetError) Error() string {
	return e.ErrType + ": " + e.Description
}

func (e *SetError) Type() string {
	return e.ErrType
}

func (e *SetError) ToMap() map[string]any {
	m := map[string]any{
		"type":        e.ErrType,
		"description": e.Description,
	}
	if e.Properties != nil {
		m["properties"] = e.Properties
	}
	return m
}

// NotFound creates a SetError when an object is not found.
func NotFound(description string) *SetError {
	return &SetError{
		ErrType:     "notFound",
		Description: description,
	}
}

// InvalidProperties creates a SetError for invalid property values.
func InvalidProperties(description string, properties []string) *SetError {
	return &SetError{
		ErrType:     "invalidProperties",
		Description: description,
		Properties:  properties,
	}
}

// TooLarge creates a SetError when an object is too large.
func TooLarge(description string) *SetError {
	return &SetError{
		ErrType:     "tooLarge",
		Description: description,
	}
}

// OverQuota creates a SetError when quota is exceeded.
func OverQuota(description string) *SetError {
	return &SetError{
		ErrType:     "overQuota",
		Description: description,
	}
}

// TooManyPending creates a SetError when there are too many pending operations.
func TooManyPending(description string) *SetError {
	return &SetError{
		ErrType:     "tooManyPending",
		Description: description,
	}
}

// BlobNotFound creates a SetError when a blob is not found.
func BlobNotFound(description string) *SetError {
	return &SetError{
		ErrType:     "blobNotFound",
		Description: description,
	}
}

// InvalidMailboxId creates a SetError for invalid mailbox IDs.
func InvalidMailboxId(description string) *SetError {
	return &SetError{
		ErrType:     "invalidMailboxId",
		Description: description,
	}
}

// InvalidEmail creates a SetError for invalid email content.
func InvalidEmail(description string) *SetError {
	return &SetError{
		ErrType:     "invalidEmail",
		Description: description,
	}
}

// SetForbidden creates a SetError when the operation is forbidden.
func SetForbidden(description string) *SetError {
	return &SetError{
		ErrType:     "forbidden",
		Description: description,
	}
}

// InvalidPatch creates a SetError for an invalid JSON Pointer patch.
func InvalidPatch(description string) *SetError {
	return &SetError{
		ErrType:     "invalidPatch",
		Description: description,
	}
}

// MailboxHasEmail creates a SetError when a mailbox cannot be deleted because it contains emails.
func MailboxHasEmail(description string) *SetError {
	return &SetError{
		ErrType:     "mailboxHasEmail",
		Description: description,
	}
}

// SetServerFail creates a SetError for server failures during set operations.
func SetServerFail(description string) *SetError {
	return &SetError{
		ErrType:     "serverFail",
		Description: description,
	}
}

// HTTPProblem represents request-level failures returned as application/problem+json.
type HTTPProblem struct {
	ProblemType string
	Title       string
	Detail      string
	Status      int
	Limit       string
}

func (e *HTTPProblem) Error() string {
	// Extract short name from URN for the error message
	shortName := e.ProblemType
	if len(e.ProblemType) > 0 {
		// Extract the last part after the last colon
		for i := len(e.ProblemType) - 1; i >= 0; i-- {
			if e.ProblemType[i] == ':' {
				shortName = e.ProblemType[i+1:]
				break
			}
		}
	}
	return shortName + ": " + e.Detail
}

func (e *HTTPProblem) Type() string {
	return e.ProblemType
}

func (e *HTTPProblem) ToMap() map[string]any {
	m := map[string]any{
		"type":   e.ProblemType,
		"title":  e.Title,
		"detail": e.Detail,
		"status": e.Status,
	}
	if e.Limit != "" {
		m["limit"] = e.Limit
	}
	return m
}

// UnknownCapability creates an HTTPProblem for unknown capabilities.
func UnknownCapability(detail string) *HTTPProblem {
	return &HTTPProblem{
		ProblemType: "urn:ietf:params:jmap:error:unknownCapability",
		Title:       "Unknown Capability",
		Detail:      detail,
		Status:      400,
	}
}

// NotJSON creates an HTTPProblem when the request is not valid JSON.
func NotJSON(detail string) *HTTPProblem {
	return &HTTPProblem{
		ProblemType: "urn:ietf:params:jmap:error:notJSON",
		Title:       "Not JSON",
		Detail:      detail,
		Status:      400,
	}
}

// NotRequest creates an HTTPProblem when the request is not a valid JMAP request.
func NotRequest(detail string) *HTTPProblem {
	return &HTTPProblem{
		ProblemType: "urn:ietf:params:jmap:error:notRequest",
		Title:       "Not Request",
		Detail:      detail,
		Status:      400,
	}
}

// Limit creates an HTTPProblem when a limit is exceeded.
func Limit(limitName, detail string) *HTTPProblem {
	return &HTTPProblem{
		ProblemType: "urn:ietf:params:jmap:error:limit",
		Title:       "Limit Exceeded",
		Detail:      detail,
		Status:      400,
		Limit:       limitName,
	}
}
