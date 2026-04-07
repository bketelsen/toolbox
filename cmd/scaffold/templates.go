package main

// templates holds all file templates for the generated project.
// Each template receives a templateData struct as its dot value.

var tmplGoMod = `module {{.Module}}

go 1.25.0

require (
	github.com/bketelsen/toolbox v0.0.0-dev
	github.com/spf13/cobra v1.10.2
)

replace github.com/bketelsen/toolbox => {{.ToolboxPath}}
`

var tmplGoSum = ``

var tmplMainGo = `package main

import (
	"os"

	"github.com/bketelsen/toolbox"
	"{{.Module}}/cmd"
)

func main() {
	app := toolbox.App{
		Version: "0.1.0",
		Commit:  "none",
		Date:    "unknown",
		BuiltBy: "source",
	}
	if err := app.Run(cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
`

var tmplCmdRootGo = `package cmd

import (
	"log/slog"
	"os"

	"github.com/bketelsen/toolbox"
	"github.com/spf13/cobra"
)

// RootCmd is the root cobra command for {{.Name}}.
var RootCmd = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.Description}}",
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
}

func init() {
	RootCmd.AddCommand(greetCmd)
	RootCmd.AddCommand(processCmd)
	RootCmd.AddCommand(completionCmd)
}
`

var tmplCmdRootTestGo = `package cmd

import "testing"

func TestRootCmdHasSubcommands(t *testing.T) {
	if len(RootCmd.Commands()) < 3 {
		t.Errorf("expected at least 3 subcommands, got %d", len(RootCmd.Commands()))
	}
}

func TestGreetCmdRegistered(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "greet" {
			found = true
			break
		}
	}
	if !found {
		t.Error("greet subcommand not registered on RootCmd")
	}
}

func TestProcessCmdRegistered(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "process" {
			found = true
			break
		}
	}
	if !found {
		t.Error("process subcommand not registered on RootCmd")
	}
}
`

var tmplCmdGreetGo = `package cmd

import (
	"fmt"
	"log/slog"

	"github.com/bketelsen/toolbox"
	"github.com/spf13/cobra"
)

type greetResult struct {
	Name    string ` + "`" + `json:"name"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
	DryRun  bool   ` + "`" + `json:"dry_run"` + "`" + `
}

var greetCmd = &cobra.Command{
	Use:   "greet [name]",
	Short: "Print a greeting",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		slog.Debug("preparing greeting", "name", name)

		result := greetResult{
			Name:    name,
			Message: fmt.Sprintf("Hello, %s!", name),
			DryRun:  toolbox.DryRun,
		}

		if toolbox.DryRun {
			slog.Info("dry run: skipping output", "name", name)
			if toolbox.OutputJSON(result) {
				return nil
			}
			fmt.Printf("(dry run) would greet: %s\n", name)
			return nil
		}

		if toolbox.OutputJSON(result) {
			return nil
		}

		fmt.Println(result.Message)
		slog.Info("greeting sent", "name", name)
		return nil
	},
}
`

var tmplCmdProcessGo = `package cmd

import (
	"fmt"
	"log/slog"

	"github.com/bketelsen/toolbox"
	"github.com/spf13/cobra"
)

type processResult struct {
	Target string ` + "`" + `json:"target"` + "`" + `
	Status string ` + "`" + `json:"status"` + "`" + `
	DryRun bool   ` + "`" + `json:"dry_run"` + "`" + `
}

var processCmd = &cobra.Command{
	Use:   "process [target]",
	Short: "Run a multi-step workflow on target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		log := slog.With("target", target)
		r := toolbox.NewReporter()

		// Step 1: Validate
		r.Step(1, 3, "Validate")
		log.Info("validating target", "target", target)
		log.Debug("validation detail", "length", len(target))
		if target == "" {
			return fmt.Errorf("target must not be empty")
		}
		r.Message("target %q is valid", target)

		// Step 2: Process (skipped in dry-run)
		r.Step(2, 3, "Process")
		if toolbox.DryRun {
			r.Warning("dry run: skipping process step for %q", target)
			log.Info("dry run mode, skipping process")
		} else {
			r.Progress(50, "processing...")
			log.Info("processing target", "target", target)
			r.Message("processed %q successfully", target)
			r.Progress(100, "done")
		}

		// Step 3: Complete
		r.Step(3, 3, "Complete")
		log.Debug("workflow completing", "target", target)

		result := processResult{
			Target: target,
			Status: "ok",
			DryRun: toolbox.DryRun,
		}

		if toolbox.OutputJSON(result) {
			return nil
		}

		r.Complete(fmt.Sprintf("Workflow complete for %q", target), result)
		return nil
	},
}
`

