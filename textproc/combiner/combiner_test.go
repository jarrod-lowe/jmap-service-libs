package combiner

import (
	"io"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// mockChunkSource is a simple ChunkProcessor for testing
type mockChunkSource struct {
	chunks []textproc.Chunk
	index  int
}

func (m *mockChunkSource) Next() (textproc.Chunk, error) {
	if m.index >= len(m.chunks) {
		return "", io.EOF
	}
	result := m.chunks[m.index]
	m.index++
	return result, nil
}

// totalChars returns the total character count of a ChunkSlice
func totalChars(cs textproc.ChunkSlice) int {
	var total int
	for _, c := range cs {
		total += len([]rune(c))
	}
	return total
}

// TestBasicAccumulation verifies multiple small chunks are combined into single ChunkSlice
func TestBasicAccumulation(t *testing.T) {
	// Three chunks of 100 characters each, with charLimit=4000
	// Should combine all into one ChunkSlice
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
		},
	}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(2))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 chunks in result, got %d", len(result))
	}
	if totalChars(result) != 300 {
		t.Errorf("expected 300 chars total, got %d", totalChars(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF on next Next(), got %v", err)
	}
}

// TestByteLimitRespected verifies char limit is not exceeded
func TestByteLimitRespected(t *testing.T) {
	// Chunks of 1500 characters each, with charLimit=4000
	// First two chunks = 3000 chars (fits)
	// Adding third chunk would be 4500 chars (exceeds limit)
	// So first result should have 2 chunks
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 1500))),
			textproc.Chunk(string(make([]rune, 1500))),
			textproc.Chunk(string(make([]rune, 1500))),
		},
	}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(1))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 chunks in first result, got %d", len(result))
	}
	if totalChars(result) > 4000 {
		t.Errorf("char limit 4000 exceeded, got %d", totalChars(result))
	}

	// Next result should overlap last chunk and include next chunk
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 chunks in second result, got %d", len(result))
	}
	if totalChars(result) > 4000 {
		t.Errorf("char limit 4000 exceeded, got %d", totalChars(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestOverlapBehavior verifies last N chunks are included in next output
func TestOverlapBehavior(t *testing.T) {
	// Chunks of 500 characters each, with charLimit=4000 and overlap=2
	// First call: [A, B, C, D, E, F, G, H] = 4000 chars
	// Second call (overlap=2): [G, H, I, J, K, L, M, N] = 4000 chars
	chunkData := make([]textproc.Chunk, 14)
	for i := range chunkData {
		// Create string with first char set to distinguish chunks
		s := make([]rune, 500)
		s[0] = rune('A' + i)
		chunkData[i] = textproc.Chunk(string(s))
	}
	src := &mockChunkSource{chunks: chunkData}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(2))

	// First result
	result1, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result1) != 8 {
		t.Errorf("expected 8 chunks, got %d", len(result1))
	}
	if len(result1) >= 8 && len(result1[0]) > 0 {
		if []rune(result1[0])[0] != 'A' || []rune(result1[7])[0] != 'H' {
			t.Errorf("first result should contain A-H, got %c to %c", []rune(result1[0])[0], []rune(result1[7])[0])
		}
	}

	// Second result should overlap last 2 chunks (G, H)
	result2, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result2) != 8 {
		t.Errorf("expected 8 chunks, got %d", len(result2))
	}
	if len(result2) >= 8 && len(result2[0]) > 0 {
		if []rune(result2[0])[0] != 'G' || []rune(result2[7])[0] != 'N' {
			t.Errorf("second result should overlap from G-H and end at N, got %c to %c", []rune(result2[0])[0], []rune(result2[7])[0])
		}
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestProgressGuarantee verifies large overlap drops from front until progress possible
func TestProgressGuarantee(t *testing.T) {
	// Overlap of 5 chunks, each 1000 characters = 5000 chars in overlap
	// Char limit is 3000, so overlap alone exceeds limit
	// Should drop from front until a new chunk fits
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
		},
	}
	p := NewProcessor(src, WithCharLimit(3000), WithOverlap(5))

	// First result should have 3 chunks (3000 chars), less than overlap count
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 chunks to fit within limit, got %d", len(result))
	}
	if totalChars(result) > 3000 {
		t.Errorf("char limit exceeded, got %d", totalChars(result))
	}

	// Second result - overlap drops from front until new chunk fits
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(result))
	}
	// Should have chunks 4, 5, 6 (dropped 1, 2, 3 to fit)
	if totalChars(result) != 3000 {
		t.Errorf("expected exactly 3000 chars, got %d", totalChars(result))
	}

	// Third result - remaining chunk 7
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(result))
	}
	if totalChars(result) != 1000 {
		t.Errorf("expected 1000 chars, got %d", totalChars(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestSingleChunkExceedsLimit verifies a single chunk > limit is returned anyway
func TestSingleChunkExceedsLimit(t *testing.T) {
	// Single chunk of 5000 characters, charLimit=4000
	// Should return it anyway
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 5000))),
		},
	}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(1))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(result))
	}
	if len(result[0]) != 5000 {
		t.Errorf("expected 5000 char chunk, got %d", len(result[0]))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestEOFHandling verifies remaining chunks are returned on final call, then EOF on next
func TestEOFHandling(t *testing.T) {
	// Chunks of 500 characters each, charLimit=4000
	// Only 5 chunks = 2500 chars total (less than limit)
	// Should return all on first Next(), then EOF on second
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 500))),
			textproc.Chunk(string(make([]rune, 500))),
			textproc.Chunk(string(make([]rune, 500))),
			textproc.Chunk(string(make([]rune, 500))),
			textproc.Chunk(string(make([]rune, 500))),
		},
	}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(2))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 5 {
		t.Errorf("expected 5 chunks, got %d", len(result))
	}
	if totalChars(result) != 2500 {
		t.Errorf("expected 2500 chars, got %d", totalChars(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF on second Next(), got %v", err)
	}
}

// TestConfigurationOptions verifies WithCharLimit and WithOverlap work correctly
func TestConfigurationOptions(t *testing.T) {
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
			textproc.Chunk(string(make([]rune, 1000))),
		},
	}

	// Test WithCharLimit
	p1 := NewProcessor(src, WithCharLimit(1500), WithOverlap(0))
	result, err := p1.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if totalChars(result) != 1000 {
		t.Errorf("WithCharLimit(1500) should give 1000 chars, got %d", totalChars(result))
	}

	// Test WithOverlap
	src2 := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
			textproc.Chunk(string(make([]rune, 100))),
		},
	}
	p2 := NewProcessor(src2, WithCharLimit(300), WithOverlap(2))
	result1, _ := p2.Next()
	result2, _ := p2.Next()

	// result2 should contain last 2 chunks from result1
	if len(result2) == 0 {
		t.Fatal("result2 should not be empty")
	}
	if len(result1) >= 2 && len(result2) >= 2 && len(result1[0]) > 0 && len(result2[0]) > 0 {
		// First chunk of result2 should be second-to-last chunk of result1
		if []rune(result2[0])[0] != []rune(result1[len(result1)-2])[0] {
			t.Errorf("overlap not working correctly")
		}
	}

	// Consume remaining chunks
	_, err = p2.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = p2.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = p2.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestNoOverlap verifies overlap=0 returns no overlap
