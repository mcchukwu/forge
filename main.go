package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Helpers that define the error system.
func fail(msg string) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
func check(err error) {
	if err != nil {
		fail(err.Error())
	}
}

// Helper that resolves the given path to an absolute path.
func resolvePath(name string) (string, error) {
	if name == "." {
		return os.Getwd()
	}

	return filepath.Abs(name)
}

// Helper that checks the given error and exits if it is not nil.
func ensureDir(path string) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("path exists and is not a directory")
	}

	return nil
}

// Helper that creates the structure of a new project.
func createStructure(base string) error {
	dirs := []string{
		"cmd",
		"internal",
		"pkg",
	}

	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(base, d), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper that initializes go modules in the given path.
func initGoModule(path string, module string) error {
	cmd := exec.Command("go", "mod", "init", module)
	cmd.Dir = path

	return cmd.Run()
}

// Helper that initializes git in the given path.
func initGit(path string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return nil // skip if git not installed
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = path

	return cmd.Run()
}

// Helper that creates the files of a new project.
func createFiles(base string, module string) error {
	mainPath := filepath.Join(base, "cmd", "main.go")

	mainContent := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Starting %s")
}
`, module)

	err := os.WriteFile(mainPath, []byte(mainContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Helper that creates a new project in the given path.
func createProject(path string, module string) error {
	check(ensureDir(path))
	check(createStructure(path))

	if module == "" {
		module = filepath.Base(path)
	}

	check(initGoModule(path, module))
	check(initGit(path))
	check(createFiles(path, module))

	return nil
}

// main function
func main() {
	if len(os.Args) < 3 {
		fail("Usage: forge new <project-name | .> [--module github.com/username/repo]")
	}

	cmd := os.Args[1]
	name := os.Args[2]

	if cmd != "new" {
		fail("Unknown command: " + cmd)
	}

	module := ""
	if len(os.Args) >= 5 && os.Args[3] == "--module" {
		module = os.Args[4]
	}

	path, err := resolvePath(name)
	check(err)

	check(createProject(path, module))
	fmt.Println("OK:", path)
}
