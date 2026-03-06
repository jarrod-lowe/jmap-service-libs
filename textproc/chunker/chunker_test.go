package chunker

import (
	"io"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

func TestNextReturnsChunks(t *testing.T) {
	// Test that Next splits data into chunks at fixed boundaries
	data := [][]byte{[]byte("hello world")}
	p := New(data, WithChunkSize(5))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}

	if string(result) != "hello" {
		t.Errorf("expected 'hello', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}

	if string(result) != " worl" {
		t.Errorf("expected ' worl', got '%s'", string(result))
	}

	// Third call returns remaining bytes
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}

	if string(result) != "d" {
		t.Errorf("expected 'd', got '%s'", string(result))
	}

	// Fourth call should return EOF
	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestNextImplementsChunkProcessor(t *testing.T) {
	// Verify Processor satisfies textproc.ChunkProcessor interface
	data := [][]byte{[]byte("test")}
	p := New(data)

	var _ textproc.ChunkProcessor = p
}
