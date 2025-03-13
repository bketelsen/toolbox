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
	"path"

	cobra "github.com/bketelsen/toolbox/cobra"
)

var (
	docsCmd = &cobra.Command{
		Use:     "docs",
		Aliases: []string{"command"},
		Short:   "Add a gendocs command and documentation template to a Cobra Application",
		Long: `Docs (starter docs) will create a new gendocs command.
		
The gendocs command will generate documentation for your application
in markdown format suitable for use with any static site generator.`,

		Run: func(cmd *cobra.Command, args []string) {

			cmd.Logger.Info("creating docs command")
			cobra.CheckErr(doDocs(cmd))

		},
	}
)

func init() {
	docsCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name (e.g. github.com/spf13/hugo)")
	docsCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
	cobra.CheckErr(docsCmd.Flags().MarkDeprecated("package", "this operation has been removed."))
}

func doDocs(cmd *cobra.Command) error {
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	modName := getModImportPath()

	SetKeyAndSaveConfig("docs", true)
	SetKeyAndSaveConfig("basepath", "/"+path.Base(modName))

	config, err := GetActiveConfig()
	if err != nil {
		return err
	}

	commandName := "gendocs"
	command := &Command{
		CmdName:   commandName,
		CmdParent: parentName,
		Project: &Project{
			PkgName:      modName,
			AbsolutePath: wd,
			Legal:        getLicense(),
			Copyright:    copyrightLine(),
			AppName:      path.Base(modName),
			Config:       &config,
		},
	}

	err = command.Docs()
	if err != nil {
		return err
	}
	// overwrite the taskfile
	return doExtras(cmd, true, false, false, false, false, false, true)

}
