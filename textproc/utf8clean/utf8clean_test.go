package utf8clean

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"strings"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc/reader"
)

// Ensure mime import is used
var _ mime.WordDecoder

// Ensure encoding/base64 import is used
var _ base64.Encoding

func TestNewProcessor(t *testing.T) {
	// Test that NewProcessor creates a processor with default block size
	r := strings.NewReader("test data")
	p, err := NewProcessor(reader.New(r))
	if err != nil {
		t.Fatal(err)
	}

	if p == nil {
		t.Fatal("expected processor to be non-nil")
	}

	// Verify default block size is 1024
	if p.blockSize != 1024 {
		t.Errorf("expected default blockSize 1024, got %d", p.blockSize)
	}
}

func TestNewProcessorWithOptions(t *testing.T) {
	// Test that NewProcessor with options sets custom block size
	r := strings.NewReader("test data")
	p, err := NewProcessor(reader.New(r), WithBlockSize(256))
	if err != nil {
		t.Fatal(err)
	}

	if p == nil {
		t.Fatal("expected processor to be non-nil")
	}

	if p.blockSize != 256 {
		t.Errorf("expected blockSize 256, got %d", p.blockSize)
	}
}

func TestNextSingleBlock(t *testing.T) {
	// Test reading a single block
	data := "hello world"
	r := strings.NewReader(data)
	p, err := NewProcessor(reader.New(r), WithBlockSize(1024))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(result) != data {
		t.Errorf("expected '%s', got '%s'", data, string(result))
	}
}

func TestNextMultipleBlocks(t *testing.T) {
	// Test reading multiple blocks
	data := "hello world, this is a test"
	r := strings.NewReader(data)
	p, err := NewProcessor(reader.New(r, reader.WithBlockSize(10)))
	if err != nil {
		t.Fatal(err)
	}

	// First block
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "hello worl" {
		t.Errorf("expected 'hello worl', got '%s'", string(result))
	}

	// Second block
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "d, this is" {
		t.Errorf("expected 'd, this is', got '%s'", string(result))
	}

	// Third block (partial)
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}
	if string(result) != " a test" {
		t.Errorf("expected ' a test', got '%s'", string(result))
	}

	// Fourth call should return EOF
	result, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

