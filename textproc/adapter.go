package textproc

// BytesToStringAdapter wraps a BytesProcessor and converts its []byte output to string.
// This provides the bridge between byte-based processors (reader, utf8clean) and string-based
// processors (htmlstrip, elider, chunker, splitter).
type BytesToStringAdapter struct {
	inner BytesProcessor
}

// NewBytesToStringAdapter creates a new adapter wrapping the given BytesProcessor.
func NewBytesToStringAdapter(inner BytesProcessor) *BytesToStringAdapter {
	return &BytesToStringAdapter{inner: inner}
}

// Next returns the next block of data as a string, converting from the inner BytesProcessor.
// Returns io.EOF when all data has been consumed.
func (a *BytesToStringAdapter) Next() (string, error) {
	b, err := a.inner.Next()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Ensure BytesToStringAdapter implements StringProcessor
var _ StringProcessor = (*BytesToStringAdapter)(nil)
