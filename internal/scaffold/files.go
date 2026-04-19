package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

func createFiles(base string, module string) error {
	// create main.go file
	mainPath := filepath.Join(base, "cmd", "main.go")
	mainContent := fmt.Sprintf(`package main
	
import (
	"fmt"
)

func main() {
	fmt.Println("Starting %s")
}
	`, module)

	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("create file main.go failed: %w", err)
	}

	// create README.md file
	readMePath := filepath.Join(base, "README.md")
	readMeContent := fmt.Sprintf(`#This project was started by forge. 
repo @ https://github.com/mcchukwu/forge
	`)

	if err := os.WriteFile(readMePath, []byte(readMeContent), 0644); err != nil {
		return fmt.Errorf("create file README.md failed: %w", err)
	}

	// create .gitignore file
	gitignorePath := filepath.Join(base, ".gitignore")
	gitignoreContent := fmt.Sprintf(`bin/*
.log

# Extend ignore patterns here
	`)

	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("create file .gitignore failed: %w", err)
	}

	return nil
}
