package textproc

import (
	"testing"
)

func TestChunkType(t *testing.T) {
	// Chunk is a type alias for []byte
	// This test verifies Chunk can be created and used
	c := Chunk([]byte("test data"))

	if len(c) != 9 {
		t.Errorf("expected Chunk length 9, got %d", len(c))
	}

	if string(c) != "test data" {
		t.Errorf("expected Chunk content 'test data', got '%s'", string(c))
	}
}

func TestChunkSliceType(t *testing.T) {
	// ChunkSlice is a slice of Chunks
	// This test verifies ChunkSlice can be created and used
	cs := ChunkSlice{
		Chunk("first"),
		Chunk("second"),
		Chunk("third"),
	}

	if len(cs) != 3 {
		t.Errorf("expected ChunkSlice length 3, got %d", len(cs))
	}

	if string(cs[0]) != "first" {
		t.Errorf("expected first chunk 'first', got '%s'", string(cs[0]))
	}

	if string(cs[1]) != "second" {
		t.Errorf("expected second chunk 'second', got '%s'", string(cs[1]))
	}

	if string(cs[2]) != "third" {
		t.Errorf("expected third chunk 'third', got '%s'", string(cs[2]))
	}
}

func TestChunkConversion(t *testing.T) {
	// Test that Chunk is a string type
	c := Chunk("hello world")

	if c == "" {
		t.Error("expected Chunk to be non-empty")
	}

	if len(c) != 11 {
		t.Errorf("expected Chunk length 11, got %d", len(c))
	}

	if c != "hello world" {
		t.Errorf("expected Chunk content 'hello world', got '%s'", c)
	}
}

func TestChunkAsString(t *testing.T) {
	// Chunk is a type alias for string
	// This test verifies Chunk can be created directly from string
	c := Chunk("test data")

	if len(c) != 9 {
		t.Errorf("expected Chunk length 9, got %d", len(c))
	}

	// As a string, no conversion needed
	if c != "test data" {
		t.Errorf("expected Chunk content 'test data', got '%s'", c)
	}
}

func TestChunkSliceAsString(t *testing.T) {
	// ChunkSlice is a slice of string Chunks
	cs := ChunkSlice{
		"first",
		"second",
		"third",
	}

	if len(cs) != 3 {
		t.Errorf("expected ChunkSlice length 3, got %d", len(cs))
	}

	// No string conversion needed
	if cs[0] != "first" {
		t.Errorf("expected first chunk 'first', got '%s'", cs[0])
	}

	if cs[1] != "second" {
		t.Errorf("expected second chunk 'second', got '%s'", cs[1])
	}

	if cs[2] != "third" {
		t.Errorf("expected third chunk 'third', got '%s'", cs[2])
	}
}
