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
		Short:   "Add a docs command and documentation template to a Cobra Application",
		Long:    `Docs (cobra-cli docs) will create a new docs command`,

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

func doDocs(_ *cobra.Command) error {
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	modName := getModImportPath()

	commandName := "docs"
	command := &Command{
		CmdName:   commandName,
		CmdParent: parentName,
		Project: &Project{
			PkgName:      modName,
			AbsolutePath: wd,
			Legal:        getLicense(),
			Copyright:    copyrightLine(),
			AppName:      path.Base(modName),
		},
	}

	return command.Docs()
}