var tmplCmdCompletionGo = `package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenBashCompletion(os.Stdout)
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenZshCompletion(os.Stdout)
	},
}

var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	},
}

var completionPowerShellCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Generate PowerShell completion script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenPowerShellCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionFishCmd)
	completionCmd.AddCommand(completionPowerShellCmd)
}
`

var tmplMakefile = `BINARY := {{.Name}}
DIST    := dist

.PHONY: fmt vet lint test build tidy check

fmt:
	gofmt -w -s .

vet:
	go vet ./...

lint: vet

test:
	go test -v ./...

build:
	go build -o $(DIST)/$(BINARY) .

tidy:
	go mod tidy

check: fmt vet test
`

var tmplGitignore = `# Binaries
dist/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage tool
*.out
coverage.html

# Dependency directories
vendor/

# Environment files
.env
.env.*

# Editor files
.idea/
.vscode/
*.swp
*.swo
`

var tmplReadme = `# {{.Name}}

{{.Description}}

## Install

` + "```" + `sh
go install {{.Module}}@latest
` + "```" + `

Or clone and build locally:

` + "```" + `sh
git clone <your-repo-url>
cd {{.Name}}
make build
` + "```" + `

## Usage

` + "```" + `
{{.Name}} [command] [flags]
` + "```" + `

### Commands

#### greet

Print a greeting. Demonstrates basic positional args, dry-run, and JSON output.

` + "```" + `sh
# Basic usage
{{.Name}} greet Alice

# Dry run (skip output)
{{.Name}} greet Alice --dry-run

# JSON output
{{.Name}} greet Alice --json
` + "```" + `

#### process

Run a multi-step workflow. Demonstrates reporter steps, progress, warnings,
dry-run, verbose logging, and JSON output.

` + "```" + `sh
# Basic usage
{{.Name}} process myfile.txt

# Dry run (skip process step)
{{.Name}} process myfile.txt --dry-run

# JSON output
{{.Name}} process myfile.txt --json

# Verbose logging
{{.Name}} process myfile.txt --verbose

# Silent (suppress progress output)
{{.Name}} process myfile.txt --silent
` + "```" + `

## Global Flags

| Flag        | Short | Default | Description                          |
|-------------|-------|---------|--------------------------------------|
| ` + "`--json`" + `    |       | false   | Output results as JSON               |
| ` + "`--verbose`" + ` | ` + "`-v`" + `   | false   | Enable verbose/debug logging         |
| ` + "`--dry-run`" + ` | ` + "`-n`" + `   | false   | Simulate actions without executing   |
| ` + "`--silent`" + `  | ` + "`-s`" + `   | false   | Suppress all progress output         |

## Shell Completions

` + "`{{.Name}}`" + ` supports shell completion for bash, zsh, fish, and PowerShell.

` + "```" + `sh
# Bash
source <({{.Name}} completion bash)

# Zsh
source <({{.Name}} completion zsh)

# Fish
{{.Name}} completion fish | source

# PowerShell
{{.Name}} completion powershell | Out-String | Invoke-Expression
` + "```" + `

## Development

` + "```" + `sh
make fmt    # Format code
make vet    # Run go vet
make test   # Run tests
make build  # Build to dist/{{.Name}}
make tidy   # go mod tidy
make check  # fmt + vet + test
` + "```" + `
`
