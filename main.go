package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	Name   string
	Module string
	Make   bool
}

// Helper that parses the command line arguments.
func parseFlags(args []string) (Options, error) {
	opts := Options{}
	opts.Name = args[2]

	for i := 3; i < len(args); i++ {
		switch args[i] {
		case "-m", "--module":
			if i+1 >= len(args) {
				return opts, errors.New("--module missing module name")
			}

			opts.Module = args[i+1]
			i++
		case "-n", "--name":
			if i+1 >= len(args) {
				return opts, errors.New("missing project name")
			}

			opts.Name = args[i+1]
			i++
		case "-M", "--make":
			opts.Make = true
		default:
			return opts, errors.New("unknown option: " + args[i])
		}
	}

	return opts, nil
}

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

// Helper that creates a Makefile in the given path.
func createMakefile(path string) error {
	makefilePath := filepath.Join(path, "Makefile")

	if _, err := os.Stat(makefilePath); err == nil {
		return nil
	}

	appName := filepath.Base(path)

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

	err := os.WriteFile(makefilePath, []byte(makefileContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Helper that creates a new project in the given path.
func createProject(path string, module string, makeEnabled bool) error {
	check(ensureDir(path))
	check(createStructure(path))

	if module == "" {
		module = filepath.Base(path)
	}

	check(initGoModule(path, module))
	check(initGit(path))
	check(createFiles(path, module))
	if makeEnabled {
		check(createMakefile(path))
	}

	return nil
}

// main function
func main() {
	if len(os.Args) < 2 {
		fail("Usage: forge new <project-name | .>")
	}

	cmd := os.Args[1]

	if cmd != "new" {
		fail("Unknown command: " + cmd)
	}

	opts, err := parseFlags(os.Args)
	check(err)

	path, err := resolvePath(opts.Name)
	check(err)

	check(createProject(path, opts.Module, opts.Make))
	fmt.Println("OK:", path)
}
