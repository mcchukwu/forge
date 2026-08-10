package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	base := t.TempDir()

	newDir := filepath.Join(base, "a", "b")
	if err := ensureDir(newDir); err != nil {
		t.Fatalf("ensureDir new: %v", err)
	}
	if info, err := os.Stat(newDir); err != nil {
		t.Fatalf("dir not created: %v", err)
	} else if !info.IsDir() {
		t.Fatal("created path is not a directory")
	}

	if err := ensureDir(newDir); err != nil {
		t.Fatalf("ensureDir existing: %v", err)
	}

	filePath := filepath.Join(base, "file")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ensureDir(filePath); err == nil {
		t.Fatal("expected error for existing file")
	}
}

func TestDirIsEmpty(t *testing.T) {
	base := t.TempDir()

	empty, err := dirIsEmpty(base)
	if err != nil {
		t.Fatalf("dirIsEmpty: %v", err)
	}
	if !empty {
		t.Fatal("expected directory to be empty")
	}

	if err := os.WriteFile(filepath.Join(base, "f"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	empty, err = dirIsEmpty(base)
	if err != nil {
		t.Fatalf("dirIsEmpty: %v", err)
	}
	if empty {
		t.Fatal("expected directory to be non-empty")
	}
}

func TestCreateStructure(t *testing.T) {
	base := t.TempDir()

	dirs, err := createStructure(base)
	if err != nil {
		t.Fatalf("createStructure: %v", err)
	}

	if len(dirs) != 3 {
		t.Fatalf("expected 3 dirs, got %v", dirs)
	}
	for _, d := range dirs {
		info, err := os.Stat(filepath.Join(base, d))
		if err != nil {
			t.Fatalf("dir %s: %v", d, err)
		}
		if !info.IsDir() {
			t.Fatalf("%s is not a directory", d)
		}
	}
}
