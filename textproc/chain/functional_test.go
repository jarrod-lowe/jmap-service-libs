package chain

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

var updateGolden = flag.Bool("update", false, "Update golden test files")

// testCaseMetadata holds optional metadata for a test case
type testCaseMetadata struct {
	Charset          string `json:"charset"`
	TransferEncoding string `json:"encoding"`
}

// testResult holds the expected output format
type testResult struct {
	Error *string `json:"error,omitempty"`
}

// TestFunctional runs end-to-end functional tests for the chain pipeline.
// It processes input files from testdata/ subdirectories and compares against expected outputs.
func TestFunctional(t *testing.T) {
	testdataDir := "testdata"

	// Find all test case directories
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	var testCaseDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			testCaseDirs = append(testCaseDirs, entry.Name())
		}
	}

	if len(testCaseDirs) == 0 {
		t.Fatal("No test case directories found in testdata/")
	}

	// Run each test case
	for _, caseName := range testCaseDirs {
		t.Run(caseName, func(t *testing.T) {
			caseDir := filepath.Join(testdataDir, caseName)

			// Load input
			input, err := loadTestCaseInput(caseDir)
			if err != nil {
				t.Fatalf("Failed to load test case input: %v", err)
			}

			// Load metadata (optional)
			meta, err := loadTestCaseMetadata(caseDir)
			if err != nil && !os.IsNotExist(err) {
				t.Fatalf("Failed to load test case metadata: %v", err)
			}

			// Get charset and transfer encoding, default to empty string if no metadata
			var charset, transferEncoding string
			if meta != nil {
				charset = meta.Charset
				transferEncoding = meta.TransferEncoding
			}

			// Run chain
			results, err := runChain(input, charset, transferEncoding)

			// Encode results
			encoded, err := encodeResults(results, err)
			if err != nil {
				t.Fatalf("Failed to encode results: %v", err)
			}

			// Update golden file or compare
			goldenPath := filepath.Join(caseDir, "expected.json")
			if *updateGolden {
				if err := os.WriteFile(goldenPath, encoded, 0600); err != nil { // #nosec G306 -- Test data file
					t.Fatalf("Failed to write golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenPath)
			} else {
				expected, err := os.ReadFile(goldenPath) // #nosec G304 -- Test data from controlled directory
				if err != nil {
					t.Fatalf("Failed to read golden file: %v", err)
				}

				compareResults(t, encoded, expected)
			}
		})
	}
}

// loadTestCaseInput reads the input.txt file from a test case directory
func loadTestCaseInput(caseDir string) (io.Reader, error) {
	inputPath := filepath.Join(caseDir, "input.txt")
	data, err := os.ReadFile(inputPath) // #nosec G304 -- Test data from controlled directory
	if err != nil {
		return nil, fmt.Errorf("failed to read input.txt: %w", err)
	}
	return strings.NewReader(string(data)), nil
}

// loadTestCaseMetadata reads the meta.json file from a test case directory
func loadTestCaseMetadata(caseDir string) (*testCaseMetadata, error) {
	metaPath := filepath.Join(caseDir, "meta.json")
	data, err := os.ReadFile(metaPath) // #nosec G304 -- Test data from controlled directory
	if err != nil {
		return nil, err
	}

	var meta testCaseMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse meta.json: %w", err)
	}

	return &meta, nil
}

// runChain executes the chain pipeline and collects all outputs
func runChain(r io.Reader, charset, transferEncoding string) ([]textproc.ChunkSlice, error) {
	chain, err := NewReaderConfigWithEncoding(r, DefaultMaxBytes, DefaultOverlap, DefaultByteLimit, charset, transferEncoding)
	if err != nil {
		return nil, err
	}

	var results []textproc.ChunkSlice
	for {
		chunk, err := chain.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return results, err
		}
		results = append(results, chunk)
	}

	return results, nil
}

// encodeResults converts chain results to JSON format for golden file comparison
func encodeResults(results []textproc.ChunkSlice, err error) ([]byte, error) {
	if err != nil {
		// Error case
		errMsg := err.Error()
		result := testResult{Error: &errMsg}
		return json.MarshalIndent(result, "", "  ")
	}

	// Success case: encode each ChunkSlice as array of base64-encoded byte arrays
	encoded := make([][]string, 0, len(results))
	for _, slice := range results {
		encodedSlice := make([]string, 0, len(slice))
		for _, chunk := range slice {
			encodedSlice = append(encodedSlice, base64.StdEncoding.EncodeToString([]byte(chunk)))
		}
		encoded = append(encoded, encodedSlice)
	}

	return json.MarshalIndent(encoded, "", "  ")
}

// compareResults compares actual and expected results
func compareResults(t *testing.T, actual, expected []byte) {
	var actualJSON, expectedJSON interface{}
	if err := json.Unmarshal(actual, &actualJSON); err != nil {
		t.Fatalf("Failed to unmarshal actual result: %v", err)
	}
	if err := json.Unmarshal(expected, &expectedJSON); err != nil {
		t.Fatalf("Failed to unmarshal expected result: %v", err)
	}

	if !jsonEqual(actualJSON, expectedJSON) {
		t.Errorf("Results differ:\nActual:\n%s\n\nExpected:\n%s", string(actual), string(expected))
	}
}

// jsonEqual compares two JSON values deeply
func jsonEqual(a, b interface{}) bool {
	switch a := a.(type) {
	case map[string]interface{}:
		bMap, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		if len(a) != len(bMap) {
			return false
		}
		for key, aVal := range a {
			bVal, ok := bMap[key]
			if !ok || !jsonEqual(aVal, bVal) {
				return false
			}
		}
		return true
	case []interface{}:
		bSlice, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if !jsonEqual(a[i], bSlice[i]) {
				return false
			}
		}
		return true
	case string:
		bStr, ok := b.(string)
		return ok && a == bStr
	case float64:
		bFloat, ok := b.(float64)
		return ok && a == bFloat
	case bool:
		bBool, ok := b.(bool)
		return ok && a == bBool
	case nil:
		return b == nil
	default:
		return false
	}
}
