package combiner

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor wraps each chunk in a ChunkSlice.
type Processor struct {
	chunks []textproc.Chunk
	index  int
}

// Option configures a Processor.
type Option func(*Processor)

// New creates a new Processor.
func New(chunks []textproc.Chunk, opts ...Option) *Processor {
	p := &Processor{chunks: chunks, index: 0}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next returns the next Chunk wrapped in a ChunkSlice.
func (p *Processor) Next() (textproc.ChunkSlice, error) {
	if p.index >= len(p.chunks) {
		return nil, io.EOF
	}
	result := textproc.ChunkSlice{p.chunks[p.index]}
	p.index++
	return result, nil
}
