package elider

import (
	"io"
	"strings"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/reader"
)

func TestNewProcessor(t *testing.T) {
	// Test that NewProcessor creates a processor with default block size
	r := strings.NewReader("test data")
	p := NewProcessor(textproc.NewBytesToStringAdapter(reader.New(r)))

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
	p := NewProcessor(textproc.NewBytesToStringAdapter(reader.New(r)), WithBlockSize(256))

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
	p := NewProcessor(textproc.NewBytesToStringAdapter(reader.New(r)), WithBlockSize(1024))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != data {
		t.Errorf("expected '%s', got '%s'", data, result)
	}
}

func TestNextMultipleBlocks(t *testing.T) {
	// Test reading multiple blocks from source
	data := "hello world, this is a test"
	r := strings.NewReader(data)
	p := NewProcessor(textproc.NewBytesToStringAdapter(reader.New(r)), WithBlockSize(1024))

	// First block - elided content (no URLs, quotes, signatures to elide)
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	// Normal text is preserved as-is when no elision rules apply
	if result != "hello world, this is a test" {
		t.Errorf("expected 'hello world, this is a test', got '%s'", result)
	}

	// Second call should return EOF
	result, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result with EOF, got '%s'", result)
	}
}

func TestNextEmptyReader(t *testing.T) {
	// Test reading from an empty reader
	r := strings.NewReader("")
	p := NewProcessor(textproc.NewBytesToStringAdapter(reader.New(r)))

	result, err := p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF immediately, got %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result with EOF, got '%s'", result)
	}
}

// NEW tests for pull-based composition with StringProcessor

type mockStringSource struct {
	blocks []string
	index  int
}

func (m *mockStringSource) Next() (string, error) {
	if m.index >= len(m.blocks) {
		return "", io.EOF
	}
	result := m.blocks[m.index]
	m.index++
	return result, nil
}

func TestNewProcessorCreatesProcessor(t *testing.T) {
	src := &mockStringSource{blocks: []string{"test"}}
	p := NewProcessor(src)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}

func TestNewProcessorPullsFromSource(t *testing.T) {
	src := &mockStringSource{blocks: []string{"hello", "world"}}
	p := NewProcessor(src)

	// With elider, normal text is preserved across blocks
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "helloworld" {
		t.Errorf("expected 'helloworld', got '%s'", result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: URL elision - keep domain, strip query params and fragments
func TestElideURLKeepsDomainStripsQueryAndFragment(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Check https://api.github.com/repos/1234?token=abc#section"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should keep normal text and domain, strip query parameters and fragments
	expected := "Check api.github.com"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: http:// protocol (not just https://)
func TestElideHTTPProtocol(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Visit http://example.com/path?id=123"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "Visit example.com"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: Quote elision - elide everything after quote marker
func TestElideQuoteMarkerRFC5322(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Important text\nOn Jan 1, 2024, John Doe wrote:\nThis is quoted\nEnd of message"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should keep text before quote, elide everything after
	// Note: trailing whitespace is trimmed
	expected := "Important text"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: Signature elision - elide everything after -- delimiter
func TestElideSignatureDelimiter(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Message body\n--\nJohn Doe\njohn@example.com"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should keep text before signature delimiter, elide everything after
	expected := "Message body"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: Non-word filtering - remove UUIDs
func TestElideUUIDs(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Your ID is 550e8400-e29b-41d4-a716-446655440000"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should keep surrounding text, remove UUID
	expected := "Your ID is"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: Non-word filtering - remove hex strings (16+ consecutive hex chars)
func TestElideHexStrings(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Key: 0123456789abcdef0123456789abcdef"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should keep surrounding text, remove hex string
	expected := "Key:"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: Non-word filtering - remove version strings (v1, v2, v1.2.3, etc.)
func TestElideVersionStrings(t *testing.T) {
	src := &mockStringSource{blocks: []string{"Running version v1.2.3 of the app"}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should keep surrounding text, remove version string
	expected := "Running version of the app"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// RED test: Block boundary handling - URL pattern split across blocks
func TestBlockBoundaryURL(t *testing.T) {
	// URL pattern split across blocks
	src := &mockStringSource{blocks: []string{
		"Check https://api.",
		"github.com/",
		"path?id=123",
	}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "Check api.github.com"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}
