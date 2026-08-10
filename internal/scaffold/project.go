package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcchukwu/forge/internal/cli"
)

// Run resolves the target path for the given options and scaffolds a new
// Go project there.
func Run(opts cli.Options) error {
	path, err := resolvePath(opts.Name)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	return createProject(path, opts)
}

func resolvePath(name string) (string, error) {
	if name == "." {
		return os.Getwd()
	}
	return filepath.Abs(name)
}

func createProject(path string, opts cli.Options) error {
	module := opts.Module
	if module == "" {
		module = filepath.Base(path)
	}

	app := filepath.Base(path)
	display := opts.Name
	if opts.Name == "." {
		display = app
	}

	fmt.Printf("Creating project %q at %s\n", display, path)
	fmt.Printf("  module %s\n", module)

	if err := ensureDir(path); err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	empty, err := dirIsEmpty(path)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	if !empty && !opts.Force {
		return fmt.Errorf("create project: directory %s is not empty (use --force to overwrite)", path)
	}

	dirs, err := createStructure(path)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	for _, d := range dirs {
		fmt.Printf("  created %s/\n", d)
	}

	if created, err := initGoModule(path, module); err != nil {
		return fmt.Errorf("create project: %w", err)
	} else if created {
		fmt.Println("  created go.mod")
	}

	files, err := createFiles(path, module)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	for _, f := range files {
		fmt.Printf("  created %s\n", f)
	}

	if opts.HasMake {
		if created, err := createMakefile(path, app); err != nil {
			return fmt.Errorf("create project: %w", err)
		} else if created {
			fmt.Println("  created Makefile")
		}
	}

	if err := initGit(path); err != nil {
		if errors.Is(err, errGitNotInstalled) {
			fmt.Println("  skipped git init (git not found in PATH)")
		} else {
			fmt.Printf("  warning: %v\n", err)
		}
	} else {
		fmt.Println("  git initialized")
	}

	fmt.Println("Done.")
	if opts.Name != "." {
		fmt.Printf("Next: cd %s && go run ./cmd/%s\n", display, app)
	}

	return nil
}
