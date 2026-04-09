# toolbox — AI Agent Guide

> Companion to `/workspace/CLAUDE.md`. Read that file first; this one covers patterns and context not already there.

## What this library is

`github.com/bketelsen/toolbox` is a **Go library** consumed by CLI tool authors. It is not a standalone application. The primary contract is the public API surface — breaking changes to exported symbols require a semver bump.

## Package map

| Package | Purpose | Key exports |
|---------|---------|-------------|
| `.` (root) | App bootstrap + flag globals | `App`, `Run()`, `JSONOutput`, `Verbose`, `DryRun`, `Silent`, `OutputJSON`, `OutputJSONError`, `NewReporter`, `ExpandPath` |
| `reporter/` | Structured progress reporting | `Reporter` interface, `TextReporter`, `JSONReporter`, `NoopReporter`, `ProgressEvent`, `EventType` constants |
| `ui/` | Console output + interactive prompts | `Console` (global), `ConsolePrinter`, `DisplayTable`, `Table`, `Prompt`, `Confirm`, `Option`, `Spinner` |
| `slug/` | `slog.Handler` (forked from tint) — **not yet present in codebase** | Structured log formatting for terminal output |

## App bootstrap pattern

```go
var app = toolbox.App{
    Version: version, Commit: commit,
    Date: date,       BuiltBy: builtBy,
}

func main() { app.Run(rootCmd) }
```

`Run()` injects persistent flags (`--json`, `--verbose`, `--dry-run`, `--silent`) onto the root Cobra command, then calls `fang.Execute`. Zero-value fields get defaults. Build-time variables are injected via `-ldflags`.

## Reporter interface

`Reporter` is the standard way to emit progress in long-running commands:

```go
type Reporter interface {
    Step(step, total int, name string)
    Progress(percent int, message string)
    Message(format string, args ...any)
    Warning(format string, args ...any)
    Error(err error, message string)
    Complete(message string, details any)
    IsJSON() bool
}
```

- `TextReporter` — human-readable to `io.Writer` (usually stderr)
- `JSONReporter` — JSON Lines to `io.Writer` (usually stdout); thread-safe
- `NoopReporter` — discards all output; use in tests to suppress noise
- Obtain via `toolbox.NewReporter()` which picks the right impl based on flag state

## Table rendering

`ui.DisplayTable(data any, sort string, filterColumns []string)` renders a slice of structs. Control columns with struct tags:

```go
type Row struct {
    Name    string `table:"name,default_sort"` // default sort column (required)
    Status  string `table:"status"`
    Hidden  string `table:"-"`                 // excluded
    Nested  Inner  `table:"inner,recursive"`   // flatten nested struct
}
```

Options: `default_sort`, `recursive`, `recursive_inline`, `nosort`, `-`.

## Console output

Use the global `ui.Console` singleton:

```go
ui.Console.Info("message")
ui.Console.Success("done")
ui.Console.Warn("caution")
ui.Console.Error("failed")
```

Prefix variants (`InfoWithPrefix`, etc.) prepend a label. Output goes to stdout/stderr via lipgloss-styled strings.

## JSON output

Commands should respect the `JSONOutput` flag:

```go
if toolbox.OutputJSON(result) {
    return nil // already written
}
// ... human-readable path
```

For errors: `return toolbox.OutputJSONError("context message", err)`.

## Error handling conventions

- Command `RunE` functions return `error`; never `os.Exit` directly inside library code
- Use `Reporter.Error(err, message)` in progress workflows, then return the error
- `OutputJSONError` wraps errors into a `{"error": "..."}` JSON envelope for machine consumers
- No panics; all errors propagate via return values

## Testing conventions

- `testify/assert` for non-fatal assertions, `testify/require` for fatal ones
- Use `NoopReporter` or pass `&bytes.Buffer{}` to `NewTextReporter`/`NewJSONReporter` for isolation
- Table tests use local struct types (`tableTest1`, `tableTest2`) with various tag combinations
- Avoid testing private functions; test behavior through exported API

## Naming conventions

- Exported types: `PascalCase`
- Flag variables: package-level `var` with `PascalCase` (`JSONOutput`, `DryRun`)
- Constructors: `New<Type>` (e.g., `NewTextReporter`, `NewReporter`)
- Struct tag key: `table` (lowercase)

## Examples

`_examples/` contains runnable demo CLIs — not part of the library's import path. They demonstrate integration patterns:

| Example | Demonstrates |
|---------|-------------|
| `greet` | `App.Run()`, flags, slog handler |
| `dashboard` | `DisplayTable`, console messages, `--interactive` flag |
| `deploy` | Multi-step `Reporter` workflow (Step/Message/Warning/Complete) |
| `migration` | Batch progress reporting with `Progress()` percent |
| `fileprocess` | File-based workflows with progress reporting |
| `healthcheck` | Health check patterns with status output |
| `sync` | Multi-subcommand CLI with Viper config + `BindViper()` |

When adding a new library feature, add or update an example to show usage.

## Common pitfalls

- **Do not call `BindViper()` before `Run()`** — flags must be registered first.
- **`DisplayTable` requires exactly one `default_sort` (or at least one `nosort`) tag** — omitting both returns an error.
- **`JSONReporter` writes to stdout; `TextReporter` to stderr** — keep this distinction when constructing reporters manually.
- **`App` fields are strings, not typed** — Version/Commit/Date are injected as raw strings via ldflags.
