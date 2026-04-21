package toolbox

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// JSONOutput is set to true when the --json flag is passed.
var JSONOutput bool

// Verbose is set to true when the --verbose or -v flag is passed.
var Verbose bool

// DryRun is set to true when the --dry-run or -n flag is passed.
var DryRun bool

// Silent is set to true when the --silent or -s flag is passed.
var Silent bool

// registerFlags adds --json, --verbose, --dry-run, --silent, and --log-file as persistent flags on cmd.
func (a *App) registerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&JSONOutput, "json", false, "output in JSON format")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "n", false, "dry run mode (no actual changes)")
	cmd.PersistentFlags().BoolVarP(&Silent, "silent", "s", false, "suppress all progress output")
	cmd.PersistentFlags().StringVar(&a.LogFile, "log-file", "", "write log output to this file (default: stderr)")
}

// BindViper binds the common flags (--json, --verbose, --dry-run, --silent, --log-file) to viper.
// Call this in a PersistentPreRunE if your app uses viper for config management.
func BindViper(cmd *cobra.Command) error {
	for _, name := range []string{"json", "verbose", "dry-run", "silent", "log-file"} {
		if err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
			return err
		}
	}
	return nil
}
