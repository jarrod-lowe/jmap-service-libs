package testing

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// StubProcessor implements textproc.BytesProcessor for testing.
// It returns the provided data strings as byte slices.
type StubProcessor struct {
	data  []string
	index int
}

// NewStubProcessor creates a new StubProcessor with the given data.
func NewStubProcessor(data []string) *StubProcessor {
	return &StubProcessor{
		data:  data,
		index: 0,
	}
}

// Next returns the next data item as a byte slice.
// Returns io.EOF when all data has been consumed.
func (s *StubProcessor) Next() ([]byte, error) {
	if s.index >= len(s.data) {
		return nil, io.EOF
	}
	result := []byte(s.data[s.index])
	s.index++
	return result, nil
}

// NewChunkSlice creates a textproc.ChunkSlice from the provided strings.
func NewChunkSlice(s ...string) textproc.ChunkSlice {
	cs := make(textproc.ChunkSlice, len(s))
	for i, str := range s {
		cs[i] = textproc.Chunk(str)
	}
	return cs
}
