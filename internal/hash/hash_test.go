package hash_test

import (
	"os"
	"testing"

	"github.com/rlindsgaard/pre-commit-check-coverage/internal/hash"
)

func TestComputeSHA256(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the file afterward

	// Write some content to the temporary file
	content := []byte("test content")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Close the file to ensure it can be read by ComputeSHA256
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Call the ComputeSHA256 function
	result, err := hash.ComputeSHA256(tmpfile.Name())
	if err != nil {
		t.Fatalf("ComputeSHA256 returned an error: %v", err)
	}

	// Expected SHA256 hash for "test content"
	expected := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}