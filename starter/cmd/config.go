package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	cobra "github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	// Verbose enables verbose output
	Verbose bool `mapstructure:"verbose"`
	// FullName is the full name of the user
	// used for copyright attribution
	FullName string `mapstructure:"full-name"`
	// Email is the email address of the user
	Email string `mapstructure:"email"`
	// UseViper is true if viper should be used
	UseViper bool `mapstructure:"viper"`
	// Docs is true if the docs command should be created
	Docs bool `mapstructure:"docs"`
	// BasePath is the base http path for the project
	// used for documentation links
	// e.g. / on a normal web server
	// or /project on a GitHub Pages site
	BasePath string `mapstructure:"basepath"`
	// Taskfile is true if the taskfile should be/is created
	Taskfile bool `mapstructure:"taskfile"`
	// Releaser is true if the GoReleaser file should be created
	Releaser bool `mapstructure:"releaser"`
	// Devcontainer is true if the devcontainer files should be created
	Devcontainer bool `mapstructure:"devcontainer"`
	// ActionsGo is true if the GitHub Actions for Go should be created
	ActionsGo bool `mapstructure:"actions-go"`
	// ActionsPages is true if the GitHub Actions for Pages should be created
	ActionsPages bool `mapstructure:"actions-pages"`
	// ActionsRelease is true if the GitHub Actions for GoReleaser should be created
	ActionsRelease bool `mapstructure:"actions-release"`
	// Installer is true if the installer script should be created
	Installer bool `mapstructure:"installer"`
	// Repository is the URL of the GitHub repository e.g https://github.com/bketelsen/toolbox
	Repository string `mapstructure:"repository"`
	// Owner is the owner of the GitHub repository e.g bketelsen
	Owner string `mapstructure:"owner"`
	// Repo is the name of the GitHub repository e.g toolbox
	Repo string `mapstructure:"repo"`
	// License is the name of the license to use e.g MIT
	License string `mapstructure:"license"`
}

var (
	keyVerbose        = "verbose"
	keyFullName       = "full-name"
	keyEmail          = "email"
	keyTaskfile       = "taskfile"
	keyReleaser       = "releaser"
	keyDevcontainer   = "devcontainer"
	keyActionsGo      = "actions-go"
	keyActionsPages   = "actions-pages"
	keyActionsRelease = "actions-release"
	keyInstaller      = "installer"
	keyDocs           = "docs"
	keyRepository     = "repository"
	keyOwner          = "owner"
	keyRepo           = "repo"
	keyLicense        = "license"
	keyBasePath       = "basepath"
	keyUseViper       = "useviper"
)

var defaultConfig map[string]interface{} = map[string]interface{}{
	keyVerbose:        false,
	keyFullName:       getName(),
	keyEmail:          "you@yourdomain.com",
	keyTaskfile:       false,
	keyReleaser:       false,
	keyDevcontainer:   false,
	keyActionsGo:      false,
	keyActionsPages:   false,
	keyActionsRelease: false,
	keyInstaller:      false,
	keyDocs:           false,
	keyRepository:     "https://github.com/you/project",
	keyOwner:          "you",
	keyRepo:           "project",
	keyLicense:        "MIT",
	keyUseViper:       true,
	keyBasePath:       "/",
}

// given a github URL, return the owner and repo
func getAndSetOwnerRepo(githubURL string) (string, string) {
	parts := strings.Split(githubURL, "/")
	if len(parts) < 5 {
		cobra.CheckErr("Invalid github URL")
	}
	owner := parts[3]
	repo := parts[4]
	defaultConfig["owner"] = owner
	defaultConfig["repo"] = repo
	SetKeyAndSaveConfig("owner", owner)
	SetKeyAndSaveConfig("repo", repo)
	return owner, repo
}

func GetActiveConfig() (Config, error) {
	var c Config
	err := viper.Unmarshal(&c)
	return c, err
}

func SaveActiveConfig() error {
	err := viper.WriteConfig()
	if err != nil {
		log.Println("Error writing configuration:", err)
		return err
	}
	return nil
}

func SetKeyAndSaveConfig(key string, value interface{}) error {
	notFound := &viper.ConfigFileNotFoundError{}

	err := viper.ReadInConfig()
	if err != nil && !errors.As(err, notFound) {
		return err
	}
	viper.Set(key, value)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func writeConfigWithDefaults(path string) error {

	// Check if the config file already exists
	if exists(path) {
		return viper.ConfigFileAlreadyExistsError(path)
	}
	err := viper.ReadInConfig()

	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !errors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && errors.As(err, notFound):
		// The config file may not exist, we shouldn't exit when the config is not found
		// but we should create later with the current settings merged with the defaults
		// and write it out
		break
	default:
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	err = viper.MergeConfigMap(defaultConfig)
	if err != nil {
		log.Println("Error merging starter configuration:", err)
		return err
	}
	// Use SafeWriteConfigAs to avoid overwriting existing files
	err = viper.SafeWriteConfigAs(path)
	if err != nil {
		log.Println("Error writing starter configuration:", err)
		return err
	}
	return nil
}

func getName() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	name := usr.Name
	if name == "" {
		return "Toolbox User"
	}
	return name
}
