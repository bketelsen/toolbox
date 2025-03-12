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
	taskfileCmd = &cobra.Command{
		Use:   "taskfile",
		Short: "Add a Taskfile template to your project",
		Long:  `Taskfile (starter taskfile) adds a Taskfile to your project.`,

		Run: func(cmd *cobra.Command, args []string) {

			err := doTaskfile(cmd)
			if err != nil {
				cmd.Logger.Error(err.Error())
				cobra.CheckErr(err)
			}

			cmd.Logger.Info("installer created")
		},
	}
)

func init() {
	taskfileCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing files")

}

func doTaskfile(_ *cobra.Command) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	modName := getModImportPath()
	binName := filepath.Base(modName)
	repository := viper.GetString("repository")
	owner, repo := getOwnerRepo(repository)

	extras := &Extras{
		Taskfile:       true,
		GoReleaser:     false,
		DevContainer:   false,
		ActionsGo:      false,
		ActionsPages:   false,
		ActionsRelease: false,
		Installer:      false,
		Overwrite:      overwrite,
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
