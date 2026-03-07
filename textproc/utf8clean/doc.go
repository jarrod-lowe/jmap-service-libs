// Package utf8clean ensures UTF-8 encoding in input data.
//
// The Processor validates UTF-8 and can convert from various charsets and transfer encodings.
// It handles multi-byte sequences split across block boundaries and replaces invalid
// UTF-8 sequences with U+FFFD (replacement character).
//
// Charset support is provided by golang.org/x/text/encoding/htmlindex, which maps
// HTML charset names (as found in email Content-Type headers) to encoding.Encodings.
// Commonly supported charsets include:
//   - ISO-8859-* (Western European, Cyrillic, Arabic, Hebrew, etc.)
//   - Windows-* codepages (1250-1258)
//   - UTF-16, UTF-32 (with BOM detection)
//   - Asian encodings (Shift_JIS, EUC-JP, GBK, Big5, etc.)
//
// For a complete list of supported charsets, see:
// https://pkg.go.dev/golang.org/x/text/encoding/htmlindex
//
// Supported transfer encodings:
//   - base64
//   - quoted-printable (QP)
//
// When charset is not specified or empty, data is treated as UTF-8.
// When charset is specified, data is converted from that charset to UTF-8.
//
// Errors:
//   - ErrInvalidCharset: returned when an unsupported charset is specified
//   - ErrInvalidTransferEncoding: returned when an unsupported transfer encoding is specified
//
// Example usage:
//
//	// Simple UTF-8 validation (data must be valid UTF-8)
//	p, err := utf8clean.NewProcessor(src)
//	if err != nil {
//	    return err
//	}
//	data, err := p.Next()
//
//	// Convert ISO-8859-1 to UTF-8
//	p, err := utf8clean.NewProcessor(src, utf8clean.WithCharset("ISO-8859-1"))
//	if err != nil {
//	    return err
//	}
//	data, err := p.Next()
//
//	// Decode base64
//	p, err := utf8clean.NewProcessor(src, utf8clean.WithTransferEncoding("base64"))
//	if err != nil {
//	    return err
//	}
//	data, err := p.Next()
//
//	// Decode quoted-printable
//	p, err := utf8clean.NewProcessor(src, utf8clean.WithTransferEncoding("quoted-printable"))
//	if err != nil {
//	    return err
//	}
//	data, err := p.Next()
//
//	// Combine charset and transfer encoding
//	p, err := utf8clean.NewProcessor(src,
//	    utf8clean.WithCharset("ISO-8859-1"),
//	    utf8clean.WithTransferEncoding("base64"))
//	if err != nil {
//	    return err
//	}
//	data, err := p.Next()
package utf8clean