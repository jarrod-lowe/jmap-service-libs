package htmlstrip

import (
	"io"
	"strings"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
	"golang.org/x/net/html"
)

// Processor reads strings from a source and returns them in blocks
// with HTML markup removed.
type Processor struct {
	src       textproc.StringProcessor
	blockSize int
	tokenizer *html.Tokenizer
	buf       strings.Builder
	done      bool
	skipTag   string // tag name being skipped ("script", "style", or "")
}

// Option configures a Processor.
type Option func(*Processor)

// WithBlockSize sets the block size for reading.
func WithBlockSize(n int) Option {
	return func(p *Processor) {
		p.blockSize = n
	}
}

// processorReader adapts a StringProcessor to an io.Reader for use with html.Tokenizer.
type processorReader struct {
	proc *Processor
	buf  []byte
	pos  int
}

// Read implements io.Reader by pulling from the StringProcessor.
func (pr *processorReader) Read(p []byte) (n int, err error) {
	// If buffer is empty, pull more data
	if pr.pos >= len(pr.buf) {
		block, err := pr.proc.src.Next()
		if err != nil {
			return 0, err
		}
		pr.buf = []byte(block)
		pr.pos = 0
	}

	// Copy from buffer to p
	n = copy(p, pr.buf[pr.pos:])
	pr.pos += n
	return n, nil
}

// NewProcessor creates a new Processor with the given StringProcessor source.
// This enables pull-based lazy evaluation.
func NewProcessor(src textproc.StringProcessor, opts ...Option) *Processor {
	p := &Processor{
		src:       src,
		blockSize: 1024,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next reads the next block of data from the source with HTML removed.
// Returns io.EOF when all data has been consumed.
func (p *Processor) Next() (string, error) {
	// Initialize tokenizer from pull-based source if needed
	if p.tokenizer == nil {
		p.tokenizer = html.NewTokenizer(&processorReader{proc: p})
	}

	if p.done && p.buf.Len() == 0 {
		return "", io.EOF
	}

	for {
		if p.buf.Len() >= p.blockSize {
			result := strings.TrimRight(p.buf.String(), " \t\n\r")
			p.buf.Reset()
			return result, nil
		}

		tokenType := p.tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := p.tokenizer.Err()
			if err == io.EOF {
				p.done = true
				if p.buf.Len() > 0 {
					result := strings.TrimRight(p.buf.String(), " \t\n\r")
					p.buf.Reset()
					return result, nil
				}
				return "", io.EOF
			}
			return "", err

		case html.TextToken:
			// Only add text if we're not inside a script or style tag
			if p.skipTag == "" {
				text := string(p.tokenizer.Text())
				p.buf.WriteString(text)
			}

		case html.StartTagToken:
			tagName, _ := p.tokenizer.TagName()
			tag := string(tagName)
			// Check for script/style tags
			if tag == "script" || tag == "style" {
				if p.skipTag == "" {
					p.skipTag = tag
				}
			} else if tag == "img" {
				// Extract alt attribute from img tags
				for {
					key, val, more := p.tokenizer.TagAttr()
					if string(key) == "alt" && string(val) != "" {
						p.buf.WriteString(string(val))
						break
					}
					if !more {
						break
					}
				}
			}

		case html.EndTagToken:
			tagName, _ := p.tokenizer.TagName()
			tag := string(tagName)
			// Check if we're closing a script or style tag
			if tag == p.skipTag {
				p.skipTag = ""
			} else if p.isCellTag(tag) {
				// Insert tab separator after table cells
				p.buf.WriteString("\t")
			} else if p.isBlockTag(tag) {
				// Insert newline after block elements
				p.buf.WriteString("\n")
			}

		case html.SelfClosingTagToken:
			tagName, _ := p.tokenizer.TagName()
			tag := string(tagName)
			if tag == "img" {
				// Extract alt attribute from img tags
				for {
					key, val, more := p.tokenizer.TagAttr()
					if string(key) == "alt" && string(val) != "" {
						p.buf.WriteString(string(val))
						break
					}
					if !more {
						break
					}
				}
			} else if tag == "br" || tag == "hr" {
				// Insert newline after br and hr
				p.buf.WriteString("\n")
			}

		case html.CommentToken, html.DoctypeToken:
			// Skip these tokens
		}
	}
}

// isCellTag returns true if the tag is a table cell element
func (p *Processor) isCellTag(tag string) bool {
	return tag == "td" || tag == "th"
}

// isBlockTag returns true if the tag is a block element that should insert a newline
func (p *Processor) isBlockTag(tag string) bool {
	switch tag {
	case "p", "div", "li", "ul", "ol",
		"h1", "h2", "h3", "h4", "h5", "h6",
		"header", "footer", "table", "tr":
		return true
	default:
		return false
	}
}
