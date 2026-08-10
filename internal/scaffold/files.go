package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

// createFiles writes the project's source files and returns the list of
// files written, relative to base.
func createFiles(base, module string) ([]string, error) {
	app := filepath.Base(base)

	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting %s")
}
`, app)
	mainPath := filepath.Join(base, "cmd", app, "main.go")
	if err := os.MkdirAll(filepath.Dir(mainPath), 0755); err != nil {
		return nil, fmt.Errorf("create cmd/%s: %w", app, err)
	}
	if err := writeFile(mainPath, []byte(mainContent)); err != nil {
		return nil, fmt.Errorf("create main.go: %w", err)
	}

	readmeContent := fmt.Sprintf(`# %s

This project was scaffolded with [forge](https://github.com/mcchukwu/forge).
`, app)
	if err := writeFile(filepath.Join(base, "README.md"), []byte(readmeContent)); err != nil {
		return nil, fmt.Errorf("create README.md: %w", err)
	}

	gitignoreContent := `bin/*
*.log

# Add project-specific ignore patterns below.
`
	if err := writeFile(filepath.Join(base, ".gitignore"), []byte(gitignoreContent)); err != nil {
		return nil, fmt.Errorf("create .gitignore: %w", err)
	}

	return []string{
		filepath.Join("cmd", app, "main.go"),
		"README.md",
		".gitignore",
	}, nil
}

func writeFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}
