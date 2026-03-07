package reader

import (
	"io"
	"strings"
	"testing"
)

func TestNewCreatesProcessor(t *testing.T) {
	r := strings.NewReader("test data")
	p := New(r)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}

func TestNextReturnsDataFromReader(t *testing.T) {
	r := strings.NewReader("test data")
	p := New(r)

	data, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != "test data" {
		t.Errorf("expected 'test data', got '%s'", string(data))
	}
}

func TestNextReturnsEOFWhenExhausted(t *testing.T) {
	r := strings.NewReader("test")
	p := New(r)

	// First read should succeed
	_, err := p.Next()
	if err != nil {
		t.Fatalf("first Next() should succeed, got %v", err)
	}

	// Second read should return EOF
	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF on second read, got %v", err)
	}
}

func TestNextRespectsBlockSize(t *testing.T) {
	r := strings.NewReader("0123456789") // 10 bytes
	p := New(r, WithBlockSize(3))

	data, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) != 3 {
		t.Errorf("expected block size 3, got %d", len(data))
	}
	if string(data) != "012" {
		t.Errorf("expected '012', got '%s'", string(data))
	}
}

func TestNextMultipleReads(t *testing.T) {
	r := strings.NewReader("0123456789") // 10 bytes
	p := New(r, WithBlockSize(3))

	data1, err := p.Next()
	if err != nil {
		t.Fatalf("first Next() failed: %v", err)
	}
	if string(data1) != "012" {
		t.Errorf("expected '012', got '%s'", string(data1))
	}

	data2, err := p.Next()
	if err != nil {
		t.Fatalf("second Next() failed: %v", err)
	}
	if string(data2) != "345" {
		t.Errorf("expected '345', got '%s'", string(data2))
	}

	data3, err := p.Next()
	if err != nil {
		t.Fatalf("third Next() failed: %v", err)
	}
	if string(data3) != "678" {
		t.Errorf("expected '678', got '%s'", string(data3))
	}

	data4, err := p.Next()
	if err != nil {
		t.Fatalf("fourth Next() failed: %v", err)
	}
	if string(data4) != "9" {
		t.Errorf("expected '9', got '%s'", string(data4))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestEmptyReader(t *testing.T) {
	r := strings.NewReader("")
	p := New(r)

	_, err := p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF for empty reader, got %v", err)
	}
}
