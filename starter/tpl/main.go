// Copyright © 2021 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	"embed"
)

//go:embed docs/*
var DocFS embed.FS

//go:embed go.mod.tpl
var GoModTemplate []byte

func MainTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`)
}

func RootTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
{{- if .Viper }}
	"errors"
	"fmt"{{ end }}
	"os"

	"github.com/bketelsen/toolbox/cobra"{{- if .Viper }}
	"github.com/spf13/viper"{{ end }}

)

{{ if .Viper -}}
var cfgFile string
{{- end }}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "A brief description of your application",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
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
{{- if .Viper }}
	cobra.OnInitialize(initConfig)
{{ end }}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
{{ if .Viper }}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
{{ else }}
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
{{ end }}
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

{{ if .Viper -}}
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".{{ .AppName }}" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".{{ .AppName }}")
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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
{{- end }}
`)
}

func AddCommandTemplate() []byte {
	return []byte(`/*
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
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("{{ .CmdName }} called")
		cmd.Logger.Info("{{ .CmdName }} called")
	},
}

func init() {
	{{ .CmdParent }}.AddCommand({{ .CmdName }}Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// {{ .CmdName }}Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// {{ .CmdName }}Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
`)
}

func AddDocsTemplate() []byte {
	return []byte(`/*
{{ .Project.Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
	"github.com/bketelsen/toolbox/cobra/doc"

)

// {{ .CmdName }}Cmd represents the {{ .CmdName }} command
var {{ .CmdName }}Cmd = &cobra.Command{
	Use:   "{{ .CmdName }}",
	Short: "Generates documentation for the project",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	Run: func(cmd *cobra.Command, args []string) {
		bp := viper.GetString("basepath")
		cmd.Logger.Info("Base path for documentation", "basepath", bp)
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return bp + "/docs/cli/" + strings.ToLower(base) + "/"
		}
		filePrepender := func(filename string) string {
			now := time.Now().Format(time.RFC3339)
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/docs/cli/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}
		err := os.MkdirAll("./docs/content/docs/cli/", 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = doc.GenMarkdownTreeCustom(rootCmd, "./docs/content/docs/cli/", filePrepender, linkHandler)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	{{ .CmdParent }}.AddCommand({{ .CmdName }}Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// {{ .CmdName }}Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// {{ .CmdName }}Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	docsCmd.Flags().StringP("basepath", "b", "", "Base path for the documentation (default is /)")
	viper.BindPFlag("basepath", docsCmd.Flags().Lookup("basepath"))
}
const fmTemplate = ` + "`" + `---
date: %s
title: "%s"
slug: %s
url: %s
---
` + "`")
}
