package chain

import (
	"io"
	"strings"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/reader"
)

func TestBytesToChunkAdaptsBytesProcessorToChunkProcessor(t *testing.T) {
	// Test that bytesToChunk adapts BytesProcessor to ChunkProcessor
	r := strings.NewReader("test data")
	src := reader.New(r)

	btc := &bytesToChunk{src: src}

	chunk, err := btc.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(chunk) != "test data" {
		t.Errorf("expected 'test data', got '%s'", string(chunk))
	}

	// Second call should return EOF
	_, err = btc.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF on second call, got %v", err)
	}
}

func TestBytesToChunkEmptyData(t *testing.T) {
	src := &mockBytesProc{blocks: [][]byte{}}
	btc := &bytesToChunk{src: src}

	_, err := btc.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF for empty source, got %v", err)
	}
}

func TestBytesToChunkMultipleBlocks(t *testing.T) {
	src := &mockBytesProc{blocks: [][]byte{[]byte("first"), []byte("second")}}
	btc := &bytesToChunk{src: src}

	chunk1, err := btc.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(chunk1) != "first" {
		t.Errorf("expected 'first', got '%s'", string(chunk1))
	}

	chunk2, err := btc.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(chunk2) != "second" {
		t.Errorf("expected 'second', got '%s'", string(chunk2))
	}

	_, err = btc.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// Ensure bytesToChunk implements textproc.ChunkProcessor
var _ textproc.ChunkProcessor = (*bytesToChunk)(nil)

// mockBytesProc is a mock BytesProcessor for testing
type mockBytesProc struct {
	blocks [][]byte
	index  int
}

func (m *mockBytesProc) Next() ([]byte, error) {
	if m.index >= len(m.blocks) {
		return nil, io.EOF
	}
	result := m.blocks[m.index]
	m.index++
	return result, nil
}
