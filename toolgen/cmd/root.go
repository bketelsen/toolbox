/*
Copyright © 2025 Brian Ketelsen <bketelsen@gmail.com>

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
	"log/slog"
	"os"
	"strings"

	"github.com/bketelsen/toolbox/cobra"
	goversion "github.com/bketelsen/toolbox/go-version"
	"github.com/bketelsen/toolbox/ui"

	"github.com/spf13/viper"
)

var cfgFile string
var appname = "toolgen"
var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

var userLicense string

var bversion = buildVersion(version, commit, date, builtBy, treeState)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     appname,
	Version: bversion.String(),
	Short:   "A generator for Go CLI applications.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set log level based on the --verbose flag
		if cmd.GlobalConfig().GetBool("verbose") {
			cmd.SetLogLevel(slog.LevelDebug)
			cmd.Logger.Debug("Debug logging enabled")
		}

		// if the --write-config flag is set, write the config file
		// and exit
		if cmd.GlobalConfig().GetBool("write-config") {
			// make a copy so we can remove some things
			//	cfgMap := cmd.GlobalConfig().AllSettings()
			cfgMap := defaultConfig
			// remove unwanted keys
			delete(cfgMap, "config")
			delete(cfgMap, "write-config")
			v := viper.New()
			if err := v.MergeConfigMap(cfgMap); err != nil {
				cmd.Println(ui.Error("Failed to merge config map"), err.Error())
				return
			}
			if err := v.WriteConfigAs(cfgFile); err != nil {
				cmd.Println(ui.Error("Failed to write config file"), err.Error())
				return
			}
			cmd.Println(ui.Info("Config file written to", cfgFile))
			cmd.Println(ui.Info("You can now edit the config file and run the command again"))
			cmd.Println(ui.Info("You can also set the config file with the --config flag"))
			os.Exit(0)

		}
	},
	InitConfig: func() *viper.Viper {
		config := viper.New()
		config.SetEnvPrefix(appname)
		config.AutomaticEnv()
		config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
		config.SetDefault("verbose", true)
		config.SetDefault("license", "MIT")
		config.SetDefault("author", "Your Name")
		config.SetDefault("email", "you@youremail.com")
		config.SetDefault("repository", "https://github.com/your/project")
		config.SetConfigType("yaml")
		config.SetConfigFile(cfgFile)
		config.AddConfigPath(".")
		_ = config.ReadInConfig()
		return config
	},

	Run: func(cmd *cobra.Command, args []string) {

		// set the slog default logger to the cobra logger
		slog.SetDefault(cmd.Logger)

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "toolgen.yaml", "use config file")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
	rootCmd.Flags().BoolP("write-config", "", false, "write config file")
}

var asciiName = `
████████╗ ██████╗  ██████╗ ██╗      ██████╗ ███████╗███╗   ██╗
╚══██╔══╝██╔═══██╗██╔═══██╗██║     ██╔════╝ ██╔════╝████╗  ██║
   ██║   ██║   ██║██║   ██║██║     ██║  ███╗█████╗  ██╔██╗ ██║
   ██║   ██║   ██║██║   ██║██║     ██║   ██║██╔══╝  ██║╚██╗██║
   ██║   ╚██████╔╝╚██████╔╝███████╗╚██████╔╝███████╗██║ ╚████║
   ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝
`

// buildVersion builds the version info for the application
// https://www.asciiart.eu/text-to-ascii-art
func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(appname, "A generator for Go CLI applications.", "https://bketelsen.github.io/toolbox"),
		goversion.WithASCIIName(asciiName),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}

		},
	)
}
