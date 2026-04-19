package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcchukwu/forge/internal/cli"
)

func Run(opts cli.Options) error {
	path, err := resolvePath(opts.Name)
	if err != nil {
		return fmt.Errorf("Could not resolve path: %v", err)
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
	if err := ensureDir(path); err != nil {
		return fmt.Errorf("create project failed: %w", err)
	}

	if err := createStructure(path); err != nil {
		return fmt.Errorf("create project failed: %w", err)
	}

	module := opts.Module
	if module == "" {
		module = filepath.Base(path)
	}

	if err := initGoModule(path, module); err != nil {
		return fmt.Errorf("create project failed: %w", err)
	}

	_ = initGit(path)

	if err := createFiles(path, module); err != nil {
		return fmt.Errorf("create project failed: %w", err)
	}

	if opts.HasMake {
		if err := createMakefile(path); err != nil {
			return fmt.Errorf("Create project failed: %w", err)
		}
	}

	return nil
}
