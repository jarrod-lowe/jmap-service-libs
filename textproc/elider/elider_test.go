package elider

import (
	"io"
	"strings"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc/reader"
)

func TestNewProcessor(t *testing.T) {
	// Test that NewProcessor creates a processor with default block size
	r := strings.NewReader("test data")
	p := NewProcessor(reader.New(r))

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
	p := NewProcessor(reader.New(r), WithBlockSize(256))

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
	p := NewProcessor(reader.New(r), WithBlockSize(1024))

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
	p := NewProcessor(reader.New(r, reader.WithBlockSize(10)))

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
	p := NewProcessor(reader.New(r))

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
	src := &mockSource{blocks: [][]byte{[]byte("test")}}
	p := NewProcessor(src)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}

func TestNewProcessorPullsFromSource(t *testing.T) {
	src := &mockSource{blocks: [][]byte{[]byte("hello"), []byte("world")}}
	p := NewProcessor(src)

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