func TestNoOverlap(t *testing.T) {
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			"aaa",
			"bbb",
			"ccc",
			"ddd",
			"eee",
		},
	}
	p := NewProcessor(src, WithCharLimit(10), WithOverlap(0))

	result1, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result2, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// result2 should not contain any chunks from result1
	for _, c1 := range result1 {
		for _, c2 := range result2 {
			if string(c1) == string(c2) {
				t.Errorf("overlap should be 0, but found overlapping chunk: %s", string(c1))
			}
		}
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

// TestEmptySource returns EOF immediately
func TestEmptySource(t *testing.T) {
	src := &mockChunkSource{chunks: []textproc.Chunk{}}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(2))

	_, err := p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF immediately, got %v", err)
	}
}

// TestChunkContentsPreserved verifies chunk data integrity
func TestChunkContentsPreserved(t *testing.T) {
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			"hello world",
			"foo bar baz",
			"qux quux",
		},
	}
	p := NewProcessor(src, WithCharLimit(1000), WithOverlap(1))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify all chunks are present with correct content
	expectedContents := []string{"hello world", "foo bar baz", "qux quux"}
	if len(result) != len(expectedContents) {
		t.Fatalf("expected %d chunks, got %d", len(expectedContents), len(result))
	}

	for i, expected := range expectedContents {
		if string(result[i]) != expected {
			t.Errorf("chunk %d: expected %q, got %q", i, expected, string(result[i]))
		}
	}
}

// TestProcessorImplementsChunkCombiner verifies Processor satisfies interface
func TestProcessorImplementsChunkCombiner(t *testing.T) {
	src := &mockChunkSource{chunks: []textproc.Chunk{"test"}}
	p := NewProcessor(src, WithCharLimit(4000), WithOverlap(2))

	var _ textproc.ChunkCombiner = p
}

// TestDefaultByteLimit uses default char limit
func TestDefaultByteLimit(t *testing.T) {
	// Create chunks that sum to > 4000 chars (default)
	src := &mockChunkSource{
		chunks: []textproc.Chunk{
			textproc.Chunk(string(make([]rune, 2000))),
			textproc.Chunk(string(make([]rune, 2000))),
			textproc.Chunk(string(make([]rune, 2000))),
		},
	}
	p := NewProcessor(src) // Use defaults

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should fit exactly 2 chunks (4000 chars) with default limit
	if totalChars(result) != 4000 {
		t.Errorf("expected 4000 chars with default limit, got %d", totalChars(result))
	}
}
