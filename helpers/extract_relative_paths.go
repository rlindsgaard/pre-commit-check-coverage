package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flags
	rootPath := flag.String("root", ".", "Path to the repository root")
	coverFile := flag.String("coverfile", "cover.out", "Path to the cover.out file")
	flag.Parse()

	// Open the cover.out file
	file, err := os.Open(*coverFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening cover file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Use a map to store unique relative paths
	uniquePaths := make(map[string]struct{})

	// Read the cover.out file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Each line is in the format "<absolute-path>:<other-data>"
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 1 {
			continue
		}

		absolutePath := parts[0]
		if absolutePath == "mode" {
			continue
		}
		// Convert the absolute path to a relative path
		relativePath, err := filepath.Rel(*rootPath, absolutePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting path to relative: %v\n", err)
			continue
		}

		// Add the relative path to the map
		uniquePaths[relativePath] = struct{}{}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading cover file: %v\n", err)
		os.Exit(1)
	}

	// Print the unique relative paths
	for path := range uniquePaths {
		fmt.Println(path)
	}
}
