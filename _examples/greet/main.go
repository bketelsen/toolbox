// Command greet demonstrates the toolbox.App wrapper with slug structured
// logging. It showcases common flags (--verbose, --json, --dry-run) and
// how slug integrates as the slog handler alongside OutputJSON.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bketelsen/toolbox"
	"github.com/bketelsen/toolbox/slug"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "greet [name]",
		Short: "Greet someone with structured logging",
		Args:  cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := slog.LevelInfo
			if toolbox.Verbose {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slug.NewHandler(os.Stderr, &slug.Options{
				Level: level,
			})))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			greeting := fmt.Sprintf("Hello, %s!", name)

			slog.Debug("resolved greeting", "name", name, "greeting", greeting)

			if toolbox.DryRun {
				slog.Info("dry run, skipping output", "greeting", greeting)
				return nil
			}

			if toolbox.OutputJSON(map[string]string{"greeting": greeting, "name": name}) {
				slog.Debug("wrote JSON output")
				return nil
			}

			fmt.Println(greeting)
			slog.Info("greeting sent", "name", name)
			return nil
		},
	}

	app := toolbox.App{Version: "0.1.0"}
	if err := app.Run(rootCmd); err != nil {
		os.Exit(1)
	}
}
