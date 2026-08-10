package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

// createMakefile writes a Makefile in base. It reports whether a new file was
// written; existing Makefiles are left untouched. Recipe lines use real tabs,
// as required by make.
func createMakefile(base, app string) (bool, error) {
	makefilePath := filepath.Join(base, "Makefile")

	if _, err := os.Stat(makefilePath); err == nil {
		return false, nil
	}

	content := fmt.Sprintf(
		"APP_NAME := %s\n"+
			"CMD_PATH := ./cmd/%s\n"+
			"BIN_PATH := bin\n"+
			"\n"+
			".PHONY: run build clean test\n"+
			"\n"+
			"run:\n"+
			"\tgo run $(CMD_PATH)\n"+
			"\n"+
			"build:\n"+
			"\tgo build -o $(BIN_PATH)/$(APP_NAME) $(CMD_PATH)\n"+
			"\n"+
			"clean:\n"+
			"\trm -rf $(BIN_PATH)\n"+
			"\n"+
			"test:\n"+
			"\tgo test ./...\n",
		app, app)

	if err := os.WriteFile(makefilePath, []byte(content), 0644); err != nil {
		return false, fmt.Errorf("create Makefile: %w", err)
	}

	return true, nil
}
