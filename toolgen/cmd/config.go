package cmd

import (
	"os/user"
	"strings"

	cobra "github.com/bketelsen/toolbox/cobra"
)

var (
	keyVerbose      = "verbose"
	keyAuthor       = "author"
	keyEmail        = "email"
	keyRepository   = "repository"
	keyLicense      = "license"
	keyAppName      = "appname"
	keyModuleName   = "modulename"
	keyTaskfile     = "taskfile"
	keyReleaser     = "releaser"
	keyDevContainer = "devcontainer"
	keyActionsGo    = "actionsgo"
	keyActionsPages = "actionspages"
	keyDocs         = "docs"
)

var defaultConfig map[string]interface{} = map[string]interface{}{
	keyVerbose:    false,
	keyAuthor:     getName(),
	keyEmail:      "you@yourdomain.com",
	keyRepository: "https://github.com/you/project",
	keyLicense:    "MIT",
}

// given a github URL, return the owner and repo
func getOwnerRepo(githubURL string) (string, string) {
	parts := strings.Split(githubURL, "/")

	if len(parts) < 5 {
		cobra.CheckErr("Invalid github URL")
	}
	owner := parts[3]
	repo := parts[4]

	return owner, repo
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
