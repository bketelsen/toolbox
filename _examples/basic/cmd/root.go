/*
Copyright © 2025 Mister Smith <you@yourdomain.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/viper"
)

var cfgFile string
var version string
var commit string

func versionString() string {
	if len(commit) > 7 {
		commit = commit[:7]
	}
	if len(commit) == 0 {
		commit = "unknown"
	}
	if len(version) == 0 {
		version = "unknown"
	}
	return fmt.Sprintf("%s (%s)", version, commit)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "basic",
	Version: versionString(),
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// set the slog default logger to the cobra logger
		slog.SetDefault(cmd.Logger)
		// set the slog level to debug if the --debug flag is set
		if viper.GetBool("verbose") {
			cmd.SetLogLevel(slog.LevelDebug)
		}

		cmd.Logger.Debug("Debug message, shows with --verbose flag")
		cmd.SetLogLevel(slog.LevelInfo)
		cmd.Logger.Debug("Debug message 2 -- won't show.")
		cmd.Logger.Info("Info message, always shows")
		cmd.SetLogLevel(slog.LevelDebug)

		action := func() {
			time.Sleep(1 * time.Second)
		}
		if err := ui.NewSpinner().Title("Preparing your order...").Action(action).Run(); err != nil {
			fmt.Println("Failed:", err)
			return
		}
		fmt.Println("Order up!")
		cmd.Println(ui.Warn("This is a warning message"))
		cmd.Println(ui.Info("This is an informational message"))
		cmd.Println(ui.Error("This is an error message"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.basic.yaml)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".basic" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".basic")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !errors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && errors.As(err, notFound):
		// The config file is optional, we shouldn't exit when the config is not found
		break
	default:
		if viper.GetBool("verbose") {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
