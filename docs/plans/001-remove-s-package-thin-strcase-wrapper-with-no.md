---
devTaskId: 001
title: Remove s/ package — thin strcase wrapper with no tests
status: Approved
created: 2026-04-02
author: architect
adrNeeded: false
designOnly: false
complexity: simple
---

# Plan: Remove s/ package — thin strcase wrapper with no tests

## Summary

Delete the `s/` sub-package from the toolbox module. It is a thin, untested wrapper around `github.com/iancoleman/strcase` with no callers anywhere in the module or its examples. Removing it eliminates dead code, reduces the public API surface, and allows the `iancoleman/strcase` direct dependency to be dropped from `go.mod`.

## Context

- `s/strings.go` exports `Pascal`, `Camel`, `Snake`, `ScreamingSnake`, `Kebab`, `Upper`, `Lower`, `Cased`, `CasedString`, and `ToCasedString` — all thin one-liners delegating to `strcase` or stdlib `strings`.
- A `grep` across all `.go` files in the repo confirms **zero** internal consumers of `github.com/bketelsen/toolbox/s`.
- `github.com/iancoleman/strcase v0.3.0` appears as a **direct** dependency in `go.mod` solely because of this package. After removal, `go mod tidy` will demote or drop it.
- No test files exist in the `s/` directory.
- The `slug/` package (a slog handler) is unrelated and must not be touched.

## Architecture

No architectural changes. This is a pure deletion with a dependency cleanup.

## Implementation Steps

1. **Delete the `s/` directory** — remove `s/strings.go` (and the directory itself).
   - Files: `s/strings.go` (delete)

2. **Tidy the module** — run `go mod tidy` to remove or demote `github.com/iancoleman/strcase` from the direct dependency list in `go.mod` / `go.sum`.
   - Files: `go.mod`, `go.sum` (modified by tooling)

3. **Verify the build** — run `go build ./...` and `go test ./...` to confirm nothing is broken.

## Files to Change

- `s/strings.go` — delete (entire file; directory removed with it)
- `go.mod` — `github.com/iancoleman/strcase` moves to indirect or is dropped entirely after `go mod tidy`
- `go.sum` — updated automatically by `go mod tidy`

## Test Strategy

```bash
# 1. Delete the package
rm -rf s/

# 2. Tidy dependencies
go mod tidy

# 3. Confirm build is clean
go build ./...

# 4. Confirm tests pass
go test ./...

# 5. Confirm strcase is no longer a direct dep
grep 'iancoleman/strcase' go.mod   # should be absent or marked // indirect
```

No new tests are required — the package had no tests, and the deletion removes the only code that needed them.

## ADR Required

needed: no

## Risks

- **Unlikely external consumer**: if any downstream project imports `github.com/bketelsen/toolbox/s`, that import will break. Given the package is undocumented, untested, and not referenced internally, this risk is very low; a major-version or deprecation notice is not required for an unannounced internal utility package. A brief entry in CHANGELOG/release notes is sufficient.
- **strcase still needed indirectly**: if another dependency pulls in `iancoleman/strcase` transitively, `go mod tidy` will keep it as `// indirect` — this is fine and expected.
