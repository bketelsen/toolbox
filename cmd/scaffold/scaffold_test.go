package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/bketelsen/toolbox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// toolboxRoot returns the absolute path of the toolbox repository root.
// cmd/scaffold lives two levels below the module root.
func toolboxRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	root := filepath.Join(wd, "..", "..")
	_, err = os.Stat(filepath.Join(root, "go.mod"))
	require.NoError(t, err, "expected go.mod at toolbox root %s", root)
	return root
}

// scaffoldProject resets flag defaults and runs scaffold via runScaffold directly.
func scaffoldProject(t *testing.T, name, outDir string) {
	t.Helper()
	flagModule = ""
	flagOutput = outDir
	flagDesc = "A toolbox-powered Go project"
	flagToolboxPath = toolboxRoot(t)
	flagInteractive = false
	toolbox.DryRun = false
	toolbox.JSONOutput = false

	if err := runScaffold(nil, []string{name}); err != nil {
		t.Fatalf("runScaffold() error = %v", err)
	}
}

// TestScaffoldProducesExpectedFiles verifies scaffold writes all required files.
func TestScaffoldProducesExpectedFiles(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "myapp")

	scaffoldProject(t, "myapp", projectDir)

	expectedFiles := []string{
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

	for _, rel := range expectedFiles {
		path := filepath.Join(projectDir, filepath.FromSlash(rel))
		_, err := os.Stat(path)
		assert.NoError(t, err, "expected file %s to exist", rel)
	}
}

// TestScaffoldGoModContent verifies the generated go.mod has correct module path and replace directive.
func TestScaffoldGoModContent(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "testmod")

	flagModule = "github.com/example/testmod"
	flagOutput = projectDir
	flagDesc = "A toolbox-powered Go project"
	flagToolboxPath = toolboxRoot(t)
	flagInteractive = false
	toolbox.DryRun = false
	toolbox.JSONOutput = false

	require.NoError(t, runScaffold(nil, []string{"testmod"}))

	content, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "module github.com/example/testmod")
	assert.Contains(t, body, "github.com/bketelsen/toolbox v0.0.0-dev")
	assert.Contains(t, body, "replace github.com/bketelsen/toolbox =>")
	assert.NotContains(t, body, "/path/to/toolbox", "toolbox path should be resolved, not placeholder")
}

// TestScaffoldMainGoContent verifies main.go imports and uses toolbox.App.
func TestScaffoldMainGoContent(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "maintest")

	scaffoldProject(t, "maintest", projectDir)

	content, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "github.com/bketelsen/toolbox")
	assert.Contains(t, body, "toolbox.App{")
	assert.Contains(t, body, "app.Run(cmd.RootCmd)")
	assert.NotContains(t, body, "slug", "main.go must not reference the removed slug package")
}

// TestScaffoldRootGoContent verifies cmd/root.go structure and flag usage.
func TestScaffoldRootGoContent(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "roottest")

	scaffoldProject(t, "roottest", projectDir)

	content, err := os.ReadFile(filepath.Join(projectDir, "cmd", "root.go"))
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "toolbox.Verbose")
	assert.Contains(t, body, "toolbox.Silent")
	assert.Contains(t, body, "toolbox.BindViper")
	assert.Contains(t, body, "RootCmd.AddCommand(processCmd)")
	assert.Contains(t, body, "RootCmd.AddCommand(greetCmd)")
	assert.NotContains(t, body, "slug")
}

// TestScaffoldProcessGoContent verifies cmd/process.go demonstrates reporter, DryRun, JSON output.
func TestScaffoldProcessGoContent(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "processtest")

	scaffoldProject(t, "processtest", projectDir)

	content, err := os.ReadFile(filepath.Join(projectDir, "cmd", "process.go"))
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "toolbox.NewReporter()")
	assert.Contains(t, body, "toolbox.DryRun")
	assert.Contains(t, body, "toolbox.OutputJSON")
	assert.NotContains(t, body, "slug")
}

// TestScaffoldGreetGoContent verifies cmd/greet.go exists and has correct content.
func TestScaffoldGreetGoContent(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "greettest")

	scaffoldProject(t, "greettest", projectDir)

	content, err := os.ReadFile(filepath.Join(projectDir, "cmd", "greet.go"))
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "greetCmd")
	assert.Contains(t, body, "toolbox.OutputJSON")
	assert.Contains(t, body, "greet")
	assert.NotContains(t, body, "slug")
}

// TestScaffoldRootTestGoContent verifies the generated cmd/root_test.go checks for both subcommands.
func TestScaffoldRootTestGoContent(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "roottesttest")

	scaffoldProject(t, "roottesttest", projectDir)

	content, err := os.ReadFile(filepath.Join(projectDir, "cmd", "root_test.go"))
	require.NoError(t, err)

	body := string(content)
	assert.Contains(t, body, "TestRootCmdHasSubcommands")
	assert.Contains(t, body, "TestGreetCmdRegistered")
	assert.Contains(t, body, "TestProcessCmdRegistered")
}

// TestScaffoldNoSlugReferences verifies no generated file references the removed slug package.
func TestScaffoldNoSlugReferences(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "cleanapp")

	scaffoldProject(t, "cleanapp", projectDir)

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		rel, _ := filepath.Rel(projectDir, path)
		assert.NotContains(t, string(content), "slug",
			"file %s must not reference removed slug package", rel)
		return nil
	})
	require.NoError(t, err)
}

