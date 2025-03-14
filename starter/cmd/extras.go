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
	overwrite                                                            bool
	taskfile, releaser, devcontainer, actionsGo, actionsPages, installer bool
	extrasCmd                                                            = &cobra.Command{
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

			err := doExtras(cmd, taskfile, releaser, devcontainer, actionsGo, actionsPages, installer, overwrite)
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
	extrasCmd.Flags().BoolVarP(&taskfile, "taskfile", "t", false, "Add Taskfile")
	extrasCmd.Flags().BoolVarP(&releaser, "goreleaser", "g", false, "Add GoReleaser")
	extrasCmd.Flags().BoolVarP(&devcontainer, "devcontainer", "d", false, "Add DevContainer")
	extrasCmd.Flags().BoolVarP(&actionsGo, "actions-go", "b", false, "Add GitHub Actions for Go")
	extrasCmd.Flags().BoolVarP(&actionsPages, "actions-pages", "p", false, "Add GitHub Actions for Pages")
	extrasCmd.Flags().BoolVarP(&installer, "installer", "i", false, "Add Installer script")
	//	viper.BindPFlags(extrasCmd.Flags())
}

func doExtras(_ *cobra.Command,
	taskfile, releaser, devcontainer, actionsGo, actionsPages, installer bool,
	replace bool,
) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	modName := getModImportPath()
	binName := filepath.Base(modName)
	repository := viper.GetString("repository")
	owner, repo := getAndSetOwnerRepo(repository)

	config, err := GetActiveConfig()
	if err != nil {
		return err
	}

	extras := &Extras{
		Taskfile:       taskfile,
		GoReleaser:     releaser,
		DevContainer:   devcontainer,
		ActionsGo:      actionsGo,
		ActionsPages:   actionsPages,
		ActionsRelease: releaser,
		Installer:      installer,
		Overwrite:      replace,
		Project: &Project{
			PkgName:      modName,
			AbsolutePath: wd,
			AppName:      binName,
			Repository:   repository,
			Owner:        owner,
			Repo:         repo,
			Config:       &config,
		},
	}
	return extras.Create()
}
