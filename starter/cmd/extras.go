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
)

var (
	overwrite bool
	extrasCmd = &cobra.Command{
		Use:   "extras",
		Short: "Add extras to your project",
		Long:  `Extras (starter extras) adds additional functionality to your project.`,

		Run: func(cmd *cobra.Command, args []string) {

			err := doExtras(cmd)
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

func doExtras(_ *cobra.Command) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	modName := getModImportPath()
	binName := filepath.Base(modName)

	extras := &Extras{
		Taskfile:       true,
		GoReleaser:     true,
		DevContainer:   true,
		ActionsGo:      true,
		ActionsPages:   true,
		ActionsRelease: true,
		Overwrite:      overwrite,
		Project: &Project{
			PkgName:      modName,
			AbsolutePath: wd,
			AppName:      binName,
		},
	}
	return extras.Create()
}
