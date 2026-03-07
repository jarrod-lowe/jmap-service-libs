package splitter

import (
	"regexp"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor splits chunks that exceed the character limit into smaller chunks.
// It tries sentence boundaries first, then word boundaries, then character boundaries.
type Processor struct {
	src       textproc.ChunkProcessor
	charLimit int
	remaining textproc.Chunk
}

// Option configures a Processor.
type Option func(*Processor)

// WithCharLimit sets the maximum character size for chunks.
func WithCharLimit(n int) Option {
	return func(p *Processor) {
		p.charLimit = n
	}
}

// NewProcessor creates a new Processor with the given ChunkProcessor source.
// Chunks larger than maxChars will be split into smaller chunks.
func NewProcessor(src textproc.ChunkProcessor, maxChars int, opts ...Option) *Processor {
	p := &Processor{
		src:       src,
		charLimit: maxChars,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// trimLeading removes leading whitespace from a Chunk.
func trimLeading(chunk textproc.Chunk) textproc.Chunk {
	for i, r := range chunk {
		if r != ' ' && r != '\t' && r != '\n' {
			return chunk[i:]
		}
	}
	return "" // All whitespace
}

// trimTrailing removes trailing whitespace from a Chunk.
func trimTrailing(chunk textproc.Chunk) textproc.Chunk {
	for i := len(chunk) - 1; i >= 0; i-- {
		if chunk[i] != ' ' && chunk[i] != '\t' && chunk[i] != '\n' {
			return chunk[:i+1]
		}
	}
	return "" // All whitespace
}

// Next returns the next Chunk, splitting large chunks if necessary.
func (p *Processor) Next() (textproc.Chunk, error) {
	// Return any remaining content first
	if len(p.remaining) > 0 {
		result := p.remaining
		p.remaining = ""
		result = trimLeading(result)
		if result == "" {
			return p.Next() // All whitespace, skip
		}
		// Use character count for limit check
		if len([]rune(result)) <= p.charLimit {
			return result, nil
		}
		// If remaining is not much larger than limit, return as-is to avoid over-splitting
		// Otherwise, split it
		if len([]rune(result)) <= p.charLimit*2 {
			return result, nil
		}
		return p.splitChunk(result)
	}

	// Get next chunk from source
	chunk, err := p.src.Next()
	if err != nil {
		return "", err
	}

	// Skip empty chunks
	if chunk == "" {
		return p.Next()
	}

	// Trim leading whitespace
	chunk = trimLeading(chunk)
	if chunk == "" {
		return p.Next() // All whitespace, skip
	}

	// If chunk fits, return it directly
	if len([]rune(chunk)) <= p.charLimit {
		return chunk, nil
	}

	// Chunk is too large, split it
	return p.splitChunk(chunk)
}

// splitChunk splits a chunk that exceeds the char limit.
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

// findSentenceBoundary finds the best sentence boundary within charLimit.
// Returns the index where remaining content starts (after the delimiter).
// The delimiter ([.!?。] and following whitespace) is NOT included in the first piece.
func (p *Processor) findSentenceBoundary(chunk textproc.Chunk) int {
	// Sentence pattern: [.!?。] followed by whitespace
	sentencePattern := regexp.MustCompile(`[.!?。]\s+`)

	// Find all matches in the string
	matches := sentencePattern.FindAllStringIndex(string(chunk), -1)
	if len(matches) == 0 {
		return 0
	}

	// Find the match that gives the first reasonable chunk within charLimit
	for _, match := range matches {
		// match[0] is position of [.!?], match[1] is end of whitespace
		// First piece should be up to match[0]+1 (includes [.!?], excludes whitespace)
		// remaining starts at match[1] (after whitespace)
		endOfSentence := match[0] + 1 // Include the [.!?] in first piece
		// Check character count up to this point
		if len([]rune(chunk[:endOfSentence])) <= p.charLimit && len([]rune(chunk[:endOfSentence])) >= p.charLimit/4 {
			return endOfSentence // Return position after [.!?], before whitespace
		}
	}

	return 0
}

// findWordBoundary finds the last word boundary within charLimit.
// Returns the index where remaining content starts (after the space).
// The space is NOT included in the first piece.
func (p *Processor) findWordBoundary(chunk textproc.Chunk) int {
	runes := []rune(chunk)
	searchEnd := p.charLimit
	if searchEnd > len(runes) {
		searchEnd = len(runes)
	}

	// Find the last space within charLimit
	for i := searchEnd - 1; i >= 0; i-- {
		if runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n' {
			// First piece is chunk up to this rune, remaining is after this rune
			// Convert rune index back to byte index
			return len(string(runes[:i+1]))
		}
	}
	return 0
}

// findCharacterBoundary finds a safe UTF-8 character boundary at or before charLimit.
// Returns the index where remaining content starts.
func (p *Processor) findCharacterBoundary(chunk textproc.Chunk) int {
	runes := []rune(chunk)
	pos := p.charLimit
	if pos > len(runes) {
		pos = len(runes)
	}

	// Return the byte index for the rune position
	return len(string(runes[:pos]))
}
