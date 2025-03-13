// Copyright © 2021 Steve Francia <spf@spf13.com>.
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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
)

var (
	initCmd = &cobra.Command{
		Use:     "init [path]",
		Aliases: []string{"initialize", "initialise", "create"},
		Short:   "Initialize a Cobra Application",
		Long: `Initialize (cobra-cli init) will create a new application, with a license
and the appropriate structure for a Cobra-based CLI application.

Cobra init must be run inside of a go module (please run "go mod init <MODNAME>" first)
`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var comps []string
			var directive cobra.ShellCompDirective
			if len(args) == 0 {
				comps = cobra.AppendActiveHelp(comps, "Optionally specify the path of the go module to initialize")
				directive = cobra.ShellCompDirectiveDefault
			} else if len(args) == 1 {
				comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
				directive = cobra.ShellCompDirectiveNoFileComp
			} else {
				comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
				directive = cobra.ShellCompDirectiveNoFileComp
			}
			return comps, directive
		},
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Logger.Info("initializing project")

			// then get the configs for viper
			if viper.GetBool("UseViper") {
				SetKeyAndSaveConfig("useViper", true)
			}

			// then get the configs for the license
			if userLicense != "" {
				SetKeyAndSaveConfig("license", userLicense)
			} else {
				SetKeyAndSaveConfig("license", getLicense())
			}
			// get this written to disk before we start doing anything
			// from here on out we will be using the config file
			// and not the flags, so we need to make sure
			// to propagate any changed flags back to the config file
			cmd.Logger.Info("writing initial config")
			err := writeConfigWithDefaults("starter.yaml")
			if err != nil {
				cmd.Logger.Error(err.Error())
			}
			projectPath, err := initializeProject(args)
			if err != nil {
				cmd.Logger.Error(err.Error())
			}

			cobra.CheckErr(err)

			err = goGet("github.com/bketelsen/toolbox/cobra")
			if err != nil {
				cmd.Logger.Error(err.Error())
			}
			cobra.CheckErr(err)
			if viper.GetBool("useViper") {
				cobra.CheckErr(goGet("github.com/spf13/viper"))
			}

			// generate docs if requested
			if viper.GetBool("docs") {
				cmd.Logger.Info("creating docs command")
				err = doDocs(cmd)
				if err != nil {
					cmd.Logger.Error(err.Error())
					cobra.CheckErr(err)
				}
				fmt.Printf("%s created in %s\n", "docs", projectPath)

			}
			// generate extras if requested
			if viper.GetBool("extras") {
				cmd.Logger.Info("creating extras")
				// overwrite the files if they exist (taskfile should already be there)
				err = doExtras(cmd, true, true, true, true, true, true, true)
				if err != nil {
					cmd.Logger.Error(err.Error())
					cobra.CheckErr(err)
				}
				fmt.Printf("Extras created in %s\n", projectPath)
			}
			fmt.Printf("Your Cobra application is ready at\n%s\n", projectPath)
		},
	}
)

func init() {
	initCmd.Flags().BoolP("docs", "d", false, "generate documentation")
	viper.BindPFlag("docs", initCmd.Flags().Lookup("docs"))
	initCmd.Flags().BoolP("extras", "e", false, "generate extra configs")
	viper.BindPFlag("extras", initCmd.Flags().Lookup("extras"))

}

func initializeProject(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if args[0] != "." {
			wd = fmt.Sprintf("%s/%s", wd, args[0])
		}
	}

	modName := getModImportPath()
	config, err := GetActiveConfig()
	if err != nil {
		return "", err
	}
	project := &Project{
		AbsolutePath: wd,
		PkgName:      modName,
		Legal:        getLicense(),
		Copyright:    copyrightLine(),
		Viper:        viper.GetBool("useViper"),
		AppName:      path.Base(modName),
		Config:       &config,
	}

	if err := project.Create(); err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}

func getModImportPath() string {
	mod, cd := parseModInfo()
	return path.Join(mod.Path, fileToURL(strings.TrimPrefix(cd.Dir, mod.Dir)))
}

func fileToURL(in string) string {
	i := strings.Split(in, string(filepath.Separator))
	return path.Join(i...)
}

func parseModInfo() (Mod, CurDir) {
	var mod Mod
	var dir CurDir

	m := modInfoJSON("-m")
	cobra.CheckErr(json.Unmarshal(m, &mod))

	// Unsure why, but if no module is present Path is set to this string.
	if mod.Path == "command-line-arguments" {
		cobra.CheckErr("Please run `go mod init <MODNAME>` before `starter init`")
	}

	e := modInfoJSON("-e")
	cobra.CheckErr(json.Unmarshal(e, &dir))

	return mod, dir
}

type Mod struct {
	Path, Dir, GoMod string
	Main             bool
}

type CurDir struct {
	Dir string
}

func goGet(mod string) error {
	return exec.Command("go", "get", mod).Run()
}

func modInfoJSON(args ...string) []byte {
	cmdArgs := append([]string{"list", "-json"}, args...)
	out, err := exec.Command("go", cmdArgs...).Output()
	cobra.CheckErr(err)

	return out
}
