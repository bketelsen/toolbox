# Toolbox

A standardized toolkit for building CLI tools in Go with opinionated defaults for configuration, output formatting, and interactive UI components.

## Features

- **Unified CLI pattern**: Wrap Cobra commands with a simple `App` struct and `Run()` method
- **Common flags**: Built-in `--json`, `--verbose`, `--dry-run`, and `--silent` flags
- **Reporter pattern**: Switch between text, JSON, and silent output modes automatically
- **Charmbracelet integration**: Ready-to-use styling, tables, spinners, and interactive prompts
- **Version injection**: Automatic version, commit, and build date handling
- **Viper config support**: Easy integration with configuration files

## Installation

### Using `go install`

```bash
go install github.com/bketelsen/toolbox/cmd/scaffold@latest
```

### Clone and Build

```bash
git clone https://github.com/bketelsen/toolbox
cd toolbox
go build ./...
```

## Quick Start

Here's a minimal CLI application using toolbox:

```go
package main

import (
	"fmt"
	"os"

	"github.com/bketelsen/toolbox"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "greet [name]",
		Short: "Greet someone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			greeting := fmt.Sprintf("Hello, %s!", name)

			// Automatically outputs JSON if --json flag is set
			if toolbox.OutputJSON(map[string]string{"greeting": greeting}) {
				return nil
			}

			fmt.Println(greeting)
			return nil
		},
	}

	// Create an app with version metadata and run
	app := toolbox.App{
		Version: "1.0.0",
		Commit:  "abc123",
		Date:    "2026-04-07",
		BuiltBy: "ci",
	}

	if err := app.Run(rootCmd); err != nil {
		os.Exit(1)
	}
}
```

Run it:
```bash
./greet Alice
./greet Alice --json
./greet Alice --verbose
```

See more examples in the [`_examples/`](./_examples/) directory.

## Global Flags

The following flags are automatically registered on all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--json` | — | `false` | Output in JSON format |
| `--verbose` | `-v` | `false` | Enable verbose output (often used with structured logging) |
| `--dry-run` | `-n` | `false` | Perform a dry run without making changes |
| `--silent` | `-s` | `false` | Suppress all progress output (takes precedence over `--json`) |

Access these flags from the `toolbox` package:
```go
if toolbox.JSONOutput {
    // Output JSON
}
if toolbox.Verbose {
    // Enable detailed logging
}
if toolbox.DryRun {
    // Skip write operations
}
if toolbox.Silent {
    // Suppress output
}
```

## Reporter Pattern

The reporter pattern provides a unified way to handle output across text and JSON modes. Use `NewReporter()` to get the appropriate reporter based on flags:

```go
reporter := toolbox.NewReporter()

// Text mode writes to stderr with styling
// JSON mode writes JSON Lines to stdout
// Silent mode produces no output
reporter.Info("Processing started")
reporter.Success("Task completed")
reporter.Error("Something went wrong")
```

The reporter automatically:
- Routes text output to **stderr** (keeps stdout clean for data piping)
- Routes JSON output to **stdout** (for easy parsing and piping)
- Suppresses all output in silent mode

### Example with Structured Data

```go
func processFile(path string) error {
	reporter := toolbox.NewReporter()

	reporter.Info("Starting file processing", map[string]any{
		"path": path,
		"dryRun": toolbox.DryRun,
	})

	if toolbox.DryRun {
		reporter.Info("Dry run mode - no changes made")
		return nil
	}

	// Process file...

	reporter.Success("File processed", map[string]any{
		"path": path,
		"size": "1.2MB",
	})
	return nil
}
```

## UI Package

The `ui` package provides styled console output, tables, spinners, and interactive prompts:

### Console Styling

```go
import "github.com/bketelsen/toolbox/ui"

ui.Console.Info("Information message")
ui.Console.Warn("Warning message")
ui.Console.Error("Error message")
ui.Console.Success("Success message")
```

### Tables

Render data with automatic column formatting using struct tags:

```go
type User struct {
	Name  string `table:"name,width:20"`
	Email string `table:"email,width:30"`
	Role  string `table:"role"`
}

users := []User{
	{Name: "Alice", Email: "alice@example.com", Role: "admin"},
	{Name: "Bob", Email: "bob@example.com", Role: "user"},
}

ui.Console.RenderTable(users)
```

### Spinners

Show progress with spinners:

```go
spinner := ui.NewSpinner("Processing...")
spinner.Start()
time.Sleep(2 * time.Second)
spinner.Stop("Done!")
```

### Interactive Prompts

Gather user input with styled prompts:

```go
// Yes/No confirmation
ok, err := ui.Console.Confirm("Continue?")

// Text input
name, err := ui.Console.Prompt("Enter your name")

// Multiple choice
choice, err := ui.Console.Option("Choose an option", []string{"Option A", "Option B", "Option C"})
```

## Command Examples

Explore complete working examples in the [`_examples/`](./_examples/) directory:

- **`greet`**: Basic greeting with flags and JSON output
- **`dashboard`**: Terminal UI dashboard example
- **`deploy`**: Deployment workflow with spinners
- **`sync`**: File synchronization with progress reporting
- **`migration`**: Database migration runner
- **`healthcheck`**: Service health monitoring
- **`fileprocess`**: Batch file processing with reporters

## Scaffolding

Use the `scaffold` command to generate a new CLI project:

```bash
go run ./cmd/scaffold [project-name]
```

This generates a ready-to-use project structure with:
- Pre-configured `App` struct
- Sample commands with proper flag handling
- Error handling patterns
- Testing setup

## Module Path

```
github.com/bketelsen/toolbox
```

## Requirements

- Go 1.25 or later
- Charmbracelet libraries: bubbles, bubbletea, fang, huh, lipgloss
- spf13/cobra for command handling
- spf13/viper for configuration

## License

MIT
