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
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

var (
	wizardCmd = &cobra.Command{
		Use:   "wizard",
		Short: "Initialize a toolbox application with interactive prompts in the current directory.",
		Long: `Wizard (toolgen wizard) will create a new application, with a license
and the appropriate structure for a Cobra-based CLI application.
It will prompt for the application name, license type, and other
information needed to create the application.
`,

		Run: func(cmd *cobra.Command, args []string) {
			var confirm bool
			var appName string
			var modName string
			var name string
			var email string
			var ghRepo string

			huh.NewConfirm().
				Title("This wizard will create a new Go CLI application. Continue?").
				Value(&confirm).
				Run()

			if !confirm {
				fmt.Println("Aborting...")
				return
			}

			// get the name of the current directory
			dir, err := os.Getwd()
			if err != nil {
				cmd.Println(ui.Error("Failed to get current directory"), err.Error())
				return
			}
			appName = filepath.Base(dir)

			// check to see if the directory we are in is empty
			entries, err := os.ReadDir(".")
			if err != nil {
				cmd.Println(ui.Error("Failed to read directory"), err.Error())
				return
			}
			if len(entries) > 0 {
				huh.NewConfirm().
					Title("The current directory is not empty. Do you want to continue?").
					Description("This will overwrite any existing files in the current directory.").
					Value(&confirm).
					Run()
				if !confirm {
					fmt.Println("Aborting...")
					return
				}
			}

			// Start with application name
			huh.NewInput().
				Title("What is the name of your application?").
				Description("Please enter the name of your application.\nApplication name should be a single word, all lowercase.").
				Placeholder("myapp").
				Prompt("> ").
				Validate(appNameValidator).
				Value(&appName).
				Run()

			// default modname to the app name
			modName = appName

			prompt("What is the name of your go module?",
				"Please enter the name of your go module.\nIt can be a single word like `myapp` or a URL like `github.com/you/yourapp`",
				"github.com/you/yourapp", &modName)

			// get the license
			ls := slices.Collect(maps.Keys(Licenses))
			huh.NewSelect[string]().
				Title("What license do you want to use?").
				Description("Please select a license for your application.").
				Options(huh.NewOptions(ls...)...).
				Value(&userLicense).
				Run()

			// get the author name and email

			prompt("What is the name of the License holder?",
				"Enter your full name, or the legal name of the entity that holds the license.",
				"The Go Authors", &name)

			prompt("What is the email address of the License holder?",
				"Enter your email address, or the address of the entity that holds the license.",
				"gophers@go.dev", &email)

			// get the GitHub repository URL

			prompt("What is the (eventual) GitHub repository URL?",
				"Enter the GitHub URL where this project will be stored.",
				"https://github.com/you/appname", &ghRepo)

			if ghRepo == "" {
				ghRepo = "https://github.com/you/appname"
			}

			huh.NewConfirm().
				Title("Ready to create application. Continue?").
				Value(&confirm).
				Run()

			if !confirm {
				fmt.Println("Aborting...")
				return
			}

			// create a configuration file
			cfgMap := defaultConfig
			// remove unwanted keys
			delete(cfgMap, "config")
			delete(cfgMap, "write-config")
			cfgMap[keyAuthor] = name
			cfgMap[keyEmail] = email // todo validate this
			cfgMap[keyLicense] = userLicense
			cfgMap[keyModuleName] = modName
			cfgMap[keyAppName] = appName
			cfgMap[keyRepository] = ghRepo // TODO normalize this to a URL, no trailing slashes, no .git etc.

			// run go mod init
			cmd.Logger.Info("Running go mod init", "module", modName)
			if err := runGoModInit(modName); err != nil {
				cmd.Println(ui.Error("Failed to run go mod init"), err.Error())
				return
			}
			newCfg := "toolgen.yaml"
			v := viper.New()
			v.SetConfigFile("toolgen.yaml")
			v.SetConfigType("yaml")
			v.AddConfigPath(".")

			if err := v.MergeConfigMap(cfgMap); err != nil {
				cmd.Println(ui.Error("Failed to merge config map"), err.Error())
				return
			}
			if err := v.WriteConfigAs(newCfg); err != nil {
				cmd.Println(ui.Error("Failed to write config file"), err.Error())
				return
			}
			cmd.Logger.Info("Creating new application", "name", appName, "module", modName, "license", userLicense, "email", email, "name", name)

			abspath, err := initializeProject([]string{"."}, v)
			if err != nil {
				cmd.Println(ui.Error("Failed to create project"), err.Error())
				return
			}
			// run go mod tidy
			cmd.Logger.Info("Tidying up", "module", modName)
			// if err := runGoModTidy(); err != nil {
			// 	cmd.Println(ui.Error("Failed to run go mod tidy"), err.Error())
			// 	return
			// }
			cmd.Logger.Info("Project created", "path", abspath)

			huh.NewConfirm().
				Title("Do you want to add some bling?").
				Value(&confirm).
				Run()

			if !confirm {
				fmt.Println("Happy trails!")
				return
			}
			var docs, taskfile, releaser, devcontainer, actionsGo, actionsPages, installer bool
			huh.NewConfirm().
				Title("Add generated documentation website starter?").
				Value(&confirm).
				Run()

			if confirm {
				cmd.Logger.Info("Adding docs")
				docs = true
			}
			huh.NewConfirm().
				Title("Add devcontainer configuration?").
				Value(&confirm).
				Run()

			if confirm {
				cmd.Logger.Info("Adding .devcontainer files")
				devcontainer = true
			}
			huh.NewConfirm().
				Title("Add GitHub Actions for Go?").
				Value(&confirm).
				Run()

			if confirm {
				cmd.Logger.Info("Adding GitHub Actions for Go")
				actionsGo = true
			}
			huh.NewConfirm().
				Title("Add bash installer script?").
				Value(&confirm).
				Run()

			if confirm {
				cmd.Logger.Info("Adding installer script")
				installer = true
			}
			huh.NewConfirm().
				Title("Add GoReleaser configuration?").
				Value(&confirm).
				Run()

			if confirm {
				cmd.Logger.Info("Adding GoReleaser config")
				releaser = true
			}
			huh.NewConfirm().
				Title("Add Taskfile?").
				Value(&confirm).
				Run()

			if confirm {
				cmd.Logger.Info("Adding Taskfile")
				taskfile = true
			}

			// install the docs templates
			if docs {
				err = doDocs(cmd, true, v)
				if err != nil {
					cmd.Logger.Error(err.Error())
					cobra.CheckErr(err)
				}
			}
			// install the extras
			// overwrite to allow for re-running some tasks
			err = doExtras(cmd, taskfile, releaser, devcontainer, actionsGo, actionsPages, installer, true, v)
			if err != nil {
				cmd.Logger.Error(err.Error())
				cobra.CheckErr(err)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(wizardCmd)

}

func runGoModInit(modName string) error {

	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func appNameValidator(source string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("application name cannot be empty")
	}
	if strings.ContainsAny(source, " ") {
		return fmt.Errorf("application name cannot contain spaces")
	}
	if strings.ContainsAny(source, "/") {
		return fmt.Errorf("application name cannot contain slashes")
	}
	if strings.ContainsAny(source, "\\") {
		return fmt.Errorf("application name cannot contain backslashes")
	}
	if strings.ContainsAny(source, ":") {
		return fmt.Errorf("application name cannot contain colons")
	}
	if strings.ContainsAny(source, "?") {
		return fmt.Errorf("application name cannot contain question marks")
	}
	if strings.ContainsAny(source, "*") {
		return fmt.Errorf("application name cannot contain asterisks")
	}
	if strings.ContainsAny(source, "_") {
		return fmt.Errorf("application name cannot contain underscores")
	}
	if strings.HasPrefix(source, "-") {
		return fmt.Errorf("application name cannot start with a dash")
	}
	if strings.HasSuffix(source, "-") {
		return fmt.Errorf("application name cannot end with a dash")
	}
	if strings.ContainsAny(source, ".") {
		return fmt.Errorf("application name cannot contain periods")
	}
	// so many more things to check...

	return nil
}

func doExtras(_ *cobra.Command,
	taskfile, releaser, devcontainer, actionsGo, actionsPages, installer bool,
	replace bool,
	config *viper.Viper,
) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	modName := getModImportPath()
	binName := filepath.Base(modName)
	repository := config.GetString("repository")
	owner, repo := getOwnerRepo(repository)

	// config, err := GetActiveConfig()
	// if err != nil {
	// 	return err
	// }

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
			Config:       config,
		},
	}
	return extras.Create()
}

func doDocs(cmd *cobra.Command, withPages bool, config *viper.Viper) error {
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	modName := getModImportPath()

	config.Set("docs", true)
	config.Set("basepath", "/"+path.Base(modName))
	repository := config.GetString("repository")
	owner, repo := getOwnerRepo(repository)
	config.Set("owner", owner)
	config.Set("repo", repo)

	commandName := "gendocs"
	command := &Command{
		CmdName:   commandName,
		CmdParent: parentName,
		Project: &Project{
			PkgName:      modName,
			AbsolutePath: wd,
			Legal:        getLicense(config),
			Copyright:    copyrightLine(config),
			AppName:      path.Base(modName),
			Config:       config,
			Repository:   repository,
			Owner:        owner,
			Repo:         repo,
		},
	}

	err = command.Docs()
	if err != nil {
		return err
	}
	// overwrite the taskfile
	return doExtras(cmd, true, false, false, false, withPages, false, true, config)

}
