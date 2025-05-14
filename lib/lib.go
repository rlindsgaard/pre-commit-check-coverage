package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rlindsgaard/pre-commit-check-coverage/internal/hash"
)

func Verify(sha25map map[string][]string) ([]string, error) {
	// Get staged files (excluding deleted files, handling renames)
	stagedFiles, err := getStagedFiles()
	if err != nil {
		fmt.Printf("Error retrieving staged files: %v\n", err)
		os.Exit(1)
	}
	
	// Check each staged file against the sha256sums.txt
	var missingFiles []string
	for _, file := range stagedFiles {
		// Compute the SHA256 checksum of the file
		checksum, err := hash.ComputeSHA256(file)
		if err != nil {
			fmt.Printf("Error computing SHA256 for %s: %v\n", file, err)
			os.Exit(1)
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

// getStagedFiles retrieves the list of staged files (excluding deleted files, handling renames)
func getStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Parse the output to filter files
	var files []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		switch status {
		case "A", "M", "C": // Added, Modified, Copied
			files = append(files, parts[1])
		case "R": // Renamed
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
