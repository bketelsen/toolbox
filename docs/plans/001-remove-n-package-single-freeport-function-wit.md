---
devTaskId: 001
title: Remove n/ package — single FreePort() function with no tests
status: Approved
created: 2026-04-02
author: architect
adrNeeded: false
designOnly: false
complexity: simple
---

# Plan: Remove n/ package — single FreePort() function with no tests

## Summary

The `n/` package contains a single function (`FreePort()`) with no tests and no consumers anywhere in the codebase. It adds surface area to the library without value. This plan removes the package entirely in a single, safe step.

## Context

- **Package location:** `n/network.go` — one file, ~15 lines, package name `n`
- **Function:** `FreePort() (port int, err error)` — asks the OS for a free TCP port via `net.ListenTCP` on `localhost:0`
- **Tests:** none
- **Internal consumers:** none (`grep -r "toolbox/n"` returns only the package file itself)
- **Example consumers:** none found in `_examples/`
- **External consumers:** this is a library; external callers may import `github.com/bketelsen/toolbox/n`, but the function is trivially replaceable with stdlib `net` directly. No changelog or deprecation notice exists for it.
- The `go.mod` module path is `github.com/bketelsen/toolbox`; removing the sub-package is a semver-compatible change since the module is still pre-1.0 (v0.x).

## Architecture

No architectural change. This is a pure deletion — no replacement, no migration, no new code.

The `net` stdlib provides identical capability in ~5 lines; any consumer can replicate `FreePort()` trivially or use a purpose-built library.

## Implementation Steps

1. **Delete the `n/` directory** — remove `n/network.go` (the only file)
   - Files: `n/network.go` (delete)
   - What: `rm -rf n/`

2. **Verify build and tests pass** — ensure nothing references the deleted package
   - Run: `go build ./...`
   - Run: `go test ./...`
   - Run: `go vet ./...`

## Files to Change

- `n/network.go` — delete; the entire `n/` package is removed

## Test Strategy

- `go build ./...` — must succeed with zero errors
- `go test ./...` — all existing tests must pass
- `go vet ./...` — no new vet errors
- `grep -r "toolbox/n" .` — must return no results after deletion

No new tests are needed; we are removing code, not adding it.

## ADR Required

needed: no

## Risks

- **External consumers breaking:** Any project that imports `github.com/bketelsen/toolbox/n` will get a compile error after this change. Since the package has no tests or documentation and the module is pre-1.0, this risk is accepted. A release note or CHANGELOG entry is recommended.
- **None inside the repo:** Grep confirms zero internal imports of the package.
