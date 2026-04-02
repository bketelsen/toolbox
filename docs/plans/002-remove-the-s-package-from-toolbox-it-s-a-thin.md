---
devTaskId: 002
title: Remove the s/ package from toolbox. It's a thin wrapper over iancoleman/strcase
status: Approved
created: 2026-04-02
author: architect
adrNeeded: false
designOnly: false
complexity: simple
---

# Plan: Remove the s/ package from toolbox. It's a thin wrapper over iancoleman/strcase with no tests. The CasedString struct and case conversion functions are used by cmd/scaffold/templates.go and cmd/scaffold/main.go. Inline the necessary functionality directly into cmd/scaffold (use strcase directly), then delete the s/ directory entirely. Remove the strcase dependency from go.mod if nothing else imports it. Run `go mod tidy` and ensure all tests pass.

## Summary

Delete the `s/` sub-package from the toolbox module. The package is a thin, untested wrapper around `github.com/iancoleman/strcase` that exposes case-conversion helpers and a `CasedString` struct. A full codebase audit reveals **zero callers** of `github.com/bketelsen/toolbox/s` anywhere in the repository — the task description's claim that scaffold files use it is not reflected in current code. Removing the package eliminates dead code, shrinks the public API surface, and allows the `iancoleman/strcase` direct dependency to be dropped from `go.mod` via `go mod tidy`.

## Context

**What exists in `s/strings.go`:**
- `Pascal`, `Camel`, `Snake`, `ScreamingSnake`, `Kebab` — thin delegates to `strcase.*`
- `Upper`, `Lower` — thin delegates to `strings.*`
- `CasedFn` type alias, `Cased` variadic helper
- `CasedString` struct and `ToCasedString` constructor

**Caller audit (verified by grep):**
- `grep -rn '"github.com/bketelsen/toolbox/s"' ./...` → **no results**
- `cmd/scaffold/main.go` imports `toolbox`, `toolbox/slug`, `toolbox/ui`, cobra — no `s` package
- `cmd/scaffold/templates.go` imports `toolbox/slug` only — no `s` package
- `_examples/*` import only `toolbox/slug`
- `slug/handler_test.go` imports `toolbox/slug` (the slog handler) — no `s` package

**Dependency:**
- `github.com/iancoleman/strcase v0.3.0` is listed as a **direct** dependency in `go.mod`
- Its sole consumer is `s/strings.go`; after `s/` is deleted, `go mod tidy` will drop it entirely (or demote to `// indirect` if pulled in transitively by another dep)

**No inlining needed:** Because no file imports `toolbox/s`, there is nothing to inline. The task description anticipated usage that was either already removed or never merged.

## Architecture

No architectural changes. This is a pure deletion with dependency cleanup — identical shape to the `n/` package removal in plan `002` (n-package).

## Implementation Steps

1. **Delete the `s/` directory** — remove `s/strings.go` (file and directory).
   - Files: `s/strings.go` (delete)
   - Command: `rm -rf s/`

2. **Tidy the module** — run `go mod tidy` to remove `github.com/iancoleman/strcase` from the direct dependency list and update `go.sum`.
   - Files: `go.mod`, `go.sum` (modified by tooling)

3. **Verify build and tests** — confirm nothing is broken.
   - Commands: `go build ./...` then `go test ./...`

## Files to Change

- `s/strings.go` — **delete** (entire file; directory removed with it)
- `go.mod` — `github.com/iancoleman/strcase` entry removed (or moved to `// indirect`) after `go mod tidy`
- `go.sum` — updated automatically by `go mod tidy`

## Test Strategy

```bash
# 1. Delete the package
rm -rf s/

# 2. Tidy dependencies
go mod tidy

# 3. Confirm build is clean
go build ./...

# 4. Confirm all tests still pass
go test ./...

# 5. Confirm strcase is no longer a direct dep
grep 'iancoleman/strcase' go.mod   # should be absent or marked // indirect
```

No new tests are required — the package had none, the deletion removes the only code that needed them, and no callers exist to regression-test.

## ADR Required

needed: no

## Risks

- **Unlikely external consumer**: if a downstream project imports `github.com/bketelsen/toolbox/s`, that import will break. Given the package is undocumented, untested, and has zero internal references, this risk is very low. A brief release-notes entry is sufficient mitigation.
- **strcase kept as indirect**: if another dependency transitively requires `iancoleman/strcase`, `go mod tidy` will retain it as `// indirect` in `go.mod`. This is acceptable and expected.
- **Scaffolding future use**: if `CasedString` was intended for future use in scaffold templates (e.g., to offer `{{.Name.Pascal}}` in template data), that intent must be addressed separately if/when scaffold templates need case variants. No action is needed now since the functionality is not currently wired up.
