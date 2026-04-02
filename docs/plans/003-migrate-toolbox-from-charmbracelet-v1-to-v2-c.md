---
devTaskId: 003
title: Migrate toolbox from Charmbracelet v1 to v2 (charm.land vanity domain)
status: Approved
created: 2026-04-02
author: architect
adrNeeded: false
designOnly: false
complexity: moderate
---

# Plan: Migrate toolbox from Charmbracelet v1 to v2 (charm.land vanity domain)

## Summary

Replace all five Charmbracelet v1 library imports (`bubbletea`, `bubbles`, `lipgloss`, `huh`, `fang`)
with their v2 equivalents on the `charm.land/` vanity domain, update call sites to match the v2 API,
and remove the direct `muesli/termenv` dependency by replacing its usage with the already-transitively-
available `github.com/charmbracelet/colorprofile` and `github.com/charmbracelet/x/ansi` packages.
The migration must be completed atomically (all v1 → v2 at once; no v1/v2 mixing).

## Current State

### go.mod direct charm dependencies (v1)
| Package | Version |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | v1.3.10 |
| `github.com/charmbracelet/bubbles` | v1.0.0 |
| `github.com/charmbracelet/lipgloss` | v1.1.0 |
| `github.com/charmbracelet/huh` | v1.0.0 |
| `github.com/charmbracelet/fang` | v1.0.0 |
| `github.com/muesli/termenv` | v0.16.0 (direct) |

`charm.land/lipgloss/v2` already appears as an indirect dep (pre-release beta) — it will be promoted.

### Files with charm/termenv imports
| File | Packages imported |
|------|------------------|
| `ui/spinner.go` | `bubbles/spinner`, `bubbletea`, `lipgloss`, `termenv` |
| `ui/styles.go` | `lipgloss`, `termenv` |
| `ui/message.go` | `lipgloss` |
| `ui/confirm.go` | `huh` |
| `ui/option.go` | `huh` |
| `ui/prompt.go` | `huh` |
| `clix.go` | `fang` |

`_examples/` does **not** import charm packages directly (uses toolbox and toolbox/ui only).

### Latent bug in `ui/styles.go`
`var color termenv.Profile` is declared but never assigned; it stays at the zero value `termenv.Ascii`.
This means `isTerm()` always returns `false`, so `Bold()`, `Red()`, `Yellow()`, and `Green()` silently
return plain, unstyled strings. The migration will fix this by replacing the pattern with
`colorprofile.Detect()` at package init time.

## Target State

| Package | New import path | Version |
|---------|----------------|---------|
| bubbletea | `charm.land/bubbletea/v2` | v2.0.2 |
| bubbles | `charm.land/bubbles/v2` | v2.1.0 |
| lipgloss | `charm.land/lipgloss/v2` | v2.0.2 |
| huh | `charm.land/huh/v2` | v2.0.3 |
| fang | `charm.land/fang/v2` | v2.0.1 |
| termenv | **removed** (replaced by `colorprofile` + `x/ansi`) | — |

## Architecture

No structural changes to the library. This is a pure dependency-upgrade + call-site migration.

### API changes by file

#### `clix.go` — LOW impact
- Change `"github.com/charmbracelet/fang"` → `"charm.land/fang/v2"`
- No API changes: `fang.Execute`, `fang.WithVersion`, `fang.WithNotifySignal` are identical in v2.

#### `ui/confirm.go`, `ui/option.go`, `ui/prompt.go` — LOW impact
- Change `"github.com/charmbracelet/huh"` → `"charm.land/huh/v2"`
- Builder API (`NewConfirm`, `NewSelect`, `NewInput`, etc.) is unchanged.

#### `ui/message.go` — LOW impact
- Change `"github.com/charmbracelet/lipgloss"` → `"charm.land/lipgloss/v2"`
- No API changes: `lipgloss.Style`, `lipgloss.NewStyle()` signature is identical.

#### `ui/styles.go` — MEDIUM impact
Three changes required:

1. **Import swap**: `"github.com/charmbracelet/lipgloss"` → `"charm.land/lipgloss/v2"`;
   remove `"github.com/muesli/termenv"`; add `"github.com/charmbracelet/colorprofile"` and `"os"`.

2. **`lipgloss.Color()` return type**: In v2, `lipgloss.Color("1")` returns `color.Color`
   (the `image/color` interface) instead of `lipgloss.Color` (was a string type). The call
   syntax is unchanged; the variable declarations (`red`, `green`, etc.) need their type
   annotation updated from `lipgloss.Color` to `color.Color`. Because they're package-level
   `var`s (untyped in current code), this will compile automatically — no explicit annotation
   is needed if the vars are written as `var red = lipgloss.Color("1")`.

