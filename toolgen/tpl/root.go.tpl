/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
	"os"

	"github.com/bketelsen/toolbox/cobra"
    goversion "github.com/bketelsen/toolbox/go-version"
)

var appname = "{{ .AppName }}"
var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

var bversion = buildVersion(version, commit, date, builtBy, treeState)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Version: bversion.String(),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// https://www.asciiart.eu/text-to-ascii-art to make your own
// just make sure the font doesn't have backticks in the letters or
// it will break the string quoting
var asciiName = `
  █████                      ████                              
 ░░███                      ░░███                              
 ███████    ██████   ██████  ░███   ███████  ██████  ████████  
░░░███░    ███░░███ ███░░███ ░███  ███░░███ ███░░███░░███░░███ 
  ░███    ░███ ░███░███ ░███ ░███ ░███ ░███░███████  ░███ ░███ 
  ░███ ███░███ ░███░███ ░███ ░███ ░███ ░███░███░░░   ░███ ░███ 
  ░░█████ ░░██████ ░░██████  █████░░███████░░██████  ████ █████
   ░░░░░   ░░░░░░   ░░░░░░  ░░░░░  ░░░░░███ ░░░░░░  ░░░░ ░░░░░ 
                                   ███ ░███                    
                                  ░░██████                     
                                   ░░░░░░                      
`

// buildVersion builds the version info for the application
func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(appname, "An application that does cool things.", "{{ .Config.GetString "repository" }}"),
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
