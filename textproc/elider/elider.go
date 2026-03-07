package elider

import (
	"io"
	"regexp"
	"strings"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// eliderMode represents the current processing state
type eliderMode int

const (
	modeNormal    eliderMode = iota
	modeURL                  // Inside a URL
	modeQuote                // In quoted reply section
	modeSignature            // After signature delimiter
)

// Processor reads strings from a source and returns them in blocks
// with noise (URLs, quoted text, signatures, non-words) removed.
type Processor struct {
	src       textproc.StringProcessor
	blockSize int

	// State machine
	mode   eliderMode
	input  string // Unprocessed input from source
	output strings.Builder

	// URL processing state
	urlBuf         strings.Builder
	skipURLContent bool // true when we've extracted domain and need to skip path/params

	// Non-word filtering
	uuidPattern    *regexp.Regexp
	hexPattern     *regexp.Regexp
	versionPattern *regexp.Regexp

	done bool
}

// Option configures a Processor.
type Option func(*Processor)

// WithBlockSize sets the block size for reading.
func WithBlockSize(n int) Option {
	return func(p *Processor) {
		p.blockSize = n
	}
}

// NewProcessor creates a new Processor with the given StringProcessor source.
// This enables pull-based lazy evaluation.
func NewProcessor(src textproc.StringProcessor, opts ...Option) *Processor {
	p := &Processor{
		src:            src,
		blockSize:      1024,
		uuidPattern:    regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`),
		hexPattern:     regexp.MustCompile(`\b[0-9a-fA-F]{16,}\b`),
		versionPattern: regexp.MustCompile(`\b[vV]\d+(?:\.\d+)*\b`),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next reads the next block of data from the source with noise elided.
// Returns io.EOF when all data has been consumed.
func (p *Processor) Next() (string, error) {
	for {
		// If we have enough output, return it
		if p.output.Len() >= p.blockSize {
			result := strings.TrimRight(p.output.String(), " \t\n\r")
			p.output.Reset()
			return result, nil
		}

		// If done and no buffered output, return EOF
		if p.done && p.output.Len() == 0 {
			return "", io.EOF
		}

		// Read more input
		if p.input == "" {
			block, err := p.src.Next()
			if err == io.EOF {
				p.done = true
				if p.output.Len() > 0 {
					result := strings.TrimRight(p.output.String(), " \t\n\r")
					p.output.Reset()
					return result, nil
				}
				return "", io.EOF
			}
			if err != nil {
				return "", err
			}
			p.input = block
		}

		// Process input
		p.processInput()
	}
}

// processInput processes the current input buffer
func (p *Processor) processInput() {
	// If we're in skip URL content mode, skip until whitespace
	if p.skipURLContent {
		for p.input != "" {
			if p.input[0] == ' ' || p.input[0] == '\t' || p.input[0] == '\n' || p.input[0] == '\r' {
				p.skipURLContent = false
				return
			}
			p.input = p.input[1:]
		}
		// Exhausted all input, continue skipping in next block
		return
	}

	for p.input != "" {
		// Re-check skipURLContent in case it was set during loop iteration
		if p.skipURLContent {
			return
		}
		switch p.mode {
		case modeNormal:
			p.processNormal()
		case modeURL:
			p.processURL()
		case modeQuote, modeSignature:
			// In quote or signature mode, elide all content
			p.input = ""
		}
	}
}

// processNormal handles normal mode processing
func (p *Processor) processNormal() {
	// Check for URL pattern http:// or https://
	if len(p.input) >= 7 && p.input[:7] == "http://" {
		p.mode = modeURL
		p.input = p.input[7:]
		return
	}
	if len(p.input) >= 8 && p.input[:8] == "https://" {
		p.mode = modeURL
		p.input = p.input[8:]
		return
	}

	// Check for signature delimiter "--" on its own line
	if len(p.input) >= 3 && p.input[:3] == "--\n" {
		p.mode = modeSignature
		p.input = "" // Elide all remaining input
		return
	}

	// Check for quote markers
	// RFC 5322: "On [date] [name] wrote:" pattern
	if len(p.input) >= 3 && p.input[:3] == "On " {
		// Look for "wrote:" after this
		wroteIdx := strings.Index(p.input, " wrote:")
		if wroteIdx > 0 && wroteIdx < 200 { // Reasonable limit for email headers
			p.mode = modeQuote
			p.input = "" // Elide all remaining input
			return
		}
	}

	// Common quote markers
	quoteMarkers := []string{
		"-----Original Message-----",
		"From:",
		"Sent:",
		"To:",
		"Subject:",
		"Cc:",
	}
	for _, marker := range quoteMarkers {
		if len(p.input) >= len(marker) && p.input[:len(marker)] == marker {
			p.mode = modeQuote
			p.input = "" // Elide all remaining input
			return
		}
	}

	// Check for UUID pattern and skip it
	uuidIdx := p.uuidPattern.FindStringIndex(p.input)
	if uuidIdx != nil && uuidIdx[0] == 0 {
		// Skip the UUID and any trailing whitespace
		p.input = p.input[uuidIdx[1]:]
		p.skipWhitespace()
		return
	}

	// Check for hex string pattern (16+ consecutive hex chars) and skip it
	hexIdx := p.hexPattern.FindStringIndex(p.input)
	if hexIdx != nil && hexIdx[0] == 0 {
		// Skip the hex string and any trailing whitespace
		p.input = p.input[hexIdx[1]:]
		p.skipWhitespace()
		return
	}

	// Check for version string pattern and skip it
	versionIdx := p.versionPattern.FindStringIndex(p.input)
	if versionIdx != nil && versionIdx[0] == 0 {
		// Skip the version string and any trailing whitespace
		p.input = p.input[versionIdx[1]:]
		p.skipWhitespace()
		return
	}

	// Copy character to output
	p.output.WriteByte(p.input[0])
	p.input = p.input[1:]
}

// processURL handles URL mode processing
func (p *Processor) processURL() {
	// URL ends at whitespace, query/fragment markers, or first slash after domain
	if p.input == "" {
		return
	}

	c := p.input[0]
	if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '?' || c == '#' || c == '/' {
		// URL ended, flush domain and return to normal mode
		p.output.WriteString(p.urlBuf.String())
		p.urlBuf.Reset()
		p.mode = modeNormal
		// Skip the current character (/ ? #) and set flag to skip rest of URL content
		p.input = p.input[1:]
		p.skipURLContent = true
		return
	}

	p.urlBuf.WriteByte(c)
	p.input = p.input[1:]
}

// skipWhitespace skips leading whitespace in the input
func (p *Processor) skipWhitespace() {
	for p.input != "" {
		c := p.input[0]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		p.input = p.input[1:]
	}
}

// Ensure Processor implements StringProcessor
var _ textproc.StringProcessor = (*Processor)(nil)
