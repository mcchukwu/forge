package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var errDirExists = errors.New("path exists and is not a directory")

// ensureDir creates path if it does not exist, or verifies that it is an
// existing directory. It returns nil when path exists and is a directory.
func ensureDir(path string) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("mkdirall %s: %w", path, err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !info.IsDir() {
		return errDirExists
	}

	return nil
}

// dirIsEmpty reports whether path is an empty directory.
func dirIsEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("read dir %s: %w", path, err)
	}
	return len(entries) == 0, nil
}

// createStructure creates the parent directories of the project layout and
// returns the directories created, relative to base.
func createStructure(base string) ([]string, error) {
	dirs := []string{"cmd", "internal", "pkg"}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(base, d), 0755); err != nil {
			return nil, fmt.Errorf("create structure: %w", err)
		}
	}

	return dirs, nil
}
