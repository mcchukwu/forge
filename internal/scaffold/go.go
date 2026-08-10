package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// initGoModule initializes a Go module at path with the given module name.
// It is a no-op when a go.mod already exists, so re-running forge is safe.
func initGoModule(path string, module string) (bool, error) {
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return false, nil
	}

	cmd := exec.Command("go", "mod", "init", module)
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err != nil {
		return false, fmt.Errorf("go mod init: %w\n%s", err, out)
	}

	return true, nil
}
