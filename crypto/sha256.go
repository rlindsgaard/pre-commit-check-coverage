package crypto

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

func ComputeSHA256(filePath string) (string, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := bufio.NewReader(file).WriteTo(hasher); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
