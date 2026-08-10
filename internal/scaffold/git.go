package scaffold

import (
	"errors"
	"fmt"
	"os/exec"
)

var errGitNotInstalled = errors.New("git not found in PATH")

// initGit initializes a git repository at path. It returns errGitNotInstalled
// when git is unavailable, so callers can skip without failing the scaffold.
func initGit(path string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return errGitNotInstalled
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init: %w\n%s", err, out)
	}

	return nil
}
