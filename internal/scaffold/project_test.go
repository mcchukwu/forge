package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mcchukwu/forge/internal/cli"
)

const testApp = "demo-app"
const testModule = "example.com/demo-app"

func TestCreateProjectStructure(t *testing.T) {
	path := filepath.Join(t.TempDir(), testApp)
	opts := cli.Options{Name: testApp, Module: testModule, HasMake: true}

	if err := createProject(path, opts); err != nil {
		t.Fatalf("createProject: %v", err)
	}

	for _, d := range []string{"cmd", "internal", "pkg", "cmd/" + testApp} {
		info, err := os.Stat(filepath.Join(path, d))
		if err != nil {
			t.Fatalf("dir %s: %v", d, err)
		}
		if !info.IsDir() {
			t.Fatalf("%s is not a directory", d)
		}
	}

	for _, f := range []string{
		"go.mod",
		"README.md",
		".gitignore",
		"Makefile",
		"cmd/" + testApp + "/main.go",
	} {
		if _, err := os.Stat(filepath.Join(path, f)); err != nil {
			t.Fatalf("file %s: %v", f, err)
		}
	}

	if _, err := os.Stat(filepath.Join(path, "cmd", "main.go")); !os.IsNotExist(err) {
		t.Fatal("cmd/main.go should not exist in the new layout")
	}
}

func TestCreateProjectContents(t *testing.T) {
	path := filepath.Join(t.TempDir(), testApp)
	opts := cli.Options{Name: testApp, Module: testModule, HasMake: true}

	if err := createProject(path, opts); err != nil {
		t.Fatalf("createProject: %v", err)
	}

	assertFileContains(t, filepath.Join(path, "go.mod"), "module "+testModule)
	assertFileContains(t, filepath.Join(path, "cmd", testApp, "main.go"), testApp)
	assertFileContains(t, filepath.Join(path, "README.md"), "# "+testApp)
	assertFileContains(t, filepath.Join(path, "Makefile"), "APP_NAME := "+testApp)
	assertFileContains(t, filepath.Join(path, "Makefile"), "CMD_PATH := ./cmd/"+testApp)
}

func TestMakefileRecipesUseTabs(t *testing.T) {
	path := filepath.Join(t.TempDir(), testApp)
	opts := cli.Options{Name: testApp, Module: testModule, HasMake: true}

	if err := createProject(path, opts); err != nil {
		t.Fatalf("createProject: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(path, "Makefile"))
	if err != nil {
		t.Fatal(err)
	}

	recipes := []string{"go run", "go build", "rm -rf", "go test"}
	for _, line := range strings.Split(string(data), "\n") {
		for _, r := range recipes {
			if strings.Contains(line, r) && !strings.HasPrefix(line, "\t") {
				t.Errorf("recipe line must start with a tab: %q", line)
			}
		}
	}
}

func TestGoldenFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), testApp)
	opts := cli.Options{Name: testApp, Module: testModule}

	if err := createProject(path, opts); err != nil {
		t.Fatalf("createProject: %v", err)
	}

	for _, f := range []string{"README.md", ".gitignore"} {
		got, err := os.ReadFile(filepath.Join(path, f))
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		if want := readGolden(t, f); string(got) != want {
			t.Errorf("%s differs from golden:\n--- got ---\n%s\n--- want ---\n%s", f, got, want)
		}
	}
}

func TestCreateProjectRefusesNonEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), testApp)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "existing.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := cli.Options{Name: testApp}
	if err := createProject(path, opts); err == nil {
		t.Fatal("expected error for non-empty directory without --force")
	}

	opts.Force = true
	if err := createProject(path, opts); err != nil {
		t.Fatalf("expected success with --force, got: %v", err)
	}

	if _, err := os.Stat(filepath.Join(path, "go.mod")); err != nil {
		t.Fatalf("go.mod missing after --force: %v", err)
	}
}

func TestRunScaffoldsProject(t *testing.T) {
	chdir(t, t.TempDir())

	opts := cli.Options{Name: testApp, Module: testModule, HasMake: true}
	if err := Run(opts); err != nil {
		t.Fatalf("Run: %v", err)
	}

	mainPath := filepath.Join(testApp, "cmd", testApp, "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("main.go missing: %v", err)
	}
	assertFileContains(t, mainPath, testApp)
}

func assertFileContains(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), want) {
		t.Errorf("%s does not contain %q\n%s", path, want, data)
	}
}

func readGolden(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "golden", name))
	if err != nil {
		t.Fatalf("read golden %s: %v", name, err)
	}
	return string(data)
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Logf("failed to restore cwd: %v", err)
		}
	})
}
