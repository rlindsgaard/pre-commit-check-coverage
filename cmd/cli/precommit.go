package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rlindsgaard/pre-commit-check-coverage/crypto"
)

// Path to sha256sums.txt
const sha256SumsFile = "sha256sums.txt"

func main() {
	// Ensure sha256sums.txt exists
	if _, err := os.Stat(sha256SumsFile); os.IsNotExist(err) {
		fmt.Printf("Error: %s file not found!\n", sha256SumsFile)
		os.Exit(1)
	}

	// Get staged files (excluding deleted files, handling renames)
	stagedFiles, err := getStagedFiles()
	if err != nil {
		fmt.Printf("Error retrieving staged files: %v\n", err)
		os.Exit(1)
	}

	// Read the sha256sums.txt file into a map (checksum -> list of filenames)
	sha256Map, err := loadSha256Sums(sha256SumsFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", sha256SumsFile, err)
		os.Exit(1)
	}

	// Check each staged file against the sha256sums.txt
	var missingFiles []string
	for _, file := range stagedFiles {
		// Compute the SHA256 checksum of the file
		checksum, err := crypto.ComputeSHA256(file)
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
		fmt.Println("Error: The following files are missing from the coverage report or have mismatched filenames:")
		for _, file := range missingFiles {
			fmt.Printf("  - %s\n", file)
		}
		fmt.Println("Commit failed. Ensure these files are included in the coverage report with the correct filenames.")
		os.Exit(1)
	}

	// If all checks pass, allow the commit
	os.Exit(0)
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

// loadSha256Sums reads sha256sums.txt into a map for quick lookup
// The map key is the checksum, and the value is a list of filenames
func loadSha256Sums(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sha256Map := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			checksum := parts[0]
			filename := parts[1]
			sha256Map[checksum] = append(sha256Map[checksum], filename)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sha256Map, nil
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
