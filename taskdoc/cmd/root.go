/*
Copyright © 2025 Brian Ketelsen

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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/bketelsen/toolbox/cobra"
)

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
	Use:     "taskdoc",
	Version: versionString(),
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error
		var summary TaskSummary
		if len(args) == 0 {
			summary, err = getTaskSummary(".")
			if err != nil {
				return err

			}
		} else {
			summary, err = getTaskSummary(args[0])
			if err != nil {
				return err
			}
		}
		if len(summary.Tasks) == 0 {
			cobra.CheckErr("No tasks found")
		}

		summaryFile, err := os.Create("TASKS.md")
		if err != nil {
			return err
		}
		defer summaryFile.Close()
		summaryTemplate := template.Must(template.New("summary").
			Parse(string(TaskSummaryTemplate())))
		err = summaryTemplate.Execute(summaryFile, summary)
		if err != nil {
			return err
		}

		return nil
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.taskdoc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type TaskSummary struct {
	Tasks []struct {
		Name     string `json:"name,omitempty"`
		Desc     string `json:"desc,omitempty"`
		Summary  string `json:"summary,omitempty"`
		Aliases  []any  `json:"aliases,omitempty"`
		UpToDate bool   `json:"up_to_date,omitempty"`
		Location struct {
			Line     int    `json:"line,omitempty"`
			Column   int    `json:"column,omitempty"`
			Taskfile string `json:"taskfile,omitempty"`
		} `json:"location,omitempty"`
	} `json:"tasks,omitempty"`
	Location string `json:"location,omitempty"`
}

// getTaskSummary returns a TaskSummary struct from the Taskfile.yml
func getTaskSummary(path string) (TaskSummary, error) {
	var taskSummary TaskSummary
	taskfile := filepath.Join(path, "Taskfile.yml")
	if !exists(taskfile) {
		cobra.CheckErr("Taskfile.yml not found")
	}
	out, err := exec.Command("task", "--list-all", "--json").Output()
	if err != nil {
		return taskSummary, err
	}

	err = json.Unmarshal(out, &taskSummary)

	return taskSummary, err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		cobra.CheckErr(err)
	}
	return true
}

func TaskSummaryTemplate() []byte {
	return []byte(`
### Available Tasks
{{ range .Tasks }}
#### {{ .Name }} - {{ .Desc }}
{{ if .Summary }}
{{ .Summary }}
{{ end }}

Run this task:

` + "```" + `
task {{ .Name }}
` + "```" + `
{{ end }}
`)
}
