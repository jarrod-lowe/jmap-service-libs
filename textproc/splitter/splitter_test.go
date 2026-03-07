package splitter

import (
	"errors"
	"io"
	"testing"
	"unicode/utf8"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// mockChunkProcessor implements textproc.ChunkProcessor for testing
type mockChunkProcessor struct {
	chunks []textproc.Chunk
	index  int
	err    error
}

func (m *mockChunkProcessor) Next() (textproc.Chunk, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.index >= len(m.chunks) {
		return nil, io.EOF
	}
	result := m.chunks[m.index]
	m.index++
	return result, nil
}

// TestPassthrough verifies chunks under the byte limit pass through unchanged.
func TestPassthrough(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("small chunk"),
		textproc.Chunk("another small"),
	}}
	p := NewProcessor(src, 1000)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "small chunk" {
		t.Errorf("expected 'small chunk', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "another small" {
		t.Errorf("expected 'another small', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestSentenceSplit verifies large chunks split at sentence boundaries.
func TestSentenceSplit(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("First sentence. Second sentence."),
	}}
	p := NewProcessor(src, 20)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "First sentence." {
		t.Errorf("expected 'First sentence.', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "Second sentence." {
		t.Errorf("expected 'Second sentence.', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestWordSplit verifies chunks split at word boundaries when no sentence markers exist.
func TestWordSplit(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("one two three four"),
	}}
	p := NewProcessor(src, 10)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "one two" {
		t.Errorf("expected 'one two', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "three four" {
		t.Errorf("expected 'three four', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestCharacterSplit verifies chunks split at character boundaries when no boundaries found.
func TestCharacterSplit(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("longwordwithoutspaces"),
	}}
	p := NewProcessor(src, 10)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "longwordwi" {
		t.Errorf("expected 'longwordwi', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "thoutspaces" {
		t.Errorf("expected 'thoutspaces', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestOversizedContent verifies that character splitting happens when no sentence/word boundaries exist.
func TestOversizedContent(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("verylongwordexceedinglimit"),
	}}
	p := NewProcessor(src, 5)

	// Should split into multiple chunks
	count := 0
	for {
		result, err := p.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) > 10 { // May get chunks up to 2x limit
			t.Errorf("expected at most 10 bytes, got %d: '%s'", len(result), string(result))
		}
		count++
	}

	// Should get at least 2 chunks
	if count < 2 {
		t.Errorf("expected at least 2 chunks, got %d", count)
	}
}

// TestMultipleSplits verifies large chunks split into multiple pieces.
func TestMultipleSplits(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("First. Second. Third. Fourth."),
	}}
	p := NewProcessor(src, 10)

	results := make([]string, 0)
	for {
		result, err := p.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		results = append(results, string(result))
	}

	// Should split into at least 2 chunks
	if len(results) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(results))
	}
}

// TestCloseToLimit verifies chunk slightly over limit splits into two pieces.
func TestCloseToLimit(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("This is slightly over the limit"),
	}}
	p := NewProcessor(src, 15)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "This is" {
		t.Errorf("expected 'This is', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Remaining chunk should be returned as-is (24 bytes < 30 = 2*15)
	if string(result) != "slightly over the limit" {
		t.Errorf("expected 'slightly over the limit', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestEmptyChunks verifies empty chunks are skipped.
func TestEmptyChunks(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk(""),
		textproc.Chunk("valid"),
		textproc.Chunk(""),
	}}
	p := NewProcessor(src, 1000)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "valid" {
		t.Errorf("expected 'valid', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestErrorPropagation verifies source errors are passed through unchanged.
func TestErrorPropagation(t *testing.T) {
	testErr := errors.New("source error")
	src := &mockChunkProcessor{err: testErr}
	p := NewProcessor(src, 100)

	_, err := p.Next()
	if err != testErr {
		t.Errorf("expected source error, got %v", err)
	}
}

// TestUTF8Safety verifies multi-byte characters are not broken during character splits.
func TestUTF8Safety(t *testing.T) {
	// Use shorter string to fit within limit
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("日本語。英語。"),
	}}
	p := NewProcessor(src, 10)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should split at character boundary since no space after period
	if string(result) != "日本語" {
		t.Errorf("expected '日本語', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Remaining chunk should be returned as-is (7 bytes < 30 = 3*10)
	if string(result) != "。英語。" {
		t.Errorf("expected '。英語。', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestCharacterSplitUTF8NoBreak verifies character split doesn't break multi-byte chars.
func TestCharacterSplitUTF8NoBreak(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{
		textproc.Chunk("日本語文字"),
	}}
	p := NewProcessor(src, 10)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should split at safe boundary - "日本語" is 9 bytes
	if string(result) != "日本語" {
		t.Errorf("expected '日本語', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "文字" {
		t.Errorf("expected '文字', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestNewProcessorCreatesProcessor verifies NewProcessor creates a valid Processor.
func TestNewProcessorCreatesProcessor(t *testing.T) {
	src := &mockChunkProcessor{chunks: []textproc.Chunk{[]byte("test")}}
	p := NewProcessor(src, 100)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}

// TestInterfaceCompliance verifies Processor implements ChunkProcessor.
func TestInterfaceCompliance(t *testing.T) {
	var _ textproc.ChunkProcessor = (*Processor)(nil)
}

// TestUTF8Validation helper to verify a string is valid UTF-8
func TestUTF8Validation(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"valid ASCII", "hello", true},
		{"valid UTF-8", "日本語", true},
		{"mixed", "hello世界", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utf8.ValidString(tt.s)
			if got != tt.want {
				t.Errorf("ValidString(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
