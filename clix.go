// Package toolbox provides CLI convenience functions for Frostyard tools,
// wrapping charmbracelet/fang and spf13/cobra with standardized version
// injection, common flags, JSON output helpers, and reporter factory.
package toolbox

import (
	"context"
	"fmt"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

// App holds build-time metadata for a CLI application.
// Create one in main() and call Run() to execute the root command.
type App struct {
	// Version is the version string of the application (e.g., "1.2.3").
	Version string
	// Commit is the git commit hash that built this application.
	Commit string
	// Date is the timestamp when this application was built.
	Date string
	// BuiltBy is the identifier of the builder (e.g., "ci" or "local").
	BuiltBy string
}

// defaults fills zero-value fields with sensible defaults.
func (a *App) defaults() {
	if a.Version == "" {
		a.Version = "dev"
	}
	if a.Commit == "" {
		a.Commit = "none"
	}
	if a.Date == "" {
		a.Date = "unknown"
	}
	if a.BuiltBy == "" {
		a.BuiltBy = "local"
	}
}

// VersionString returns a formatted version string including commit, date,
// and builder info. Example: "1.2.3 (Commit: abc) (Date: 2026-01-01) (Built by: ci)"
func (a *App) VersionString() string {
	a.defaults()
	return fmt.Sprintf("%s (Commit: %s) (Date: %s) (Built by: %s)",
		a.Version, a.Commit, a.Date, a.BuiltBy)
}

// Run registers common persistent flags on cmd, then executes the command
// via fang.Execute with the formatted version string and signal handling.
func (a *App) Run(cmd *cobra.Command) error {
	a.defaults()
	registerFlags(cmd)
	return fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(a.VersionString()),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	)
}
