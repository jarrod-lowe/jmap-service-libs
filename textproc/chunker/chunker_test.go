package chunker

import (
	"errors"
	"io"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// mockBytesProcessor implements textproc.BytesProcessor for testing
type mockBytesProcessor struct {
	blocks [][]byte
	index  int
	err    error
}

func (m *mockBytesProcessor) Next() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.index >= len(m.blocks) {
		return nil, io.EOF
	}
	result := m.blocks[m.index]
	m.index++
	return result, nil
}

func TestBasicParagraphSplitting(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("Para1\n\nPara2\n\nPara3")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "Para2" {
		t.Errorf("expected 'Para2', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}
	if string(result) != "Para3" {
		t.Errorf("expected 'Para3', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestCrossPlatformEndings(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("Para1\r\n\r\nPara2")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "Para2" {
		t.Errorf("expected 'Para2', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestHorizontalRules(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("Para1\n---\nPara2\n***\nPara3")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "Para2" {
		t.Errorf("expected 'Para2', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}
	if string(result) != "Para3" {
		t.Errorf("expected 'Para3', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestWhitespaceTrimming(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("  Para1  \n\n  Para2  ")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "Para2" {
		t.Errorf("expected 'Para2', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestEmptyParagraphSkipping(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("Para1\n\n\n\nPara2")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "Para2" {
		t.Errorf("expected 'Para2', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestBoundarySplitAcrossBlocks(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("Para1\n"), []byte("\nPara2")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "Para2" {
		t.Errorf("expected 'Para2', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestInterfaceCompliance(t *testing.T) {
	var _ textproc.ChunkProcessor = (*Processor)(nil)
}

func TestEOFHandling(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("Para1")}}
	p := NewProcessor(src)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "Para1" {
		t.Errorf("expected 'Para1', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestErrorPropagation(t *testing.T) {
	testErr := errors.New("test error")
	src := &mockBytesProcessor{err: testErr}
	p := NewProcessor(src)

	_, err := p.Next()
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}

func TestNewProcessorCreatesProcessor(t *testing.T) {
	src := &mockBytesProcessor{blocks: [][]byte{[]byte("test")}}
	p := NewProcessor(src)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}
