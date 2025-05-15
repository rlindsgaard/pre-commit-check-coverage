package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestVerifyEmptyDiff(t *testing.T) {
	mockRunner := new(MockCommandRunner)
	
	mockOutput := []byte("\n")
	
	mockRunner.On("Run").Return(nil)
	mockRunner.On("Output").Return(mockOutput, nil)
	
	checksums := map[string][]string{
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855": {"path/to/file1.txt", "path/to/file2.txt"},
		"9d5e3ecdeb89fbb6de1f2b1aebc3c6a2f4c8e9348d7d3e0c5e5b9eb7b8b1a8f9": {"path/to/file3.txt"},
	}
	
	files, err := Verify(checksums, mockRunner)
	if err != nil {
		t.Fatalf("Verify return an err: %v", err)
	}
	
	if files != nil {
		t.Errorf("Expected no missing files but found: %v", files)
	}
}

func TestDiffFilesCovered(t *testing.T) {
	// Arrange
	mockRunner := new(MockCommandRunner)
	
	mockOutput := []byte("A\tpath/to/file1.txt\nM\tpath/to/file3.txt")
	mockRunner.On("Run").Return(nil)
	mockRunner.On("Output").Return(mockOutput)
	
	checksums := map[string][]string{
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855": {"path/to/file1.txt", "path/to/file2.txt"},
		"9d5e3ecdeb89fbb6de1f2b1aebc3c6a2f4c8e9348d7d3e0c5e5b9eb7b8b1a8f9": {"path/to/file3.txt"},
	}
	// Act
	missingFiles, err := Verify(checksums, mockRunner)
	
	// Assert
	assert.Nil(t, missingFiles)
	assert.Nil(t, err)
}


func TestGetStagedFiles(t *testing.T) {
	// Create a mock CommandRunner
	mockRunner := new(MockCommandRunner)

	// Define the mocked output of the command
	mockOutput := []byte("A\tnewfile.txt\nM\tmodifiedfile.txt\nR100\toldfile.txt\tnewfilelocation.txt\nD\tdeletedfile.txt\n")
	mockRunner.On("Run").Return(nil)
	mockRunner.On("Output").Return(mockOutput, nil)

	// Call the function under test
	files, err := getStagedFiles(mockRunner)
	if err != nil {
		t.Fatalf("getStagedFiles returned an error: %v", err)
	}

	// Validate the results
	expectedFiles := []string{"newfile.txt", "modifiedfile.txt", "newfilelocation.txt"}
	if len(files) != len(expectedFiles) {
		t.Fatalf("Expected %d files, got %d (%v)", len(expectedFiles), len(files), files)
	}
	for i, file := range files {
		if file != expectedFiles[i] {
			t.Errorf("Expected file %s, got %s", expectedFiles[i], file)
		}
	}

	// Assert that the mock methods were called as expected
	mockRunner.AssertExpectations(t)
}
