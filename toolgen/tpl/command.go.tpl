/*
{{ .Project.Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
	"github.com/bketelsen/toolbox/cobra"
)

// {{ .CmdName }}Cmd represents the {{ .CmdName }} command
var {{ .CmdName }}Cmd = &cobra.Command{
	Use:   "{{ .CmdName }}",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("{{ .CmdName }} called")
		cmd.Logger.Info("toggle","value",cmd.Config().GetBool("toggle"))
	},
}

func init() {
	{{ .CmdParent }}.AddCommand({{ .CmdName }}Cmd)

	// Here you will define your flags and configuration settings.

	// It is not necessary to add Viper bindings to your flags.
	// rootCmd's persistent flags are added to `cmd.GlobalConfig()`
	// and this command's flags are added to `cmd.Config()`

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// {{ .CmdName }}Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	{{ .CmdName }}Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}