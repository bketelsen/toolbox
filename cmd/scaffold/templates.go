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
	app.Run(cmd.RootCmd)
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
	RootCmd.AddCommand(exampleCmd)
}
`

var tmplCmdExampleGo = `package cmd

import (
	"log/slog"

	"github.com/bketelsen/toolbox"
	"github.com/spf13/cobra"
)

type exampleResult struct {
	Name    string ` + "`" + `json:"name"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
	DryRun  bool   ` + "`" + `json:"dry_run"` + "`" + `
}

var exampleCmd = &cobra.Command{
	Use:   "example [name]",
	Short: "Run an example workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		log := slog.With("name", name)
		r := toolbox.NewReporter()

		// Step 1: Validate
		r.Step(1, 3, "Validate")
		log.Info("validating input", "name", name)
		log.Debug("validation detail", "length", len(name))

		// Step 2: Process (skipped in dry-run)
		r.Step(2, 3, "Process")
		if toolbox.DryRun {
			r.Warning("dry run: skipping process step")
			log.Info("dry run mode, skipping process")
		} else {
			log.Info("processing", "name", name)
			r.Message("processed %q", name)
		}

		// Step 3: Complete
		r.Step(3, 3, "Complete")
		log.Debug("completing workflow")

		result := exampleResult{
			Name:    name,
			Message: "Hello, " + name + "!",
			DryRun:  toolbox.DryRun,
		}

		if toolbox.OutputJSON(result) {
			return nil
		}

		r.Complete("Workflow complete", result)
		return nil
	},
}
`

var tmplMakefile = `BINARY := {{.Name}}
DIST    := dist

.PHONY: fmt lint test build tidy check

fmt:
	gofmt -w -s .

lint:
	go vet ./...

test:
	go test -v ./...

build:
	go build -o $(DIST)/$(BINARY) .

tidy:
	go mod tidy

check: fmt lint test
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

#### example

Run an example workflow demonstrating reporter steps, dry-run support, and JSON output.

` + "```" + `sh
# Basic usage
{{.Name}} example myinput

# Dry run (skip process step)
{{.Name}} example myinput --dry-run

# JSON output
{{.Name}} example myinput --json
` + "```" + `

## Global Flags

| Flag        | Short | Default | Description                          |
|-------------|-------|---------|--------------------------------------|
| ` + "`--json`" + `    |       | false   | Output results as JSON               |
| ` + "`--verbose`" + ` | ` + "`-v`" + `   | false   | Enable verbose/debug logging         |
| ` + "`--dry-run`" + ` | ` + "`-n`" + `   | false   | Simulate actions without executing   |
| ` + "`--silent`" + `  | ` + "`-s`" + `   | false   | Suppress all progress output         |

## Development

` + "```" + `sh
make fmt    # Format code
make lint   # Run go vet
make test   # Run tests
make build  # Build to dist/{{.Name}}
make tidy   # go mod tidy
make check  # fmt + lint + test
` + "```" + `
`
