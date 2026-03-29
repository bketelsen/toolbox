// Command dashboard demonstrates the toolbox ui package alongside toolbox.App.
// It simulates a service dashboard that displays styled console messages,
// tabular data, a spinner, and optional interactive prompts.
//
// Run with no flags for non-interactive output:
//
//	go run ./_examples/dashboard
//
// Add --interactive to exercise the huh-based Confirm, Prompt, and Option widgets:
//
//	go run ./_examples/dashboard --interactive
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bketelsen/toolbox"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/cobra"
)

type service struct {
	Name   string `table:"name,default_sort"`
	Status string `table:"status"`
	Uptime string `table:"uptime"`
	Port   int    `table:"port"`
}

var (
	interactive bool
	showTable   bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Display a service dashboard",
		Long: ui.Long(
			"A demo CLI that exercises the toolbox/ui package: styled console messages, "+
				"tables, spinners, and interactive prompts.",
			ui.Example{
				Description: "Show the dashboard with default settings",
				Command:     "dashboard",
			},
			ui.Example{
				Description: "Show the dashboard with interactive prompts",
				Command:     "dashboard --interactive",
			},
		),
		RunE: runDashboard,
	}

	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "enable interactive prompts")
	rootCmd.Flags().BoolVar(&showTable, "table", false, "show only the service table")

	app := toolbox.App{Version: "0.3.0"}
	if err := app.Run(rootCmd); err != nil {
		os.Exit(1)
	}
}

func runDashboard(cmd *cobra.Command, args []string) error {
	services := []service{
		{"api-gateway", "running", "14d 3h", 8080},
		{"auth-service", "running", "14d 3h", 8081},
		{"worker-pool", "degraded", "2h 15m", 9090},
		{"cache", "running", "14d 3h", 6379},
		{"scheduler", "stopped", "0s", 9091},
	}

	// --- Table-only mode ---
	if showTable {
		table, err := ui.DisplayTable(services, "name", nil)
		if err != nil {
			return toolbox.OutputJSONError("failed to render table", err)
		}
		fmt.Println(table)
		return nil
	}

	// --- Console messages ---
	ui.Info("Service Dashboard", "Gathering service status...")
	ui.InfoPrefix("cluster-01", "Cluster", "Primary region: us-east-1")

	// --- Styled output ---
	fmt.Println()
	for _, svc := range services {
		switch svc.Status {
		case "degraded":
			ui.WarnPrefix(ui.LinePrefixWarning, "Degraded",
				fmt.Sprintf("%s on port %d", svc.Name, svc.Port),
				fmt.Sprintf("uptime: %s", svc.Uptime),
			)
		case "stopped":
			ui.ErrorPrefix(ui.LinePrefixCross, "Stopped",
				svc.Name,
				fmt.Sprintf("expected on port %d", svc.Port),
			)
		}
	}

	// --- Key-value pairs and style helpers ---
	fmt.Println()
	ui.Info("Metadata",
		ui.KeyValuePair("Region", "us-east-1"),
		ui.KeyValuePair("Environment", "production"),
		fmt.Sprintf("Last checked: %s", ui.Timestamp(time.Now())),
		fmt.Sprintf("Config file: %s", ui.Code("~/.config/dashboard.yaml")),
	)

	// --- Spinner ---
	fmt.Println()
	err := ui.NewSpinner().
		Title("Refreshing metrics...").
		Type(ui.Dots).
		ActionWithErr(func(ctx context.Context) error {
			time.Sleep(800 * time.Millisecond)
			return nil
		}).
		Run()
	if err != nil {
		return err
	}
	ui.Success("Metrics refreshed", "All dashboards up to date")

	// --- JSON output ---
	if toolbox.OutputJSON(services) {
		return nil
	}

	// --- Interactive prompts (guarded) ---
	if !interactive {
		return nil
	}

	fmt.Println()
	ui.TitleInfo("Interactive Mode")

	// Option select
	var action string
	ui.Option(
		"Choose an action",
		"Select what to do with the stopped service",
		&action,
		[]string{"Restart", "Investigate", "Ignore"},
	)

	// Text prompt
	var reason string
	ui.Prompt(
		"Reason",
		"Why are you taking this action?",
		"e.g. scheduled maintenance",
		&reason,
	)

	// Confirm
	var confirmed bool
	ui.Confirm(
		fmt.Sprintf("%s scheduler?", action),
		fmt.Sprintf("Reason: %s", reason),
		&confirmed,
	)

	if confirmed {
		ui.SuccessPrefix(ui.LinePrefixCheck, "Done",
			fmt.Sprintf("Action %q applied to scheduler", action),
			fmt.Sprintf("Reason: %s", reason),
		)
	} else {
		ui.WarnPrefix(ui.LinePrefixWarning, "Cancelled", "No changes made")
	}

	return nil
}
