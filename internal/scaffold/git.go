package scaffold

import "os/exec"

// initGit initializes git on a given path and returns err or nil
func initGit(path string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return nil
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path

	return cmd.Run()
}
