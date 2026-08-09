package scaffold

import (
	"os/exec"
)

// initGoModule initializes go module on a given path and returns err or nil
func initGoModule(path string, module string) error {
	cmd := exec.Command("go", "mod", "init", module)
	cmd.Dir = path

	return cmd.Run()
}
