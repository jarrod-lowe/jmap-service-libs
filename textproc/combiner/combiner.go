package combiner

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor combines chunks with overlap into ChunkSlices.
type Processor struct {
	chunks       []textproc.Chunk        // Deprecated: use src instead
	src          textproc.ChunkProcessor // NEW: pull-based source
	index        int
	overlapCount int
	history      []textproc.Chunk
}

// Option configures a Processor.
type Option func(*Processor)

// New creates a new Processor.
// Deprecated: Use NewProcessor with ChunkProcessor for pull-based composition.
func New(chunks []textproc.Chunk, opts ...Option) *Processor {
	p := &Processor{chunks: chunks, index: 0}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// NewProcessor creates a new Processor with the given ChunkProcessor source.
// Overlaps chunks by including overlapCount previous chunks in each ChunkSlice.
func NewProcessor(src textproc.ChunkProcessor, overlapCount int, opts ...Option) *Processor {
	p := &Processor{
		src:          src,
		overlapCount: overlapCount,
		history:      make([]textproc.Chunk, 0, overlapCount+1),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next returns the next ChunkSlice with overlapping chunks.
func (p *Processor) Next() (textproc.ChunkSlice, error) {
	// Use pull-based source if available
	if p.src != nil {
		return p.nextFromSource()
	}

	// Fall back to []Chunk for backward compatibility
	if p.index >= len(p.chunks) {
		return nil, io.EOF
	}
	result := textproc.ChunkSlice{p.chunks[p.index]}
	p.index++
	return result, nil
}

// nextFromSource pulls chunks from source and returns them with overlap.
func (p *Processor) nextFromSource() (textproc.ChunkSlice, error) {
	chunk, err := p.src.Next()
	if err != nil {
		return nil, err
	}

	// Add current chunk to history
	p.history = append(p.history, chunk)

	// Build result with overlap
	var result textproc.ChunkSlice
	start := 0
	if len(p.history) > p.overlapCount {
		start = len(p.history) - p.overlapCount - 1
	}
	result = p.history[start:]

	return result, nil
}
