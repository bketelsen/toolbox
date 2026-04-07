package ui_test

// Tests for prompt.go require a real terminal and interactive input via huh v2.
// The Prompt() function directly calls huh.NewInput().Run() which blocks on terminal I/O
// and requires actual user input from a TTY. Unit testing this function without a real
// terminal would require either:
// 1. Mocking the huh library's internal state (fragile, version-dependent)
// 2. Running with a pseudo-terminal (complex test setup)
//
// The Prompt function is a thin wrapper around huh's input prompt and has no
// testable business logic beyond parameter forwarding.
//
// For integration testing, use manual testing or e2e tests with a real terminal.
// See SO-6 for coverage goals and architecture decisions.
