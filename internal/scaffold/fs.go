package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var errDirExists = errors.New("path exists and is not a directory")

// ensureDir checks if the path does not exists and returns an error or nil if it does
func ensureDir(path string) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("MkdirAll failed: %w", err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("Stat failed: %w", err)
	}

	if !info.IsDir() {
		return errDirExists
	}

	return nil
}

// createStructure creates the parent folders or the code base
func createStructure(base string) error {
	dirs := []string{
		"cmd",
		"internal",
		"pkg",
	}

	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(base, d), 0755)
		if err != nil {
			return fmt.Errorf("Create structure failed: %w", err)
		}
	}

	return nil
}
