---
devTaskId: "002"
title: "Remove the n/ package from toolbox"
status: Draft
created: 2026-04-02
author: architect
adrNeeded: false
designOnly: false
complexity: simple
---

# Plan: Remove the n/ package from toolbox

## Summary

The `n/` package contains a single 16-line `FreePort()` function with no tests and no callers anywhere in the repository. We are removing it entirely to reduce package surface area. No inlining is required because nothing imports it.

## Context

- `n/network.go` exports a single function `FreePort() (int, error)` that resolves a free TCP port via `net.ListenTCP`.
- A `grep` for `"github.com/bketelsen/toolbox/n"` across all `.go` files returns zero results — the package is unused within the repo.
- No other packages, examples, or tests import `n/`.
- The change is purely subtractive: delete one file and one directory.

## Architecture

No architectural changes are required. This is a dead-code removal. There are no callers to update, no interfaces to change, and no data-flow impact.

## Implementation Steps

1. **Delete the n/ directory** — remove `n/network.go` and the `n/` directory.
   - Files: `n/network.go` (delete), `n/` (rmdir)
   - Command: `rm -rf n/`

2. **Run go mod tidy** — ensure `go.mod` / `go.sum` are clean after deletion.
   - Command: `go mod tidy`

3. **Verify the build** — confirm nothing broke.
   - Command: `go build ./...`

4. **Run the test suite** — confirm all tests still pass.
   - Command: `go test ./...`

## Files to Change

- `n/network.go` — delete entirely; the package is unused
- `n/` (directory) — remove after deleting the file

## Test Strategy

No new tests are needed (there were none to begin with). Verification consists of:

```bash
rm -rf n/
go mod tidy
go build ./...
go test ./...
```

All four commands should complete without errors. If `go test ./...` was green before the change it will remain green after, since no callers existed.

## ADR Required

needed: no

## Risks

- **Low risk overall** — no callers exist, so there is no breakage possible inside the repo.
- **External consumers** — if any downstream module (outside this repo) imports `github.com/bketelsen/toolbox/n`, deleting the package is a breaking change. The Go module is not v2+, so a major-version bump is not strictly required; however, the maintainer should confirm no external consumers exist (e.g., via `pkg.go.dev` or a GitHub code search) before merging. Since this is a library, releasing a patch or minor version bump after removal is advisable.
- **go.sum drift** — `go mod tidy` should handle any orphaned checksum entries; verify `go.sum` has no unexpected changes.
