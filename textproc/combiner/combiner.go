package combiner

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor accumulates chunks with overlap into ChunkSlices up to a byte limit.
type Processor struct {
	src       textproc.ChunkProcessor
	byteLimit int
	overlap   int
	buffer    []textproc.Chunk
	byteCount int
	exhausted bool
	pending   textproc.Chunk // Chunk that was read but didn't fit
}

// Option configures a Processor.
type Option func(*Processor)

// WithByteLimit sets the maximum bytes per ChunkSlice. Default is 4000.
func WithByteLimit(n int) Option {
	return func(p *Processor) {
		p.byteLimit = n
	}
}

// WithOverlap sets the number of chunks to overlap between outputs. Default is 2.
func WithOverlap(n int) Option {
	return func(p *Processor) {
		p.overlap = n
	}
}

// NewProcessor creates a new Processor with the given ChunkProcessor source.
// By default: byteLimit=4000, overlap=2.
func NewProcessor(src textproc.ChunkProcessor, opts ...Option) *Processor {
	p := &Processor{
		src:       src,
		byteLimit: 4000,
		overlap:   2,
		buffer:    make([]textproc.Chunk, 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next returns the next ChunkSlice with accumulated chunks and overlap.
// Returns io.EOF when all chunks have been consumed.
func (p *Processor) Next() (textproc.ChunkSlice, error) {
	// If we already returned everything and no pending chunks, return EOF
	if p.exhausted && len(p.buffer) == 0 && p.pending == nil {
		return nil, io.EOF
	}

	// Handle overlap: keep last overlap chunks from previous output
	if len(p.buffer) > p.overlap {
		p.buffer = p.buffer[len(p.buffer)-p.overlap:]
		// Recalculate byte count after trimming
		p.byteCount = 0
		for _, c := range p.buffer {
			p.byteCount += len(c)
		}

		// Progress guarantee: if overlap alone exceeds limit, drop from front
		for p.byteCount > p.byteLimit && len(p.buffer) > 0 {
			p.buffer = p.buffer[1:]
			p.byteCount = 0
			for _, c := range p.buffer {
				p.byteCount += len(c)
			}
		}
	}

	// Add pending chunk if we have one
	if p.pending != nil {
		p.buffer = append(p.buffer, p.pending)
		p.byteCount += len(p.pending)
		p.pending = nil
	}

	// Accumulate chunks until byte limit is reached
	for {
		chunk, err := p.src.Next()
		if err == io.EOF {
			p.exhausted = true
			break
		}
		if err != nil {
			return nil, err
		}

		// If adding this chunk would exceed limit
		if p.byteCount+len(chunk) > p.byteLimit {
			// Special case: single chunk exceeds limit, return it anyway
			if len(p.buffer) == 0 {
				result := textproc.ChunkSlice{chunk}
				return result, nil
			}

			// Otherwise, save as pending and return what we have
			p.pending = chunk
			break
		}

		// Add chunk to buffer
		p.buffer = append(p.buffer, chunk)
		p.byteCount += len(chunk)
	}

	// Return buffer as result
	if len(p.buffer) == 0 {
		return nil, io.EOF
	}

	result := make(textproc.ChunkSlice, len(p.buffer))
	copy(result, p.buffer)

	// If source is exhausted, clear buffer for EOF on next call
	// Don't keep overlap on final output
	if p.exhausted {
		p.buffer = make([]textproc.Chunk, 0)
		p.byteCount = 0
	} else {
		// Keep overlap chunks for next iteration
		if len(p.buffer) > p.overlap {
			p.buffer = p.buffer[len(p.buffer)-p.overlap:]
		} else {
			p.buffer = make([]textproc.Chunk, 0)
		}

		// Recalculate byte count
		p.byteCount = 0
		for _, c := range p.buffer {
			p.byteCount += len(c)
		}
	}

	return result, nil
}
