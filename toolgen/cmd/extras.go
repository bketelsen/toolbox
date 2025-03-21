// Copyright © 2015 Steve Francia <spf@spf13.com>.
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

package cmd

import (
	cobra "github.com/bketelsen/toolbox/cobra"
)

var (
	extrasCmd = &cobra.Command{
		Use:     "extras",
		Aliases: []string{"bling"},
		Short:   "Add optional configuration to your project.",
		Long: `Extras (starter extras) adds additional functionality to your project.

Options available:
* Taskfile - automated task runner
* GoReleaser - release automation
* DevContainer - VSCode devcontainer
* GitHub Actions for Go - CI/CD
* GitHub Actions for Pages - GitHub Pages deployment
* GitHub Actions for Release - GitHub Release automation for GoReleaser
* Installer script - GitHub download/install script for your project`,

		Run: func(cmd *cobra.Command, args []string) {

			cmd.Logger.Info("adding extras")
			// install the extras
			// overwrite to allow for re-running some tasks
			err := doExtras(cmd,
				cmd.Config().GetBool("taskfile"),
				cmd.Config().GetBool("goreleaser"),
				cmd.Config().GetBool("devcontainer"),
				cmd.Config().GetBool("github-actions"),
				cmd.Config().GetBool("github-pages"),
				cmd.Config().GetBool("installer"),
				cmd.Config().GetBool("overwrite"),
				cmd.GlobalConfig())
			if err != nil {
				cmd.Logger.Error(err.Error())
				cobra.CheckErr(err)
			}

		},
	}
)

func init() {
	extrasCmd.Flags().BoolP("taskfile", "", false, "add a Taskfile to your project")
	extrasCmd.Flags().BoolP("goreleaser", "", false, "add a GoReleaser config to your project")
	extrasCmd.Flags().BoolP("devcontainer", "", false, "add a devcontainer to your project")
	extrasCmd.Flags().BoolP("github-actions", "", false, "add GitHub Actions to your project")
	extrasCmd.Flags().BoolP("github-pages", "", false, "add GitHub Pages to your project")
	extrasCmd.Flags().BoolP("installer", "", false, "add an installer script to your project")
	extrasCmd.Flags().BoolP("overwrite", "o", false, "allow overwriting existing files")
	rootCmd.AddCommand(extrasCmd)
}
