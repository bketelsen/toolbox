// Command scaffold generates a new Go project that uses toolbox as its foundation.
package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bketelsen/toolbox"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/cobra"
)

// templateData is the value passed to every file template.
type templateData struct {
	Name        string
	Module      string
	Description string
	ToolboxPath string // absolute path to toolbox clone, used in go.mod replace directive
}

// generatedFile records a single file that would be or was written.
type generatedFile struct {
	Path    string `json:"path"`
	Written bool   `json:"written"`
}

var (
	flagModule      string
	flagOutput      string
	flagDesc        string
	flagToolboxPath string
	flagInteractive bool
)

func main() {
	app := toolbox.App{
		Version: "0.1.0",
		Commit:  "none",
		Date:    "unknown",
		BuiltBy: "source",
	}
	_ = app.Run(rootCmd())
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold <project-name>",
		Short: "Generate a new Go project powered by toolbox",
		Args:  cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level := slog.LevelInfo
			if toolbox.Verbose {
				level = slog.LevelDebug
			}
			if toolbox.Silent {
				level = slog.LevelError + 1
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: level,
			})))
			return toolbox.BindViper(cmd.Root())
		},
		RunE: runScaffold,
	}

	cmd.Flags().StringVarP(&flagModule, "module", "m", "", "Go module path (default: github.com/bketelsen/<project-name>)")
	cmd.Flags().StringVarP(&flagOutput, "output", "o", "", "Output directory (default: ./<project-name>)")
	cmd.Flags().StringVar(&flagDesc, "desc", "A toolbox-powered Go project", "One-line project description")
	cmd.Flags().StringVar(&flagToolboxPath, "toolbox-path", "", "Path to local toolbox clone for replace directive (auto-detected if omitted)")
	cmd.Flags().BoolVar(&flagInteractive, "interactive", false, "Use interactive prompts to fill in details")

	return cmd
}

