package htmlstrip

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Processor reads bytes from an io.Reader and returns them in blocks
// with HTML markup removed.
type Processor struct {
	r         io.Reader
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

// New creates a new Processor with the given reader and options.
// The default block size is 1024 bytes.
func New(r io.Reader, opts ...Option) *Processor {
	p := &Processor{
		r:         r,
		blockSize: 1024,
		tokenizer: html.NewTokenizer(r),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next reads the next block of data from the reader with HTML removed.
// Returns io.EOF when all data has been consumed.
func (p *Processor) Next() ([]byte, error) {
	if p.done && p.buf.Len() == 0 {
		return nil, io.EOF
	}

	for {
		if p.buf.Len() >= p.blockSize {
			result := strings.TrimRight(p.buf.String(), " \t\n\r")
			p.buf.Reset()
			return []byte(result), nil
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
					return []byte(result), nil
				}
				return nil, io.EOF
			}
			return nil, err

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
