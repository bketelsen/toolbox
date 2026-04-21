// Package toolbox provides CLI convenience functions for Frostyard tools,
// wrapping charmbracelet/fang and spf13/cobra with standardized version
// injection, common flags, JSON output helpers, and reporter factory.
package toolbox

import (
	"context"
	"fmt"
	"log/slog"
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

	// LogFile is the path to write log output. Empty means stderr.
	LogFile string
	// Logger is the slog.Logger initialized during PersistentPreRunE.
	Logger *slog.Logger

	// logFile holds the open file handle for cleanup.
	logFile *os.File
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
//
// Run installs a PersistentPreRunE on cmd that initializes the logger and
// chains any existing PersistentPreRunE/PersistentPreRun set on cmd before
// this call. Consumers MUST set their custom pre-run logic on the root
// command (before calling Run), NOT on subcommands — a subcommand's
// PersistentPreRunE would shadow the root's and bypass logger setup.
func (a *App) Run(cmd *cobra.Command) error {
	a.defaults()
	a.registerFlags(cmd)

	capturedPreRunE := cmd.PersistentPreRunE
	capturedPreRun := cmd.PersistentPreRun

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if err := a.setupLogger(); err != nil {
			return err
		}
		if capturedPreRunE != nil {
			return capturedPreRunE(c, args)
		}
		if capturedPreRun != nil {
			capturedPreRun(c, args)
		}
		return nil
	}
	cmd.PersistentPreRun = nil

	defer a.Close()
	return fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(a.VersionString()),
		fang.WithNotifySignal(os.Interrupt, os.Kill),
	)
}
