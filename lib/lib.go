package lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rlindsgaard/pre-commit-check-coverage/internal/hash"
)

// CommandRunner interface to abstract command execution
type CommandRunner interface {
	Run() error
	Output() ([]byte, error)
}

// RealCommandRunner is the real implementation of CommandRunner using exec.Command
type RealCommandRunner struct {
	cmd *exec.Cmd
}

func (r *RealCommandRunner) Run() error {
	return r.cmd.Run()
}

func (r *RealCommandRunner) Output() ([]byte, error) {
	return r.cmd.Output()
}

func Verify(sha256Map map[string][]string, commandRunner CommandRunner) ([]string, error) {
	// Get staged files (excluding deleted files, handling renames)
	
	stagedFiles, err := getStagedFiles(commandRunner)
	if err != nil {	
        return nil, fmt.Errorf("Could not retrieve stages files: %v", err)
	}
	
	// Check each staged file against the sha256sums.txt
	var missingFiles []string
	for _, file := range stagedFiles {
		// Compute the SHA256 checksum of the file
		checksum, err := hash.ComputeSHA256(file)
		if err != nil {
			return nil, fmt.Errof("Error computing SHA256 for %s: %v\n", file, err)
		
		}

		// Check if the checksum exists in the sha256sums.txt
		if filenames, exists := sha256Map[checksum]; exists {
			// Verify that the filename is in the list of filenames for this checksum
			if !stringInSlice(file, filenames) {
				missingFiles = append(missingFiles, file)
			}
		} else {
			// If the checksum is not found, the file is missing
			missingFiles = append(missingFiles, file)
		}
	}

	// If any files are missing, fail the commit
	if len(missingFiles) > 0 {
		err := errors.New("Found staged files not tested")
		return missingFiles, err
	}
	return nil, nil
}

// Helper function to create a new CommandRunner for the git command
func newGitCommandRunner() CommandRunner {
	return &RealCommandRunner{
		cmd: exec.Command("git", "diff", "--cached", "--name-status"),
	}
}

// getStagedFiles retrieves the list of staged files (excluding deleted files, handling renames)
func getStagedFiles(runner CommandRunner) ([]string, error) {
	runner.Run()

	output, err := runner.Output()
	if err != nil {
		return nil, err
	}

	// Parse the output to filter files
	var files []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Parsing line: '%s'", line)
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		status := fmt.Sprintf("%c", parts[0][0])
		fmt.Printf("%s", status)
		switch status {
		case `A`, `M`, `C`: // Added, Modified, Copied
			files = append(files, parts[1])
		case `R`: // Renamed
			files = append(files, parts[2]) // Only include the destination file
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

// stringInSlice checks if a string is present in a slice
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
