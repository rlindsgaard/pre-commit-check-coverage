package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockCommandRunner is a mock implementation of the CommandRunner interface
type MockCommandRunner struct {
	mock.Mock
}

func (m *MockCommandRunner) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCommandRunner) Output() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestGetStagedFiles(t *testing.T) {
	// Create a mock CommandRunner
	mockRunner := new(MockCommandRunner)

	// Define the mocked output of the command
	mockOutput := []byte("A\tnewfile.txt\nM\tmodifiedfile.txt\nR100\toldfile.txt\tnewfile.txt\nD\tdeletedfile.txt\n")
	mockRunner.On("Run").Return(nil)
	mockRunner.On("Output").Return(mockOutput, nil)

	// Call the function under test
	files, err := getStagedFiles(mockRunner)
	if err != nil {
		t.Fatalf("getStagedFiles returned an error: %v", err)
	}

	// Validate the results
	expectedFiles := []string{"newfile.txt", "modifiedfile.txt", "newfile.txt"}
	if len(files) != len(expectedFiles) {
		t.Fatalf("Expected %d files, got %d", len(expectedFiles), len(files))
	}
	for i, file := range files {
		if file != expectedFiles[i] {
			t.Errorf("Expected file %s, got %s", expectedFiles[i], file)
		}
	}

	// Assert that the mock methods were called as expected
	mockRunner.AssertExpectations(t)
}