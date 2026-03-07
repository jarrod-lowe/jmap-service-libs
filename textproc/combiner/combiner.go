package combiner

import (
	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor combines chunks with overlap into ChunkSlices.
type Processor struct {
	src          textproc.ChunkProcessor
	overlapCount int
	history      []textproc.Chunk
}

// Option configures a Processor.
type Option func(*Processor)

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
	return p.nextFromSource()
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
