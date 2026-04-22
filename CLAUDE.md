# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`github.com/bketelsen/toolbox` is a **Go library** (not an application) providing a standardized toolkit for building CLI tools. It wraps Cobra, Charmbracelet (bubbletea/huh/lipgloss/fang), and Viper with opinionated defaults for version injection, common flags, output formatting, progress reporting, and interactive UI components.

## Commands

```bash
# Build
go build ./...

# Test
go test -v ./...          # or: make test / task test

# Format
gofmt -w -s .             # or: make fmt / task format

# Lint / Vet
go vet ./...              # or: task vet
staticcheck ./...         # or: task staticcheck
golangci-lint run         # optional, may not be installed

# Full check suite
make check                # fmt + lint + test
task checks               # staticcheck + vet + test + format + tidy

# Coverage
make test-cover           # generates coverage.html

# Version bump (requires clean tree + svu)
make bump
```

## Architecture

The root package (`toolbox`) provides the main API surface:

- **`App`** (`clix.go`) - Build-time metadata + `Run()` method that registers flags and executes via fang
- **Flags** (`flags.go`) - Package-level globals (`JSONOutput`, `Verbose`, `DryRun`, `Silent`) auto-registered as persistent Cobra flags; `BindViper()` for config integration
- **Output** (`output.go`) - `OutputJSON`/`OutputJSONError` helpers conditioned on the `JSONOutput` flag
- **Reporter factory** (`reporter.go`) - `NewReporter()` returns Text/JSON/Noop reporter based on flag state
- **Path** (`path.go`) - `ExpandPath()` for resolving `~` and relative paths

Key subpackages:

- **`reporter/`** - `Reporter` interface with `TextReporter` (stderr), `JSONReporter` (stdout, JSON Lines with timestamps), `NoopReporter` (silent)
- **`ui/`** - Console styling (`Info`/`Warn`/`Error`/`Success`), table rendering with struct tags (`table:"field,options"`), spinners, interactive prompts (`Confirm`/`Prompt`/`Option` via huh)
- **Logging** (`logger.go`) - `setupLogger()` configures `a.Logger` using stdlib `log/slog` with JSON or text handler based on `JSONOutput`; log level follows `Verbose`

The `_examples/` directory contains demo CLI apps showing usage patterns (greet, dashboard, deploy, etc.).

## CI

GitHub Actions runs on push to `main` and on PRs:
- **Go test** (`cobra.yml`): `go test ./...` with Go 1.25
- **Lint** (`lint.yml`): `gofmt -s -d .` and `go vet ./...`

## Testing Notes

- Uses standard `testing` package with `github.com/stretchr/testify` for assertions
- Run a single test: `go test -v -run TestName ./path/to/package`
