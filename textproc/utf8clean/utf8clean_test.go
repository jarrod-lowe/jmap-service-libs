package utf8clean

import (
	"io"
	"strings"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	// Test that New creates a processor with default block size
	r := strings.NewReader("test data")
	p := New(r)

	if p == nil {
		t.Fatal("expected processor to be non-nil")
	}

	// Verify default block size is 1024
	if p.blockSize != 1024 {
		t.Errorf("expected default blockSize 1024, got %d", p.blockSize)
	}
}

func TestNewProcessorWithOptions(t *testing.T) {
	// Test that New with options sets custom block size
	r := strings.NewReader("test data")
	p := New(r, WithBlockSize(256))

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
	p := New(r, WithBlockSize(1024))

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
	p := New(r, WithBlockSize(10))

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
	p := New(r)

	result, err := p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF immediately, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}
