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
	installerCmd = &cobra.Command{
		Use:   "installer",
		Short: "Add an installer to your project",
		Long: `Installer (starter installer) adds a GitHub downloader script to your project.
		
Included:
* Installer script - GitHub download/install script for your project`,

		Run: func(cmd *cobra.Command, args []string) {

			err := doExtras(cmd, false, false, false, false, false, false, true, overwrite)
			if err != nil {
				cmd.Logger.Error(err.Error())
				cobra.CheckErr(err)
			}

			cmd.Logger.Info("installer created")
		},
	}
)

func init() {
	installerCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing files")

}
