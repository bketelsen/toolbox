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
	docsCmd = &cobra.Command{
		Use:     "docs",
		Aliases: []string{"command"},
		Short:   "Add a documentation generation and docs template to your project.",
		Long: `Docs (toolgen docs) will create a new gendocs command.
		
The gendocs command will generate a command named "gendocs" which will
generate documentation for your application in markdown format suitable
for use with any static site generator.

Optionally, you can add support for GitHub Pages by using the -g flag.`,

		Run: func(cmd *cobra.Command, args []string) {

			cmd.Logger.Info("creating docs command")
			cobra.CheckErr(doDocs(cmd, cmd.Config().GetBool("ghpages"), cmd.GlobalConfig()))

		},
	}
)

func init() {
	docsCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
	docsCmd.Flags().BoolP("ghpages", "g", false, "add github pages support")
	rootCmd.AddCommand(docsCmd)
}
