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
	"os"
	"path/filepath"

	cobra "github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
)

var (
	overwrite bool
	extrasCmd = &cobra.Command{
		Use:   "extras",
		Short: "Add all the extras to your project",
		Long: `Extras (starter extras) adds additional functionality to your project.

Included:
* Taskfile - automated task runner
* GoReleaser - release automation
* DevContainer - VSCode devcontainer
* GitHub Actions for Go - CI/CD
* GitHub Actions for Pages - GitHub Pages deployment
* GitHub Actions for Release - GitHub Release automation for GoReleaser
* Installer script - GitHub download/install script for your project`,

		Run: func(cmd *cobra.Command, args []string) {

			err := doExtras(cmd, true, true, true, true, true, true, true, overwrite)
			if err != nil {
				cmd.Logger.Error(err.Error())
				cobra.CheckErr(err)
			}

			cmd.Logger.Info("extras added")
		},
	}
)

func init() {
	extrasCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing files")

}

func doExtras(_ *cobra.Command,
	taskfile bool,
	releaser bool,
	devcontainer bool,
	actionsGo bool,
	actionsPages bool,
	actionsRelease bool,
	installer bool,
	replace bool,
) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	modName := getModImportPath()
	binName := filepath.Base(modName)
	repository := viper.GetString("repository")
	owner, repo := getOwnerRepo(repository)

	extras := &Extras{
		Taskfile:       taskfile,
		GoReleaser:     releaser,
		DevContainer:   devcontainer,
		ActionsGo:      actionsGo,
		ActionsPages:   actionsPages,
		ActionsRelease: actionsRelease,
		Installer:      installer,
		Overwrite:      replace,
		Project: &Project{
			PkgName:      modName,
			AbsolutePath: wd,
			AppName:      binName,
			Repository:   repository,
			Owner:        owner,
			Repo:         repo,
		},
	}
	return extras.Create()
}
