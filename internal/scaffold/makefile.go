package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

func createMakefile(base string) error {
	makefilePath := filepath.Join(base, "Makefile")

	if _, err := os.Stat(makefilePath); err == nil {
		return nil
	}

	appName := filepath.Base(base)

	makefileContent := fmt.Sprintf(`APP_NAME := %s
CMD_PATH := ./cmd
BIN_PATH := bin

.PHONY: run build clean test

run: 
	go run $(CMD_PATH)

build:
	go build -o $(BIN_PATH)/$(APP_NAME) $(CMD_PATH)

clean:
	rm -rf $(BIN_PATH)

test:
	go test ./...
	`, appName)

	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("create file Makefile failed: %w", err)
	}

	return nil
}
