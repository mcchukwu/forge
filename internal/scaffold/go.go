package scaffold

import (
	"os/exec"
)

func initGoModule(path string, module string) error {
	cmd := exec.Command("go", "mod", "init", module)
	cmd.Dir = path

	return cmd.Run()
}
