package utf8clean

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"strings"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
)

// Processor reads bytes from a source and validates they are UTF-8.
type Processor struct {
	src              textproc.BytesProcessor
	blockSize        int
	charset          string
	transferEncoding string
	buf              []byte      // Buffer for incomplete UTF-8 sequences
	charsetEncoding  encoding.Encoding // Charset decoder for conversion
}

// Option configures a Processor.
type Option func(*Processor)

// WithBlockSize sets the block size for reading.
// Note: In pull-based mode, block size is determined by the source processor.
func WithBlockSize(n int) Option {
	return func(p *Processor) {
		p.blockSize = n
	}
}

// WithCharset sets the character set for encoding conversion.
func WithCharset(charset string) Option {
	return func(p *Processor) {
		p.charset = charset
	}
}

// getEncoding returns the encoding for a given charset name.
// Returns nil for empty charset or UTF-8 (no conversion needed).
// Uses htmlindex.Get() for standard charset name mapping and alias resolution.
func getEncoding(charset string) (encoding.Encoding, error) {
	// Empty charset means treat as UTF-8 (no conversion)
	if charset == "" {
		return nil, nil
	}

	// UTF-8 is the default, no conversion needed
	if strings.EqualFold(charset, "utf-8") || strings.EqualFold(charset, "utf8") {
		return nil, nil
	}

	// Use htmlindex.Get() for standard charset mapping
	enc, err := htmlindex.Get(charset)
	if err != nil {
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
	return enc, nil
}

// WithTransferEncoding sets the content-transfer-encoding for decoding.
func WithTransferEncoding(encoding string) Option {
	return func(p *Processor) {
		p.transferEncoding = encoding
	}
}

// NewProcessor creates a new Processor with the given BytesProcessor source.
// This enables pull-based lazy evaluation where the processor calls Next() on its source.
// Returns error for invalid charset or transfer encoding.
func NewProcessor(src textproc.BytesProcessor, opts ...Option) (*Processor, error) {
	p := &Processor{
		src:       src,
		blockSize: 1024,
	}
	for _, opt := range opts {
		opt(p)
	}
	// Validate options - get charset encoding using htmlindex.Get()
	var err error
	p.charsetEncoding, err = getEncoding(p.charset)
	if err != nil {
		return nil, ErrInvalidCharset{Charset: p.charset}
	}

	// Validate transfer encoding
	if p.transferEncoding != "" {
		normalized := strings.ToLower(strings.ReplaceAll(p.transferEncoding, " ", ""))
		switch normalized {
		case "base64", "quoted-printable", "qp":
			// Valid
		default:
			return nil, ErrInvalidTransferEncoding{Encoding: p.transferEncoding}
		}
	}

	return p, nil
}

// Next reads the next block of data from the source with validated UTF-8.
// Returns io.EOF when all data has been consumed.
func (p *Processor) Next() ([]byte, error) {
	// Build output from buffered data and new data
	var output []byte

	for len(output) < p.blockSize {
		// If we have buffered data, try to process it first
		if len(p.buf) > 0 {
			valid, remaining, hasIncomplete := p.validateAndBuffer(p.buf)
			output = append(output, valid...)
			p.buf = remaining

			if hasIncomplete {
				// Still incomplete, need more data - continue to read from source
			} else if len(output) > 0 {
				// Have output and no incomplete data, return it
				return output, nil
			}
			// No output and no incomplete data, continue to read from source
		}

		// Read more data from source
		block, err := p.src.Next()
		if err != nil {
			if err == io.EOF {
				// End of input - flush any remaining buffered data
				if len(p.buf) > 0 {
					// Invalid UTF-8 at end - replace with U+FFFD
					output = append(output, '\xef', '\xbf', '\xbd') // U+FFFD in UTF-8
					p.buf = nil
				}
				if len(output) > 0 {
					return output, nil
				}
				return nil, io.EOF
			}
			return nil, err
		}

		// Apply transfer encoding decoding if set
		if p.transferEncoding != "" {
			var err error
			block, err = p.convertTransferEncoding(block)
			if err != nil {
				// Invalid transfer encoding, return error
				return nil, err
			}
		}

		// Apply charset conversion if charset is set
		if p.charsetEncoding != nil {
			block = p.convertCharset(block)
		}

		// Append new data to buffer and process it
		p.buf = append(p.buf, block...)
		valid, remaining, hasIncomplete := p.validateAndBuffer(p.buf)
		output = append(output, valid...)
		p.buf = remaining

		// If we have output and no incomplete data, or we've reached blockSize, return
		if (!hasIncomplete && len(output) > 0) || len(output) >= p.blockSize {
			return output, nil
		}
	}

	if len(output) > 0 {
		return output, nil
	}

	return nil, io.EOF
}

// validateAndBuffer processes a byte slice and returns valid UTF-8 bytes, remaining bytes,
// and whether there's an incomplete sequence at the end.
func (p *Processor) validateAndBuffer(data []byte) (valid, remaining []byte, hasIncomplete bool) {
	for i := 0; i < len(data); {
		r, size := decodeUTF8Rune(data[i:])

		if r == 0xFFFD && size == 1 {
			// Invalid UTF-8 sequence - check if it might be incomplete
			b := data[i]
			if (b&0xE0) == 0xC0 && len(data)-i < 2 {
				// Incomplete 2-byte sequence - buffer it
				return valid, data[i:], true
			}
			if (b&0xF0) == 0xE0 && len(data)-i < 3 {
				// Incomplete 3-byte sequence - buffer it
				return valid, data[i:], true
			}
			if (b&0xF8) == 0xF0 && len(data)-i < 4 {
				// Incomplete 4-byte sequence - buffer it
				return valid, data[i:], true
			}
			// Just an invalid byte - replace with U+FFFD
			valid = append(valid, '\xef', '\xbf', '\xbd')
			i++
			continue
		}

		// Valid UTF-8 sequence
		valid = append(valid, data[i:i+size]...)
		i += size
	}

	return valid, nil, false
}

// convertCharset converts bytes from the source charset to UTF-8.
func (p *Processor) convertCharset(data []byte) []byte {
	if p.charsetEncoding == nil {
		return data
	}

	// Create a transformer for the charset
	decoder := p.charsetEncoding.NewDecoder()

	// Convert the data
	result, err := decoder.Bytes(data)
	if err != nil {
		// On conversion error, replace with U+FFFD
		return []byte{'\xef', '\xbf', '\xbd'}
	}

	return result
}

// convertTransferEncoding decodes Content-Transfer-Encoding.
func (p *Processor) convertTransferEncoding(data []byte) ([]byte, error) {
	if p.transferEncoding == "" {
		return data, nil
	}

	// Normalize encoding name
	normalized := strings.ToLower(strings.ReplaceAll(p.transferEncoding, " ", ""))

	switch normalized {
	case "base64":
		decoder := base64.StdEncoding
		result, err := decoder.DecodeString(string(data))
		if err != nil {
			// On decode error, return as-is (will fail validation)
			return data, nil
		}
		return []byte(result), nil
	case "quoted-printable", "qp":
		// Quoted-printable encoding: =XX where XX is hex value of byte
		decoder := new(mime.WordDecoder)
		decoded := string(data)
		result, err := decoder.DecodeHeader(decoded)
		if err != nil {
			// Fallback: try to decode =XX sequences manually
			decoded = decodeQuotedPrintable(decoded)
		} else if result == decoded {
			// No change, try manual decoding
			decoded = decodeQuotedPrintable(decoded)
		}
		return []byte(decoded), nil
	default:
		// Unknown encoding, return nil data with error
		return nil, ErrInvalidTransferEncoding{Encoding: p.transferEncoding}
	}
}

// decodeQuotedPrintable decodes quoted-printable =XX sequences manually.
func decodeQuotedPrintable(s string) string {
	var result []rune
	for i := 0; i < len(s); i++ {
		if s[i] == '=' && i+2 < len(s) {
			// Try to parse hex value
			hexStr := s[i+1 : i+3]
			if len(hexStr) == 2 {
				var b byte
				_, err := fmt.Sscanf(hexStr, "%02X", &b)
				if err == nil {
					result = append(result, rune(b))
					i += 2
					continue
				}
			}
		}
		result = append(result, rune(s[i]))
	}
	return string(result)
}

// decodeUTF8Rune decodes a UTF-8 rune from a byte slice.
// Returns the rune and the number of bytes consumed.
// Returns U+FFFD and size 1 for invalid sequences.
func decodeUTF8Rune(data []byte) (rune, int) {
	if len(data) == 0 {
		return 0xFFFD, 0
	}

	b := data[0]
	switch {
	case b < 0x80:
		return rune(b), 1
	case b < 0xC0:
		return 0xFFFD, 1 // Invalid continuation byte
	case b < 0xE0:
		if len(data) < 2 {
			return 0xFFFD, 1
		}
		if data[1]&0xC0 != 0x80 {
			return 0xFFFD, 1
		}
		return rune(b&0x1F)<<6 | rune(data[1]&0x3F), 2
	case b < 0xF0:
		if len(data) < 3 {
			return 0xFFFD, 1
		}
		if data[1]&0xC0 != 0x80 || data[2]&0xC0 != 0x80 {
			return 0xFFFD, 1
		}
		return rune(b&0x0F)<<12 | rune(data[1]&0x3F)<<6 | rune(data[2]&0x3F), 3
	case b < 0xF8:
		if len(data) < 4 {
			return 0xFFFD, 1
		}
		if data[1]&0xC0 != 0x80 || data[2]&0xC0 != 0x80 || data[3]&0xC0 != 0x80 {
			return 0xFFFD, 1
		}
		return rune(b&0x07)<<18 | rune(data[1]&0x3F)<<12 | rune(data[2]&0x3F)<<6 | rune(data[3]&0x3F), 4
	default:
		return 0xFFFD, 1
	}
}
