package testing

import (
	"io"
	"testing"
)

func TestStubProcessor(t *testing.T) {
	// Test that StubProcessor implements BytesProcessor
	data := []string{"first", "second", "third"}
	sp := NewStubProcessor(data)

	if sp == nil {
		t.Fatal("expected StubProcessor to be non-nil")
	}

	// Test First Next()
	result, err := sp.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "first" {
		t.Errorf("expected 'first', got '%s'", string(result))
	}

	// Test Second Next()
	result, err = sp.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "second" {
		t.Errorf("expected 'second', got '%s'", string(result))
	}

	// Test Third Next()
	result, err = sp.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}
	if string(result) != "third" {
		t.Errorf("expected 'third', got '%s'", string(result))
	}

	// Test EOF
	result, err = sp.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

func TestNewChunkSlice(t *testing.T) {
	// Test that NewChunkSlice creates ChunkSlice from strings
	input := []string{"first", "second", "third"}
	cs := NewChunkSlice(input...)

	if cs == nil {
		t.Fatal("expected ChunkSlice to be non-nil")
	}

	if len(cs) != 3 {
		t.Fatalf("expected ChunkSlice length 3, got %d", len(cs))
	}

	if string(cs[0]) != "first" {
		t.Errorf("expected 'first', got '%s'", string(cs[0]))
	}

	if string(cs[1]) != "second" {
		t.Errorf("expected 'second', got '%s'", string(cs[1]))
	}

	if string(cs[2]) != "third" {
		t.Errorf("expected 'third', got '%s'", string(cs[2]))
	}
}

func TestNewChunkSliceEmpty(t *testing.T) {
	// Test that NewChunkSlice handles empty input
	cs := NewChunkSlice()

	if cs == nil {
		t.Fatal("expected ChunkSlice to be non-nil")
	}

	if len(cs) != 0 {
		t.Errorf("expected ChunkSlice length 0, got %d", len(cs))
	}
}

func TestStubProcessorEmpty(t *testing.T) {
	// Test that StubProcessor handles empty data
	sp := NewStubProcessor([]string{})

	result, err := sp.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF immediately, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}
