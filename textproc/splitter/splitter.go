package splitter

import (
	"regexp"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor splits chunks that exceed the byte limit into smaller chunks.
// It tries sentence boundaries first, then word boundaries, then character boundaries.
type Processor struct {
	src       textproc.ChunkProcessor
	byteLimit int
	remaining textproc.Chunk
}

// Option configures a Processor.
type Option func(*Processor)

// WithByteLimit sets the maximum byte size for chunks.
func WithByteLimit(n int) Option {
	return func(p *Processor) {
		p.byteLimit = n
	}
}

// NewProcessor creates a new Processor with the given ChunkProcessor source.
// Chunks larger than maxBytes will be split into smaller chunks.
func NewProcessor(src textproc.ChunkProcessor, maxBytes int, opts ...Option) *Processor {
	p := &Processor{
		src:       src,
		byteLimit: maxBytes,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// trimLeading removes leading whitespace from a Chunk.
func trimLeading(chunk textproc.Chunk) textproc.Chunk {
	for i := 0; i < len(chunk); i++ {
		if chunk[i] != ' ' && chunk[i] != '\t' && chunk[i] != '\n' {
			return chunk[i:]
		}
	}
	return nil // All whitespace
}

// trimTrailing removes trailing whitespace from a Chunk.
func trimTrailing(chunk textproc.Chunk) textproc.Chunk {
	for i := len(chunk) - 1; i >= 0; i-- {
		if chunk[i] != ' ' && chunk[i] != '\t' && chunk[i] != '\n' {
			return chunk[:i+1]
		}
	}
	return nil // All whitespace
}

// Next returns the next Chunk, splitting large chunks if necessary.
func (p *Processor) Next() (textproc.Chunk, error) {
	// Return any remaining content first
	if len(p.remaining) > 0 {
		result := p.remaining
		p.remaining = nil
		result = trimLeading(result)
		if result == nil {
			return p.Next() // All whitespace, skip
		}
		if len(result) <= p.byteLimit {
			return result, nil
		}
		// If remaining is not much larger than limit, return as-is to avoid over-splitting
		// Otherwise, split it
		if len(result) <= p.byteLimit*2 {
			return result, nil
		}
		return p.splitChunk(result)
	}

	// Get next chunk from source
	chunk, err := p.src.Next()
	if err != nil {
		return nil, err
	}

	// Skip empty chunks
	if len(chunk) == 0 {
		return p.Next()
	}

	// Trim leading whitespace
	chunk = trimLeading(chunk)
	if chunk == nil {
		return p.Next() // All whitespace, skip
	}

	// If chunk fits, return it directly
	if len(chunk) <= p.byteLimit {
		return chunk, nil
	}

	// Chunk is too large, split it
	return p.splitChunk(chunk)
}

// splitChunk splits a chunk that exceeds the byte limit.
func (p *Processor) splitChunk(chunk textproc.Chunk) (textproc.Chunk, error) {
	// Try sentence boundary first
	splitIdx := p.findSentenceBoundary(chunk)
	if splitIdx > 0 {
		p.remaining = chunk[splitIdx:]
		return trimTrailing(chunk[:splitIdx]), nil
	}

	// Try word boundary
	splitIdx = p.findWordBoundary(chunk)
	if splitIdx > 0 {
		p.remaining = chunk[splitIdx:]
		return trimTrailing(chunk[:splitIdx]), nil
	}

	// No sentence/word boundaries found
	// Split at character boundary (UTF-8 safe)
	splitIdx = p.findCharacterBoundary(chunk)
	if splitIdx > 0 {
		p.remaining = chunk[splitIdx:]
		return chunk[:splitIdx], nil
	}

	// Worst case: return as-is
	return chunk, nil
}

// findSentenceBoundary finds the best sentence boundary within byteLimit.
// Returns the index where remaining content starts (after the delimiter).
// The delimiter ([.!?。] and following whitespace) is NOT included in the first piece.
func (p *Processor) findSentenceBoundary(chunk textproc.Chunk) int {
	// Sentence pattern: [.!?。] followed by whitespace
	sentencePattern := regexp.MustCompile(`[.!?。]\s+`)

	// Find all matches
	matches := sentencePattern.FindAllIndex(chunk, -1)
	if len(matches) == 0 {
		return 0
	}

	// Find the match that gives the first reasonable chunk within byteLimit
	for _, match := range matches {
		// match[0] is position of [.!?], match[1] is end of whitespace
		// First piece should be up to match[0]+1 (includes [.!?], excludes whitespace)
		// remaining starts at match[1] (after whitespace)
		endOfSentence := match[0] + 1 // Include the [.!?] in first piece
		if endOfSentence <= p.byteLimit && endOfSentence >= p.byteLimit/4 {
			return endOfSentence // Return position after [.!]?, before whitespace
		}
	}

	return 0
}

// findWordBoundary finds the last word boundary within byteLimit.
// Returns the index where remaining content starts (after the space).
// The space is NOT included in the first piece.
func (p *Processor) findWordBoundary(chunk textproc.Chunk) int {
	// Find the last space within byteLimit
	searchEnd := p.byteLimit
	if searchEnd > len(chunk) {
		searchEnd = len(chunk)
	}

	for i := searchEnd - 1; i >= 0; i-- {
		if chunk[i] == ' ' || chunk[i] == '\t' || chunk[i] == '\n' {
			// First piece is chunk[:i], remaining is chunk[i+1:]
			return i + 1 // Remaining starts after the space
		}
	}
	return 0
}

// findCharacterBoundary finds a safe UTF-8 character boundary at or before byteLimit.
// Returns the index where remaining content starts.
func (p *Processor) findCharacterBoundary(chunk textproc.Chunk) int {
	pos := p.byteLimit
	if pos > len(chunk) {
		pos = len(chunk)
	}

	// Find nearest valid UTF-8 boundary at or before pos
	for i := pos; i > 0; i-- {
		// Check if byte i is the start of a UTF-8 rune
		b := chunk[i]
		if b&0xC0 != 0x80 { // Not a continuation byte
			return i
		}
	}

	return 0
}
