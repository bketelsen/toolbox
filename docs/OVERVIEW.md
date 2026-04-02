# toolbox

## Purpose

`github.com/bketelsen/toolbox` is a Go library providing a standardized toolkit for building CLI tools.
It wraps Cobra, Charmbracelet (bubbletea / huh / lipgloss / fang), and Viper with opinionated defaults
for version injection, common flags, output formatting, progress reporting, and interactive UI components.
The library is pre-1.0; its public API is subject to change.

## Architecture

```
github.com/bketelsen/toolbox   ← root package (main API surface)
├── clix.go          App struct — build-time metadata + Run() entry point (fang integration)
├── flags.go         Package-level flag globals (JSONOutput, Verbose, DryRun, Silent) + BindViper()
├── output.go        OutputJSON / OutputJSONError — flag-conditioned JSON helpers
├── reporter.go      NewReporter() factory — returns Text/JSON/Noop based on flag state
├── path.go          ExpandPath() — resolve ~ and relative paths
├── trace.go         Tracef() — stderr trace logging gated on TOOLBOX_TRACING=on
│
├── reporter/        Reporter interface + TextReporter / JSONReporter / NoopReporter
├── ui/              Console styling, table rendering, spinners, interactive prompts (huh)
├── slug/            Custom slog.Handler (forked from tint) for structured logging
├── s/               String case helpers (thin wrapper around iancoleman/strcase)
├── scaffold/        Scaffolding utilities
├── cmd/scaffold/    Cobra sub-command for the scaffold tool
├── dashboard/       Bubbletea dashboard component
│
└── _examples/       Demo CLI apps (greet, dashboard, deploy, …) — not imported by library
```

Data flow for a typical CLI command:

1. `App.Run()` — registers flags via `flags.go`, hands off to fang/Cobra
2. Command handler reads `toolbox.JSONOutput` / `toolbox.Verbose` etc.
3. Output helpers (`OutputJSON`, `reporter.Reporter`) format and write results
4. `ui/` functions render interactive or styled terminal output

## Key Patterns

### App bootstrapping
`App` (defined in `clix.go`) carries `Name`, `Version`, `Commit`, `Date` set at link-time via
`-ldflags`. Call `App.Run(rootCmd)` as the entry point; fang handles shell completion, man-page
generation, and environment-variable binding automatically.

### Global flags
Four boolean flags are registered as persistent Cobra flags in `flags.go`:

| Flag | Variable | Purpose |
|------|----------|---------|
| `--json` | `JSONOutput` | Emit machine-readable JSON |
| `--verbose` | `Verbose` | Enable verbose logging |
| `--dry-run` | `DryRun` | Simulate without side effects |
| `--silent` | `Silent` | Suppress all non-error output |

`BindViper()` connects these to Viper for config-file / env-var override.

### Reporter pattern
`NewReporter()` returns a `reporter.Reporter` based on the active flags:
- `JSONOutput=true` → `JSONReporter` (stdout, JSON Lines with timestamps)
- `Silent=true` → `NoopReporter`
- default → `TextReporter` (stderr, human-readable)

### UI helpers (`ui/`)
- `console.go` — `Info`, `Warn`, `Error`, `Success` styled messages
- `table.go` — struct-tag–driven table rendering (`table:"field,options"`)
- `spinner.go` — bubbletea spinner wrapper
- `confirm.go`, `prompt.go`, `option.go` — interactive prompts via huh

### Structured logging (`slug/`)
A custom `slog.Handler` forked from [tint](https://github.com/lmittmann/tint).
See `slug/README.md` for handler construction options.

### String utilities (`s/`)
Thin wrappers around `iancoleman/strcase` — camelCase, snake_case, kebab-case conversions.
No additional logic beyond delegation.

### Tracing
`Tracef(format, args...)` writes caller-annotated lines to stderr when `TOOLBOX_TRACING=on`.
Use `SetTracing(bool)` to toggle programmatically.

## Configuration

| Mechanism | Details |
|-----------|---------|
| Build-time version | Inject via `-ldflags "-X main.Version=..."` (see `clix.go`) |
| Config file | Viper-managed; path resolved with `ExpandPath()` |
| Env vars | Any flag key, upper-cased, prefixed per Viper config |
| Tracing | `TOOLBOX_TRACING=on` enables `Tracef()` output |

## Development Commands

```bash
go build ./...          # build all packages
go test ./...           # run full test suite
go vet ./...            # static analysis
gofmt -w -s .           # format
go mod tidy             # clean module graph

make check              # fmt + lint + test
task checks             # staticcheck + vet + test + format + tidy
make test-cover         # generate coverage.html
make bump               # version bump (requires clean tree + svu)
```

CI runs on push to `main` and on PRs via GitHub Actions:
- **Go test** (`cobra.yml`): `go test ./...` with Go 1.25
- **Lint** (`lint.yml`): `gofmt -s -d .` and `go vet ./...`

## Detailed Documentation

- [slug/README.md](../slug/README.md) — slog handler construction and options
- [docs/plans/](plans/) — implementation plans for significant changes

## Change History (notable)

| Version | Change |
|---------|--------|
| pre-1.0 | `n/` package (`FreePort()`) removed — had zero internal callers; external consumers must remove the import |
| pre-1.0 | `s/` package retained as thin strcase wrapper |