3. **Replace `termenv.Profile` / `isTerm()`**: Remove `var color termenv.Profile`. Add:
   ```go
   import (
       "github.com/charmbracelet/colorprofile"
       "os"
   )

   var isTTY bool

   func init() {
       isTTY = colorprofile.Detect(os.Stdout, os.Environ()) != colorprofile.NoTTY
   }

   func isTerm() bool { return isTTY }
   ```
   This fixes the latent bug (was always false) and replaces the termenv dependency.

#### `ui/spinner.go` — HIGH impact
Five distinct API changes:

1. **Import paths**: swap all three charm imports; remove `termenv`.
   Add `"github.com/charmbracelet/x/ansi"`.

2. **`View() string` → `View() tea.View`**:
   ```go
   // v1
   func (s *Spinner) View() string {
       return s.spinner.View() + title
   }

   // v2
   func (s *Spinner) View() tea.View {
       return tea.NewView(s.spinner.View() + title)
   }
   ```

3. **`tea.KeyMsg` → `tea.KeyPressMsg`**:
   ```go
   // v1
   case tea.KeyMsg:
       switch msg.String() { case "ctrl+c": return s, tea.Interrupt }

   // v2
   case tea.KeyPressMsg:
       switch msg.String() { case "ctrl+c": return s, tea.Interrupt }
   ```
   Note: `tea.Interrupt` still exists in v2 (confirmed: it's a `func() Msg` command).

4. **`spinner.Tick` → `s.spinner.Tick`**: In bubbles v2, `Tick` is a method on `Model`,
   not a standalone function. Update `Init()`:
   ```go
   // v1
   return tea.Batch(s.spinner.Tick, ...)

   // v2
   return tea.Batch(s.spinner.Tick, ...)  // same call — Tick is now a method, but
   // s.spinner.Tick is already a method expression; this compiles correctly if the
   // existing code uses s.spinner.Tick (it does). No change needed here.
   ```
   Actually `spinner.Tick` (the package-level function in v1) does not exist in v2; only
   `Model.Tick()` exists. In the existing code the call is `s.spinner.Tick` (method value),
   which will work correctly in v2 as `s.spinner.Tick` is still a method value.

5. **`lipgloss.AdaptiveColor` removed in v2**: Replace with `compat.AdaptiveColor` from
   `charm.land/lipgloss/v2/compat`, or simply use a single `lipgloss.Color()` value.
   Because the spinner title color (`#00020A` / `#FFFDF5`) is white-on-dark / black-on-light,
   and we don't need adaptive behavior for a spinner, use a simple neutral color:
   ```go
   // Replace in NewSpinner():
   titleStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")),
   ```
   This is the simplest approach; adaptive color support can be added later with `compat` if needed.

6. **`termenv` cursor management in `runAccessible()`**: Replace `termenv.NewOutput` cursor
   calls with direct ANSI escape writes via `github.com/charmbracelet/x/ansi`:
   ```go
   // v1
   output := termenv.NewOutput(tty)
   output.HideCursor()
   defer output.ShowCursor()

   // v2
   _, _ = fmt.Fprint(tty, ansi.HideCursor)
   defer fmt.Fprint(tty, ansi.ShowCursor)
   ```
   `ansi.HideCursor` and `ansi.ShowCursor` are string constants in `x/ansi` v0.11.6
   (already a transitive dependency).

## Implementation Steps

1. **Update go.mod — add v2 deps, remove v1 and termenv**
   - Files: `go.mod`, `go.sum`
   - Run:
     ```bash
     go get charm.land/bubbletea/v2@v2.0.2
     go get charm.land/bubbles/v2@v2.1.0
     go get charm.land/lipgloss/v2@v2.0.2
     go get charm.land/huh/v2@v2.0.3
     go get charm.land/fang/v2@v2.0.1
     ```
   - Remove old direct deps after all source changes:
     ```bash
     go mod edit -droprequire github.com/charmbracelet/bubbletea
     go mod edit -droprequire github.com/charmbracelet/bubbles
     go mod edit -droprequire github.com/charmbracelet/lipgloss
     go mod edit -droprequire github.com/charmbracelet/huh
     go mod edit -droprequire github.com/charmbracelet/fang
     go mod edit -droprequire github.com/muesli/termenv
     go mod tidy
     ```

2. **Migrate `clix.go`** (LOW)
   - Files: `clix.go` (modify)
   - What: Change import `github.com/charmbracelet/fang` → `charm.land/fang/v2`

3. **Migrate `ui/confirm.go`, `ui/option.go`, `ui/prompt.go`** (LOW)
   - Files: three files (modify)
   - What: Change import `github.com/charmbracelet/huh` → `charm.land/huh/v2` in each

4. **Migrate `ui/message.go`** (LOW)
   - Files: `ui/message.go` (modify)
   - What: Change import `github.com/charmbracelet/lipgloss` → `charm.land/lipgloss/v2`

5. **Migrate `ui/styles.go`** (MEDIUM)
   - Files: `ui/styles.go` (modify)
   - What:
     - Swap lipgloss import path
     - Remove `muesli/termenv` import
     - Add `github.com/charmbracelet/colorprofile` and `os` imports
     - Remove `var color termenv.Profile`
     - Add `var isTTY bool` and an `init()` func using `colorprofile.Detect`
     - Change `isTerm()` to return `isTTY`
     - The ANSI color var declarations (`red`, `green`, etc.) don't need type annotations
       changed; they're inferred from `lipgloss.Color()` return value

6. **Migrate `ui/spinner.go`** (HIGH)
   - Files: `ui/spinner.go` (modify)
   - What:
     - Swap all charm import paths (bubbletea → `charm.land/bubbletea/v2`,
       bubbles/spinner → `charm.land/bubbles/v2/spinner`,
       lipgloss → `charm.land/lipgloss/v2`)
     - Remove `muesli/termenv` import; add `github.com/charmbracelet/x/ansi`
     - Change `View() string` → `View() tea.View`, wrap return in `tea.NewView(...)`
     - Change `case tea.KeyMsg:` → `case tea.KeyPressMsg:`
     - Replace `lipgloss.AdaptiveColor{Light: ..., Dark: ...}` with
       `lipgloss.Color("#FFFDF5")` in `NewSpinner()`
     - Replace `termenv.NewOutput(tty).HideCursor()` with `fmt.Fprint(tty, ansi.HideCursor)`
     - Replace `output.ShowCursor()` / `output.CursorBack(...)` with
       `fmt.Fprint(tty, ansi.ShowCursor)` in defer (drop CursorBack — it used termenv)

7. **Verify and tidy**
   - Run `go build ./...` — must succeed with zero errors
   - Run `go vet ./...` — must be clean
   - Run `go test ./...` — all tests must pass
   - Run `go mod tidy` — removes any remaining orphan indirect deps

## Files to Change

| File | Change type | Summary |
|------|-------------|---------|
| `go.mod` | modify | Add v2 deps, drop v1 + termenv direct deps, run tidy |
| `clix.go` | modify | Import path: fang |
| `ui/confirm.go` | modify | Import path: huh |
| `ui/option.go` | modify | Import path: huh |
| `ui/prompt.go` | modify | Import path: huh |
| `ui/message.go` | modify | Import path: lipgloss |
| `ui/styles.go` | modify | Import path: lipgloss; remove termenv; add colorprofile; fix isTerm() |
| `ui/spinner.go` | modify | Import paths: all charm + termenv; View() return type; KeyPressMsg; AdaptiveColor; cursor management |

No new files are created. `_examples/` requires no changes.

## Test Strategy

After all changes:

```bash
# 1. Build check
go build ./...

# 2. Vet
go vet ./...

# 3. Unit tests
go test -v ./...

# 4. Module graph clean
go mod tidy
git diff go.mod go.sum   # confirm v1 deps removed, v2 deps present

# 5. Smoke test spinner (manual, if desired)
go run ./_examples/dashboard
```

Existing tests cover:
- `ui/table_test.go` — lipgloss style usage
- `clix_test.go`, `flags_test.go`, `output_test.go` — toolbox root package
- `reporter/` tests — no charm deps

There are no spinner/styles unit tests today; visual correctness is verified by the dashboard example.

## ADR Required

needed: no  
(This is a mechanical dependency upgrade following published migration guides; no novel architectural
decision is being made.)

## Risks

| Risk | Likelihood | Mitigation |
|------|-----------|-----------|
| `spinner.Tick` package-level vs. method: compilation error | Low | Confirmed `Model.Tick()` method exists in bubbles v2; existing code already uses `s.spinner.Tick` (method value) |
| lipgloss `Color()` return type `color.Color` breaks other call sites | Low | `lipgloss.NewStyle().Foreground(lipgloss.Color(...))` accepted: `Foreground` takes `color.Color` in v2 |
| `ansi.HideCursor` / `ShowCursor` are string constants not functions | Confirmed | Verified in `x/ansi@v0.11.6`: they are string constants used with `fmt.Fprint` |
| huh v2 builder API changed beyond import path | Low | Confirmed quick-start guide says "import path only" for most users; huh v2 README shows identical builder API |
| Transitive dependency conflicts after tidy | Medium | Run `go mod tidy` and `go build ./...` together; resolve any replace directives or version conflicts before committing |
| `CursorBack` removal in spinner accessible mode: mild visual regression | Low | `CursorBack` only cleaned up after accessible spinner render — drop it entirely; the cursor is reset by ShowCursor escape anyway |
| `isTerm()` behavior change (now correctly detects TTY) | Intentional | This fixes a latent bug; callers get styled output on actual terminals for the first time |
