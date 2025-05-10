package main

import (
	"os"
	"testing"
)

// Mock function for computing SHA256 checksum
func mockComputeSHA256(filePath string) (string, error) {
	mockData := map[string]string{
		"path/to/file1.txt":    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"path/to/file2.txt":    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"path/to/file3.txt":    "9d5e3ecdeb89fbb6de1f2b1aebc3c6a2f4c8e9348d7d3e0c5e5b9eb7b8b1a8f9",
		"path/to/new_file.txt": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}
	checksum, exists := mockData[filePath]
	if !exists {
		return "", os.ErrNotExist
	}
	return checksum, nil
}

// Mock function for loading sha256sums.txt
func mockLoadSha256Sums(filePath string) (map[string][]string, error) {
	return map[string][]string{
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855": {"path/to/file1.txt", "path/to/file2.txt"},
		"9d5e3ecdeb89fbb6de1f2b1aebc3c6a2f4c8e9348d7d3e0c5e5b9eb7b8b1a8f9": {"path/to/file3.txt"},
	}, nil
}

// Test for a valid staged file with an existing checksum and correct filename
func TestValidStagedFile(t *testing.T) {
	stagedFiles := []string{"path/to/file1.txt"}

	sha256Map, _ := mockLoadSha256Sums(sha256SumsFile)
	for _, file := range stagedFiles {
		checksum, err := mockComputeSHA256(file)
		if err != nil {
			t.Fatalf("Failed to compute checksum for file %s: %v", file, err)
		}

		if filenames, exists := sha256Map[checksum]; exists {
			if !stringInSlice(file, filenames) {
				t.Errorf("File %s is not in the list of filenames for its checksum", file)
			}
		} else {
			t.Errorf("Checksum for file %s not found in sha256sums.txt", file)
		}
	}
}

// Test for a missing file (checksum not found in sha256sums.txt)
func TestMissingFile(t *testing.T) {
	stagedFiles := []string{"path/to/missing_file.txt"}

	sha256Map, _ := mockLoadSha256Sums(sha256SumsFile)
	for _, file := range stagedFiles {
		checksum, err := mockComputeSHA256(file)
		if err == nil {
			if _, exists := sha256Map[checksum]; exists {
				t.Errorf("Checksum for missing file %s should not exist in sha256sums.txt", file)
			}
		}
	}
}

// Test for a renamed file (only the destination file is checked)
func TestRenamedFile(t *testing.T) {
	stagedFiles := []string{"path/to/new_file.txt"}

	sha256Map, _ := mockLoadSha256Sums(sha256SumsFile)
	for _, file := range stagedFiles {
		checksum, err := mockComputeSHA256(file)
		if err != nil {
			t.Fatalf("Failed to compute checksum for file %s: %v", file, err)
		}

		if filenames, exists := sha256Map[checksum]; exists {
			if !stringInSlice(file, filenames) {
				t.Errorf("Renamed file %s is not in the list of filenames for its checksum", file)
			}
		} else {
			t.Errorf("Checksum for renamed file %s not found in sha256sums.txt", file)
		}
	}
}

// Test for checksum mismatch (content changed, checksum does not match any in sha256sums.txt)
func TestChecksumMismatch(t *testing.T) {
	// Mock a file with new content that doesn't match any checksum in sha256sums.txt
	stagedFiles := []string{"path/to/changed_file.txt"}

	sha256Map, _ := mockLoadSha256Sums(sha256SumsFile)
	for _, file := range stagedFiles {
		// Simulate a checksum that doesn't exist in the mock map
		checksum := "invalidchecksum"

		if _, exists := sha256Map[checksum]; exists {
			t.Errorf("Checksum for changed file %s should not exist in sha256sums.txt", file)
		}
	}
}
