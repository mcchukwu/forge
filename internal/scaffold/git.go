package scaffold

import "os/exec"

func initGit(path string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return nil
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path

	return cmd.Run()
}
