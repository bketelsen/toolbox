// Command sync demonstrates a multi-subcommand CLI using toolbox.App with
// slug structured logging, NewReporter, BindViper, and context-propagated
// loggers. It simulates syncing files between local and remote storage.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bketelsen/toolbox"
	"github.com/bketelsen/toolbox/slug"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync files between local and remote storage",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := toolbox.BindViper(cmd.Root()); err != nil {
				return err
			}

			level := slog.LevelInfo
			if toolbox.Verbose {
				level = slog.LevelDebug
			}
			if toolbox.Silent {
				level = slog.LevelError + 1 // suppress everything
			}
			slog.SetDefault(slog.New(slug.NewHandler(os.Stderr, &slug.Options{
				Level: level,
			})))
			return nil
		},
	}

	rootCmd.AddCommand(pushCmd(), statusCmd())

	app := toolbox.App{Version: "0.2.0"}
	if err := app.Run(rootCmd); err != nil {
		os.Exit(1)
	}
}

func pushCmd() *cobra.Command {
	var remote string
	cmd := &cobra.Command{
		Use:   "push [path]",
		Short: "Push local files to remote",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			log := slog.With("remote", remote, "path", path)
			r := toolbox.NewReporter()

			log.Debug("starting push", "dry_run", toolbox.DryRun)

			// Phase 1: Scan
			r.Step(1, 3, "Scanning local files")
			time.Sleep(150 * time.Millisecond)
			files := []string{"config.yaml", "data.json", "README.md"}
			log.Info("scan complete", "files_found", len(files))

			// Phase 2: Diff
			r.Step(2, 3, "Comparing with remote")
			time.Sleep(150 * time.Millisecond)
			changed := files[:2] // simulate 2 changed files
			log.Info("diff complete",
				slog.Group("stats",
					slog.Int("changed", len(changed)),
					slog.Int("unchanged", len(files)-len(changed)),
				),
			)
			if len(changed) == 0 {
				r.Message("Already up to date")
				return nil
			}

			// Phase 3: Upload
			r.Step(3, 3, "Uploading changes")
			if toolbox.DryRun {
				for _, f := range changed {
					r.Message("would upload %s", f)
					log.Debug("skipped upload (dry run)", "file", f)
				}
				r.Complete("Dry run complete, no files uploaded", nil)
				return nil
			}

			for i, f := range changed {
				time.Sleep(100 * time.Millisecond)
				log.Debug("uploaded file", "file", f, "bytes", 1024*(i+1))
				r.Progress((i+1)*100/len(changed), f)
			}

			result := map[string]any{
				"remote":   remote,
				"uploaded": len(changed),
				"files":    changed,
			}
			if toolbox.OutputJSON(result) {
				return nil
			}
			r.Complete(fmt.Sprintf("Pushed %d files to %s", len(changed), remote), result)
			return nil
		},
	}
	cmd.Flags().StringVar(&remote, "remote", "s3://my-bucket", "remote storage URL")
	return cmd
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show sync status",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := slog.With("command", "status")

			log.Debug("checking sync status")

			status := struct {
				LastSync string `json:"last_sync"`
				Pending  int    `json:"pending"`
				Errors   int    `json:"errors"`
			}{
				LastSync: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
				Pending:  3,
				Errors:   1,
			}

			if toolbox.OutputJSON(status) {
				return nil
			}

			r := toolbox.NewReporter()
			r.Message("Last sync: %s", status.LastSync)
			r.Message("Pending: %d files", status.Pending)
			if status.Errors > 0 {
				log.Warn("sync errors detected", slug.Err(errors.New("1 file failed checksum")))
				r.Warning("%d sync errors", status.Errors)
			}
			return nil
		},
	}
}