func TestNextEmptyReader(t *testing.T) {
	// Test reading from an empty reader
	r := strings.NewReader("")
	p, err := NewProcessor(reader.New(r))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF immediately, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

// NEW tests for pull-based composition with BytesProcessor

func TestNewProcessorCreatesProcessor(t *testing.T) {
	// Test that NewProcessor creates a processor with BytesProcessor source
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
	_ = p // Use p to avoid linter about unused variable
}

func TestNewProcessorPullsFromSource(t *testing.T) {
	// Test that Next() pulls data from source BytesProcessor
	src := &mockSource{blocks: [][]byte{[]byte("hello"), []byte("world")}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "hello" {
		t.Errorf("expected 'hello', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second call, got %v", err)
	}
	if string(result) != "world" {
		t.Errorf("expected 'world', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

type mockSource struct {
	blocks [][]byte
	index  int
}

func (m *mockSource) Next() ([]byte, error) {
	if m.index >= len(m.blocks) {
		return nil, io.EOF
	}
	result := m.blocks[m.index]
	m.index++
	return result, nil
}

// TDD Step 1: Core Options and Configuration - RED tests

func TestWithCharset(t *testing.T) {
	// Test that WithCharset sets the charset field
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p, _ := NewProcessor(src, WithCharset("ISO-8859-1"))

	if p.charset != "ISO-8859-1" {
		t.Errorf("expected charset 'ISO-8859-1', got '%s'", p.charset)
	}
}

func TestWithTransferEncoding(t *testing.T) {
	// Test that WithTransferEncoding sets the transferEncoding field
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p, _ := NewProcessor(src, WithTransferEncoding("base64"))

	if p.transferEncoding != "base64" {
		t.Errorf("expected transferEncoding 'base64', got '%s'", p.transferEncoding)
	}
}

func TestWithMultipleOptions(t *testing.T) {
	// Test that multiple options can be applied
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p, _ := NewProcessor(src, WithCharset("Windows-1252"), WithTransferEncoding("quoted-printable"))

	if p.charset != "Windows-1252" {
		t.Errorf("expected charset 'Windows-1252', got '%s'", p.charset)
	}
	if p.transferEncoding != "quoted-printable" {
		t.Errorf("expected transferEncoding 'quoted-printable', got '%s'", p.transferEncoding)
	}
}

func TestDefaultCharsetIsEmpty(t *testing.T) {
	// Test that default charset is empty
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	if p.charset != "" {
		t.Errorf("expected empty charset by default, got '%s'", p.charset)
	}
}

func TestDefaultTransferEncodingIsEmpty(t *testing.T) {
	// Test that default transferEncoding is empty
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	if p.transferEncoding != "" {
		t.Errorf("expected empty transferEncoding by default, got '%s'", p.transferEncoding)
	}
}

// TDD Step 2: UTF-8 Validation and Split Sequence Handling - RED tests

func TestValidUTF8Passthrough(t *testing.T) {
	// Test that valid UTF-8 passes through unchanged
	input := "Hello, 世界! 🌍"
	src := &mockSource{blocks: [][]byte{[]byte(input)}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(result) != input {
		t.Errorf("expected '%s', got '%s'", input, string(result))
	}
}

func TestUTF8SequenceSplitAcrossBlocks(t *testing.T) {
	// Test that UTF-8 sequence split across blocks is buffered and reassembled
	// "世界" is "\xe4\xb8\x96\xe7\x95\x8c" in UTF-8
	// Split after first byte: "\xe4" + "\xb8\x96\xe7\x95\x8c"
	input := "世界"
	src := &mockSource{blocks: [][]byte{
		[]byte("\xe4"),           // First byte of first character
		[]byte("\xb8\x96\xe7\x95\x8c"), // Rest of characters
	}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	// Should get reassembled UTF-8
	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != input {
		t.Errorf("expected '%s', got '%s'", input, string(result))
	}

	// Second call should return EOF
	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF on second call, got %v", err)
	}
}

func TestInvalidUTF8ReplacedWithFFFD(t *testing.T) {
	// Test that invalid UTF-8 sequences are replaced with U+FFFD
	// \xff is invalid in UTF-8
	// \x80\x81 is an invalid continuation byte sequence
	src := &mockSource{blocks: [][]byte{
		[]byte("Hello\xffWorld\x80\x81End"),
	}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Convert to string to validate proper UTF-8 with replacement
	resultStr := string(result)
	expectedStr := "Hello\ufffdWorld\ufffd\ufffdEnd"

	// Check rune-by-rune since byte comparison might differ
	resultRunes := []rune(resultStr)
	expectedRunes := []rune(expectedStr)

	if len(resultRunes) != len(expectedRunes) {
		t.Errorf("expected %d runes, got %d", len(expectedRunes), len(resultRunes))
	}

	for i := range expectedRunes {
		if i >= len(resultRunes) {
			t.Errorf("missing rune at position %d: expected %U", i, expectedRunes[i])
		} else if resultRunes[i] != expectedRunes[i] {
			t.Errorf("at position %d: expected %U, got %U", i, expectedRunes[i], resultRunes[i])
		}
	}
}

func TestIncompleteUTF8AtEnd(t *testing.T) {
	// Test that incomplete UTF-8 at end of stream is replaced with U+FFFD
	// \xe4\xb8 is incomplete (missing \x96)
	src := &mockSource{blocks: [][]byte{
		[]byte("Test\xe4\xb8"),
	}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	resultStr := string(result)
	resultRunes := []rune(resultStr)
	// Expected: "Test" (4 runes) + U+FFFD (1 rune) = 5 runes total
	// Our processor recognizes \xe4\xb8 as an incomplete sequence and replaces with ONE U+FFFD
	expectedRunes := []rune("Test\ufffd")

	if len(resultRunes) != len(expectedRunes) {
		t.Errorf("expected %d runes, got %d", len(expectedRunes), len(resultRunes))
	}

	for i := range expectedRunes {
		if i >= len(resultRunes) {
			t.Errorf("missing rune at position %d: expected %U", i, expectedRunes[i])
		} else if resultRunes[i] != expectedRunes[i] {
			t.Errorf("at position %d: expected %U, got %U", i, expectedRunes[i], resultRunes[i])
		}
	}
}

func TestMultipleIncompleteSequences(t *testing.T) {
	// Test multiple incomplete UTF-8 sequences across blocks that get reassembled
	// Split multi-byte character, then another incomplete sequence
	src := &mockSource{blocks: [][]byte{
		[]byte("\xe4\xb8"),      // Incomplete 3-byte sequence (missing \x96)
		[]byte("\x96\xe7"),      // Complete first char, incomplete second (missing \x95\x8c)
		[]byte("\x95\x8c"),      // Complete the sequence
	}}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	resultStr := string(result)
	resultRunes := []rune(resultStr)
	// When properly reassembled, the three blocks form "世界" (two Chinese characters)
	expectedRunes := []rune("世界")

	if len(resultRunes) != len(expectedRunes) {
		t.Errorf("expected %d runes, got %d: expected %U, got %U", len(expectedRunes), len(resultRunes), expectedRunes, resultRunes)
	}

	for i := range expectedRunes {
		if i >= len(resultRunes) {
			t.Errorf("missing rune at position %d: expected %U", i, expectedRunes[i])
		} else if resultRunes[i] != expectedRunes[i] {
			t.Errorf("at position %d: expected %U, got %U", i, expectedRunes[i], resultRunes[i])
		}
	}
}

// TDD Step 3: Charset Conversion - RED tests

func TestISO88591ToUTF8(t *testing.T) {
	// Test ISO-8859-1 to UTF-8 conversion
	// ISO-8859-1 byte 0xE9 is "é" in Latin-1
	input := []byte{0xE9} // é in ISO-8859-1
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithCharset("ISO-8859-1"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// UTF-8 encoding of é is 0xC3 0xA9
	expected := []byte{0xC3, 0xA9}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestWindows1252ToUTF8(t *testing.T) {
	// Test Windows-1252 to UTF-8 conversion
	// Windows-1252 byte 0x80 is the Euro sign
	input := []byte{0x80} // € in Windows-1252
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithCharset("Windows-1252"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// UTF-8 encoding of € is 0xE2 0x82 0xAC
	expected := []byte{0xE2, 0x82, 0xAC}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestUTF16LEToUTF8(t *testing.T) {
	// Test UTF-16LE to UTF-8 conversion
	// UTF-16LE encoding of "A" is 0x41 0x00
	input := []byte{0x41, 0x00} // "A" in UTF-16LE
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithCharset("UTF-16LE"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := []byte{'A'}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestUTF16BEToUTF8(t *testing.T) {
	// Test UTF-16BE to UTF-8 conversion
	// UTF-16BE encoding of "A" is 0x00 0x41
	input := []byte{0x00, 0x41} // "A" in UTF-16BE
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithCharset("UTF-16BE"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := []byte{'A'}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestShiftJISToUTF8(t *testing.T) {
	// Test Shift_JIS to UTF-8 conversion
	// Shift_JIS encoding of "あ" (hiragana 'a') is 0x82 0xA0
	input := []byte{0x82, 0xA0} // "あ" in Shift_JIS
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithCharset("Shift_JIS"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// UTF-8 encoding of "あ" is 0xE3 0x81 0x82
	expected := []byte{0xE3, 0x81, 0x82}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestNoCharsetPassthrough(t *testing.T) {
	// Test that without charset option, data passes through as UTF-8
	input := []byte("Hello")
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src) // No charset option
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if !bytes.Equal(result, input) {
		t.Errorf("expected %x, got %x", input, result)
	}
}

// TDD Step 4: Content-Transfer-Encoding Decoding - RED tests

func TestBase64Decoding(t *testing.T) {
	// Test base64 decoding
	// "Hello" in base64 is "SGVsbG8="
	input := []byte("SGVsbG8=")
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithTransferEncoding("base64"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := []byte("Hello")
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestQuotedPrintableDecoding(t *testing.T) {
	// Test quoted-printable decoding
	// "Hello World" in quoted-printable is "Hello=20World"
	input := []byte("Hello=20World")
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithTransferEncoding("quoted-printable"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := []byte("Hello World")
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestCharsetWithTransferEncoding(t *testing.T) {
	// Test combined charset + transfer encoding
	// "é" (U+00E9) in ISO-8859-1 is 0xE9, base64 encoded is "6Q=="
	input := []byte("6Q==")
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src, WithCharset("ISO-8859-1"), WithTransferEncoding("base64"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// After base64 decode: 0xE9 (ISO-8859-1 é)
	// After charset conversion: 0xC3 0xA9 (UTF-8 é)
	expected := []byte{0xC3, 0xA9}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestNoTransferEncodingPassthrough(t *testing.T) {
	// Test that without transfer encoding option, data passes through
	input := []byte("Hello")
	src := &mockSource{blocks: [][]byte{input}}
	p, err := NewProcessor(src) // No transfer encoding option
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if !bytes.Equal(result, input) {
		t.Errorf("expected %x, got %x", input, result)
	}
}

// TDD Step 5: Charset Detection - RED tests

func TestCharsetDetectionHTMLIndexLookup(t *testing.T) {
	// Test charset detection using htmlindex.Get
	// htmlindex.Get maps charset names to encodings
	src := &mockSource{blocks: [][]byte{[]byte("Hello")}}
	p, err := NewProcessor(src, WithCharset("")) // Empty charset means auto-detect
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// Default should be UTF-8 (auto-detect defaults to UTF-8)
	expected := []byte("Hello")
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

func TestCharsetDetectionFallsBackToUTF8(t *testing.T) {
	// Test that when charset is empty, data is treated as UTF-8
	input := "Hello 世界" // Mix of ASCII and Chinese
	src := &mockSource{blocks: [][]byte{[]byte(input)}}
	p, err := NewProcessor(src) // No charset option
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := []byte(input)
	if !bytes.Equal(result, expected) {
		t.Errorf("expected %x, got %x", expected, result)
	}
}

// TDD Step 6: Error Cases - RED tests

func TestInvalidCharsetError(t *testing.T) {
	// Test that invalid charset returns error from NewProcessor
	src := &mockSource{blocks: [][]byte{[]byte("Hello")}}
	_, err := NewProcessor(src, WithCharset("invalid-charset-name"))
	_ = err // Use err to avoid linter

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should be ErrInvalidCharset
	if _, ok := err.(ErrInvalidCharset); !ok {
		t.Errorf("expected ErrInvalidCharset, got %T", err)
	}
}

func TestInvalidTransferEncodingError(t *testing.T) {
	// Test that invalid transfer encoding returns error from NewProcessor
	src := &mockSource{blocks: [][]byte{[]byte("Hello")}}
	_, err := NewProcessor(src, WithTransferEncoding("invalid-transfer-encoding"))
	_ = err // Use err to avoid linter

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should be ErrInvalidTransferEncoding
	if _, ok := err.(ErrInvalidTransferEncoding); !ok {
		t.Errorf("expected ErrInvalidTransferEncoding, got %T", err)
	}
}

func TestSourceErrorPropagation(t *testing.T) {
	// Test that errors from source are propagated
	expectedErr := fmt.Errorf("source error")
	src := &errorSource{err: expectedErr}
	p, err := NewProcessor(src)
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Next()
	if err != expectedErr {
		t.Fatalf("expected source error %v, got %v", expectedErr, err)
	}
}

type errorSource struct {
	err error
}

func (e *errorSource) Next() ([]byte, error) {
	return nil, e.err
}

