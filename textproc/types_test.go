package textproc

import (
	"testing"
)

func TestChunkType(t *testing.T) {
	// Chunk is a type alias for []byte
	// This test verifies Chunk can be created and used
	var c Chunk
	c = []byte("test data")

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
	var cs ChunkSlice
	cs = ChunkSlice{
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
	// Test conversion between []byte and Chunk
	data := []byte("hello world")
	c := Chunk(data)

	if c == nil {
		t.Error("expected Chunk to be non-nil")
	}

	if len(c) != len(data) {
		t.Errorf("expected Chunk length %d, got %d", len(data), len(c))
	}

	// Verify it's a true conversion
	for i := range data {
		if c[i] != data[i] {
			t.Errorf("mismatch at index %d: expected %d, got %d", i, data[i], c[i])
		}
	}
}
