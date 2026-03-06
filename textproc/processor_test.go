package textproc

import (
	"io"
	"testing"
)

// testBytesProcessor implements BytesProcessor for testing
type testBytesProcessor struct {
	data  [][]byte
	index int
}

func newTestBytesProcessor(data [][]byte) *testBytesProcessor {
	return &testBytesProcessor{data: data, index: 0}
}

func (t *testBytesProcessor) Next() ([]byte, error) {
	if t.index >= len(t.data) {
		return nil, io.EOF
	}
	result := t.data[t.index]
	t.index++
	return result, nil
}

func TestBytesProcessorInterface(t *testing.T) {
	// Test that a type implementing Next() ([]byte, error) satisfies BytesProcessor
	data := [][]byte{
		[]byte("first block"),
		[]byte("second block"),
	}

	var p BytesProcessor = newTestBytesProcessor(data)

	// First Next() call
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "first block" {
		t.Errorf("expected 'first block', got '%s'", string(result))
	}

	// Second Next() call
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "second block" {
		t.Errorf("expected 'second block', got '%s'", string(result))
	}

	// Third call should return EOF
	result, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

// testChunkProcessor implements ChunkProcessor for testing
type testChunkProcessor struct {
	chunks []Chunk
	index  int
}

func newTestChunkProcessor(chunks []Chunk) *testChunkProcessor {
	return &testChunkProcessor{chunks: chunks, index: 0}
}

func (t *testChunkProcessor) Next() (Chunk, error) {
	if t.index >= len(t.chunks) {
		return nil, io.EOF
	}
	result := t.chunks[t.index]
	t.index++
	return result, nil
}

func TestChunkProcessorInterface(t *testing.T) {
	// Test that a type implementing Next() (Chunk, error) satisfies ChunkProcessor
	chunks := []Chunk{
		Chunk("first chunk"),
		Chunk("second chunk"),
	}

	var p ChunkProcessor = newTestChunkProcessor(chunks)

	// First Next() call
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "first chunk" {
		t.Errorf("expected 'first chunk', got '%s'", string(result))
	}

	// Second Next() call
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "second chunk" {
		t.Errorf("expected 'second chunk', got '%s'", string(result))
	}

	// Third call should return EOF
	result, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

// testChunkCombiner implements ChunkCombiner for testing
type testChunkCombiner struct {
	chunkSlices []ChunkSlice
	index       int
}

func newTestChunkCombiner(slices []ChunkSlice) *testChunkCombiner {
	return &testChunkCombiner{chunkSlices: slices, index: 0}
}

func (t *testChunkCombiner) Next() (ChunkSlice, error) {
	if t.index >= len(t.chunkSlices) {
		return nil, io.EOF
	}
	result := t.chunkSlices[t.index]
	t.index++
	return result, nil
}

func TestChunkCombinerInterface(t *testing.T) {
	// Test that a type implementing Next() (ChunkSlice, error) satisfies ChunkCombiner
	slices := []ChunkSlice{
		ChunkSlice{Chunk("a"), Chunk("b")},
		ChunkSlice{Chunk("c"), Chunk("d"), Chunk("e")},
	}

	var c ChunkCombiner = newTestChunkCombiner(slices)

	// First Next() call
	result, err := c.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 chunks in first slice, got %d", len(result))
	}
	if string(result[0]) != "a" {
		t.Errorf("expected 'a', got '%s'", string(result[0]))
	}

	// Second Next() call
	result, err = c.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 chunks in second slice, got %d", len(result))
	}

	// Third call should return EOF
	result, err = c.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

func TestProcessorIntegration(t *testing.T) {
	// Test that we can compose processors through the interfaces
	// This verifies the interfaces work together as expected

	// BytesProcessor produces []byte
	byteData := [][]byte{[]byte("hello"), []byte("world")}
	bp := newTestBytesProcessor(byteData)

	// Convert []byte to Chunk for ChunkProcessor
	chunks := make([]Chunk, len(byteData))
	for i, d := range byteData {
		chunks[i] = Chunk(d)
	}
	cp := newTestChunkProcessor(chunks)

	// Convert chunks to ChunkSlice for ChunkCombiner
	slices := []ChunkSlice{{chunks[0]}, {chunks[1]}}
	cc := newTestChunkCombiner(slices)

	// Verify we can iterate through all
	var bpInterface BytesProcessor = bp
	b1, _ := bpInterface.Next()
	if string(b1) != "hello" {
		t.Errorf("expected 'hello' from BytesProcessor, got '%s'", string(b1))
	}

	var cpInterface ChunkProcessor = cp
	c1, _ := cpInterface.Next()
	if string(c1) != "hello" {
		t.Errorf("expected 'hello' from ChunkProcessor, got '%s'", string(c1))
	}

	var ccInterface ChunkCombiner = cc
	s1, _ := ccInterface.Next()
	if len(s1) != 1 {
		t.Errorf("expected 1 chunk in ChunkSlice, got %d", len(s1))
	}
}
