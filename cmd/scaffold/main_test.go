package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/bketelsen/toolbox"
)

// allTemplates maps a filename to its template string for parse testing.
var allTemplates = map[string]string{
	"go.mod":           tmplGoMod,
	"go.sum":           tmplGoSum,
	"main.go":          tmplMainGo,
	"cmd/root.go":      tmplCmdRootGo,
	"cmd/root_test.go": tmplCmdRootTestGo,
	"cmd/greet.go":     tmplCmdGreetGo,
	"cmd/process.go":   tmplCmdProcessGo,
	"Makefile":         tmplMakefile,
	".gitignore":       tmplGitignore,
	"README.md":        tmplReadme,
}

func TestTemplatesParse(t *testing.T) {
	for name, tmpl := range allTemplates {
		t.Run(name, func(t *testing.T) {
			_, err := template.New(name).Parse(tmpl)
			if err != nil {
				t.Errorf("template %q failed to parse: %v", name, err)
			}
		})
	}
}

func TestRunScaffoldCreatesFiles(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "testapp")

	flagModule = "github.com/test/testapp"
	flagOutput = projectDir
	flagDesc = "A test project"
	flagToolboxPath = "/tmp/toolbox"
	flagInteractive = false
	toolbox.DryRun = false
	toolbox.JSONOutput = false

	if err := runScaffold(nil, []string{"testapp"}); err != nil {
		t.Fatalf("runScaffold() error = %v", err)
	}

	expected := []string{
		"go.mod",
		"go.sum",
		"main.go",
		"cmd/root.go",
		"cmd/root_test.go",
		"cmd/greet.go",
		"cmd/process.go",
		"Makefile",
		".gitignore",
		"README.md",
	}

	for _, rel := range expected {
		path := filepath.Join(projectDir, filepath.FromSlash(rel))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %q was not created", rel)
		}
	}
}

func TestRunScaffoldGoModModuleName(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "myproject")

	flagModule = "github.com/example/myproject"
	flagOutput = projectDir
	flagDesc = "Test description"
	flagToolboxPath = "/tmp/toolbox"
	flagInteractive = false
	toolbox.DryRun = false
	toolbox.JSONOutput = false

	if err := runScaffold(nil, []string{"myproject"}); err != nil {
		t.Fatalf("runScaffold() error = %v", err)
	}

	gomod, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if !strings.Contains(string(gomod), "module github.com/example/myproject") {
		t.Errorf("go.mod does not contain expected module declaration; got:\n%s", gomod)
	}
}

func TestRunScaffoldNoSlugReferences(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "cleanapp")

	flagModule = "github.com/test/cleanapp"
	flagOutput = projectDir
	flagDesc = "Clean app with no old slug package"
	flagToolboxPath = "/tmp/toolbox"
	flagInteractive = false
	toolbox.DryRun = false
	toolbox.JSONOutput = false

	if err := runScaffold(nil, []string{"cleanapp"}); err != nil {
		t.Fatalf("runScaffold() error = %v", err)
	}

	// "toolbox/slug" is the removed package that must not appear in generated files.
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(content), "toolbox/slug") {
			t.Errorf("file %q contains unexpected 'toolbox/slug' reference", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking project dir: %v", err)
	}
}

func TestRunScaffoldDryRunProducesNoFiles(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "dryrunapp")

	flagModule = "github.com/test/dryrunapp"
	flagOutput = projectDir
	flagDesc = "Dry run test"
	flagToolboxPath = "/tmp/toolbox"
	flagInteractive = false
	toolbox.DryRun = true
	toolbox.JSONOutput = false

	defer func() { toolbox.DryRun = false }()

	if err := runScaffold(nil, []string{"dryrunapp"}); err != nil {
		t.Fatalf("runScaffold() error = %v", err)
	}

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		t.Errorf("dry run should not create project directory %q", projectDir)
	}
}