func runScaffold(cmd *cobra.Command, args []string) error {
	name := args[0]
	log := slog.With("project", name)

	// Apply defaults.
	module := flagModule
	if module == "" {
		module = "github.com/bketelsen/" + name
	}
	outDir := flagOutput
	if outDir == "" {
		outDir = filepath.Join(".", name)
	}
	description := flagDesc

	// Resolve toolbox path for the replace directive.
	toolboxPath := flagToolboxPath
	if toolboxPath == "" {
		var err error
		toolboxPath, err = detectToolboxPath()
		if err != nil {
			log.Debug("could not auto-detect toolbox path, falling back to GOPATH module cache note", "err", err)
			// Use a placeholder the user must update.
			toolboxPath = "/path/to/toolbox"
		}
	}
	// Make toolboxPath absolute so the replace directive is unambiguous.
	if abs, err := filepath.Abs(toolboxPath); err == nil {
		toolboxPath = abs
	}

	// Interactive mode: prompt for overrides.
	if flagInteractive {
		if err := ui.Prompt(
			"Module path",
			"Go module path for go.mod",
			module,
			&module,
		); err != nil {
			return fmt.Errorf("prompt cancelled: %w", err)
		}
		if err := ui.Prompt(
			"Description",
			"One-line project description",
			description,
			&description,
		); err != nil {
			return fmt.Errorf("prompt cancelled: %w", err)
		}
		if err := ui.Prompt(
			"Output directory",
			"Where to create the project",
			outDir,
			&outDir,
		); err != nil {
			return fmt.Errorf("prompt cancelled: %w", err)
		}
		if err := ui.Prompt(
			"Toolbox path",
			"Local path to the toolbox clone (for go.mod replace directive)",
			toolboxPath,
			&toolboxPath,
		); err != nil {
			return fmt.Errorf("prompt cancelled: %w", err)
		}
	}

	data := templateData{
		Name:        name,
		Module:      module,
		Description: description,
		ToolboxPath: toolboxPath,
	}

	log.Debug("scaffold parameters",
		"module", module,
		"output", outDir,
		"description", description,
		"toolbox_path", toolboxPath,
		"dry_run", toolbox.DryRun,
	)

	r := toolbox.NewReporter()

	// Step 1: Validate
	r.Step(1, 3, "Validate")
	log.Info("validating scaffold parameters")

	if !toolbox.DryRun {
		if _, err := os.Stat(outDir); err == nil {
			return fmt.Errorf("output directory %q already exists; remove it or choose a different --output path", outDir)
		}
	}
	log.Debug("output directory available", "path", outDir)

	// Build the file manifest.
	type fileSpec struct {
		RelPath  string
		Template string
	}
	files := []fileSpec{
		{"go.mod", tmplGoMod},
		{"go.sum", tmplGoSum},
		{"main.go", tmplMainGo},
		{"cmd/root.go", tmplCmdRootGo},
		{"cmd/root_test.go", tmplCmdRootTestGo},
		{"cmd/greet.go", tmplCmdGreetGo},
		{"cmd/process.go", tmplCmdProcessGo},
		{"Makefile", tmplMakefile},
		{".gitignore", tmplGitignore},
		{"README.md", tmplReadme},
	}

	// Step 2: Generate files
	r.Step(2, 3, "Generate files")

	var generated []generatedFile

	for _, f := range files {
		absPath := filepath.Join(outDir, filepath.FromSlash(f.RelPath))

		if toolbox.DryRun {
			r.Message("would write %s", absPath)
			log.Debug("dry run: skipping write", "file", f.RelPath)
			generated = append(generated, generatedFile{Path: absPath, Written: false})
			continue
		}

		log.Debug("generating file", "file", f.RelPath)

		// Ensure parent directory exists.
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			return fmt.Errorf("create directory for %s: %w", f.RelPath, err)
		}

		tmpl, err := template.New(f.RelPath).Parse(f.Template)
		if err != nil {
			return fmt.Errorf("parse template for %s: %w", f.RelPath, err)
		}

		fh, err := os.Create(absPath)
		if err != nil {
			return fmt.Errorf("create file %s: %w", absPath, err)
		}

		if err := tmpl.Execute(fh, data); err != nil {
			_ = fh.Close()
			return fmt.Errorf("render template for %s: %w", f.RelPath, err)
		}
		if err := fh.Close(); err != nil {
			return fmt.Errorf("close %s: %w", f.RelPath, err)
		}

		r.Message("wrote %s", f.RelPath)
		generated = append(generated, generatedFile{Path: absPath, Written: true})
	}

	// Step 3: Report completion
	r.Step(3, 3, "Complete")

	type scaffoldResult struct {
		Name        string          `json:"name"`
		Module      string          `json:"module"`
		OutputDir   string          `json:"output_dir"`
		Description string          `json:"description"`
		ToolboxPath string          `json:"toolbox_path"`
		DryRun      bool            `json:"dry_run"`
		Files       []generatedFile `json:"files"`
	}

	result := scaffoldResult{
		Name:        name,
		Module:      module,
		OutputDir:   outDir,
		Description: description,
		ToolboxPath: toolboxPath,
		DryRun:      toolbox.DryRun,
		Files:       generated,
	}

	if toolbox.OutputJSON(result) {
		return nil
	}

	if toolbox.DryRun {
		r.Complete("Dry run complete — no files written", result)
		return nil
	}

	r.Complete(fmt.Sprintf("Project %q created at %s", name, outDir), result)
	printNextSteps(outDir)

	return nil
}

// detectToolboxPath walks up from the running executable to find the
// directory containing go.mod with "module github.com/bketelsen/toolbox".
func detectToolboxPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	for {
		candidate := filepath.Join(dir, "go.mod")
		if isToolboxGoMod(candidate) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("toolbox go.mod not found walking up from %s", filepath.Dir(exe))
}

// isToolboxGoMod reports whether the file at path is the toolbox go.mod.
func isToolboxGoMod(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close() //nolint:errcheck // best-effort cleanup
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "module github.com/bketelsen/toolbox" {
			return true
		}
	}
	return false
}

// printNextSteps prints the post-scaffold instructions to stderr via ui.
func printNextSteps(outDir string) {
	steps := []string{
		fmt.Sprintf("cd %s", outDir),
		"go mod tidy",
		"go run .",
	}
	ui.Info("Next steps:")
	for i, s := range steps {
		ui.Info(fmt.Sprintf("  %d. %s", i+1, s))
	}
}