// TestScaffoldDryRun verifies dry-run mode produces no output files.
func TestScaffoldDryRun(t *testing.T) {
	outDir := t.TempDir()
	projectDir := filepath.Join(outDir, "dryrunapp")

	flagModule = ""
	flagOutput = projectDir
	flagDesc = "A toolbox-powered Go project"
	flagToolboxPath = toolboxRoot(t)
	flagInteractive = false
	toolbox.DryRun = true
	toolbox.JSONOutput = false
	defer func() { toolbox.DryRun = false }()

	require.NoError(t, runScaffold(nil, []string{"dryrunapp"}))

	_, statErr := os.Stat(projectDir)
	assert.True(t, os.IsNotExist(statErr), "dry-run must not create output directory")
}

// TestScaffoldTemplatesRenderCorrectly verifies templates render with correct substitutions.
func TestScaffoldTemplatesRenderCorrectly(t *testing.T) {
	data := templateData{
		Name:        "myproject",
		Module:      "github.com/myorg/myproject",
		Description: "My test project",
		ToolboxPath: "/some/path/toolbox",
	}

	tests := []struct {
		name        string
		tmpl        string
		contains    []string
		notContains []string
	}{
		{
			name: "go.mod",
			tmpl: tmplGoMod,
			contains: []string{
				"module github.com/myorg/myproject",
				"github.com/bketelsen/toolbox v0.0.0-dev",
				"replace github.com/bketelsen/toolbox => /some/path/toolbox",
			},
			notContains: []string{"slug"},
		},
		{
			name: "main.go",
			tmpl: tmplMainGo,
			contains: []string{
				"github.com/myorg/myproject/cmd",
				"toolbox.App{",
				"app.Run(cmd.RootCmd)",
			},
			notContains: []string{"slug"},
		},
		{
			name: "cmd/root.go",
			tmpl: tmplCmdRootGo,
			contains: []string{
				`Use:   "myproject"`,
				"My test project",
				"toolbox.Verbose",
				"toolbox.Silent",
				"RootCmd.AddCommand(processCmd)",
				"RootCmd.AddCommand(greetCmd)",
			},
			notContains: []string{"slug"},
		},
		{
			name: "cmd/root_test.go",
			tmpl: tmplCmdRootTestGo,
			contains: []string{
				"TestRootCmdHasSubcommands",
				"TestGreetCmdRegistered",
				"TestProcessCmdRegistered",
			},
			notContains: []string{"slug"},
		},
		{
			name: "cmd/greet.go",
			tmpl: tmplCmdGreetGo,
			contains: []string{
				"greetCmd",
				"toolbox.OutputJSON",
				"toolbox.DryRun",
			},
			notContains: []string{"slug"},
		},
		{
			name: "cmd/process.go",
			tmpl: tmplCmdProcessGo,
			contains: []string{
				"toolbox.NewReporter()",
				"toolbox.DryRun",
				"toolbox.OutputJSON",
				"r.Step(",
			},
			notContains: []string{"slug"},
		},
		{
			name: "Makefile",
			tmpl: tmplMakefile,
			contains: []string{
				"BINARY := myproject",
				"go test -v ./...",
				"go build",
			},
			notContains: []string{"slug"},
		},
		{
			name: "README.md",
			tmpl: tmplReadme,
			contains: []string{
				"# myproject",
				"My test project",
				"--json",
				"--verbose",
				"--dry-run",
				"--silent",
			},
			notContains: []string{"slug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := template.New(tt.name).Parse(tt.tmpl)
			require.NoError(t, err)

			var sb strings.Builder
			err = parsed.Execute(&sb, data)
			require.NoError(t, err)

			body := sb.String()
			for _, want := range tt.contains {
				assert.Contains(t, body, want, "template %s should contain %q", tt.name, want)
			}
			for _, notWant := range tt.notContains {
				assert.NotContains(t, body, notWant, "template %s must not contain %q", tt.name, notWant)
			}
		})
	}
}

// TestScaffoldFailsIfOutputExists verifies scaffold returns error when output dir already exists.
func TestScaffoldFailsIfOutputExists(t *testing.T) {
	existingDir := t.TempDir()

	flagModule = ""
	flagOutput = existingDir
	flagDesc = "A toolbox-powered Go project"
	flagToolboxPath = toolboxRoot(t)
	flagInteractive = false
	toolbox.DryRun = false
	toolbox.JSONOutput = false

	err := runScaffold(nil, []string{"conflict"})
	assert.Error(t, err, "scaffold should fail when output directory already exists")
	assert.Contains(t, err.Error(), "already exists")
}

// TestIsToolboxGoMod verifies the helper correctly identifies toolbox's go.mod.
func TestIsToolboxGoMod(t *testing.T) {
	root := toolboxRoot(t)
	assert.True(t, isToolboxGoMod(filepath.Join(root, "go.mod")))
	assert.False(t, isToolboxGoMod("/nonexistent/go.mod"))

	tmpFile, err := os.CreateTemp(t.TempDir(), "go.mod")
	require.NoError(t, err)
	_, err = tmpFile.WriteString("module github.com/other/project\n\ngo 1.25.0\n")
	require.NoError(t, err)
	_ = tmpFile.Close()
	assert.False(t, isToolboxGoMod(tmpFile.Name()))
}
