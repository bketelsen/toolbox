package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"text/template"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/starter/tpl"
)

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	PkgName string
	// Repository is the full URL of the repository
	Repository string
	// Owner is the owner of the repository
	Owner string
	// Repo is the name of the repository
	Repo         string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
	Summary      *TaskSummary
	Config       *Config
}

type Command struct {
	CmdName   string
	CmdParent string
	*Project
}

type Extras struct {
	Taskfile       bool
	Installer      bool
	GoReleaser     bool
	DevContainer   bool
	ActionsGo      bool
	ActionsPages   bool
	ActionsRelease bool
	Overwrite      bool
	*Project
}

func (e *Extras) Create() error {
	if e.Installer {
		tf := fmt.Sprintf("%s/install.sh", e.AbsolutePath)
		if exists(tf) && !e.Overwrite {
			return fmt.Errorf("found existing install.sh, use -o to overwrite")
		}
		installer, err := os.Create(tf)
		if err != nil {
			return err
		}
		defer installer.Close()
		// change the delimiters to [[ and ]]
		// this is because the template uses {{ and }} for the template
		installTemplate := template.Must(template.New("installer").
			Delims("[[", "]]").
			Parse(string(tpl.InstallScriptTemplate)))
		err = installTemplate.Execute(installer, e)
		if err != nil {
			return err
		}
		// make the file executable
		if err := os.Chmod(tf, 0755); err != nil {
			return err
		}
		SetKeyAndSaveConfig("installer", true)

	}
	if e.Taskfile {
		tf := fmt.Sprintf("%s/Taskfile.yml", e.AbsolutePath)
		if exists(tf) && !e.Overwrite {
			return fmt.Errorf("found existing Taskfile.yml, use -o to overwrite")
		}
		taskfile, err := os.Create(tf)
		if err != nil {
			return err
		}
		defer taskfile.Close()
		// change the delimiters to [[ and ]]
		// this is because the template uses {{ and }} for the template
		taskfileTemplate := template.Must(template.New("taskfile").
			Delims("[[", "]]").
			Parse(string(tpl.TaskfileTemplate)))
		err = taskfileTemplate.Execute(taskfile, e)
		if err != nil {
			return err
		}
		tfchecks := fmt.Sprintf("%s/Taskfile.checks.yml", e.AbsolutePath)
		if exists(tfchecks) && !e.Overwrite {
			return fmt.Errorf("found existing Taskfile.checks.yml, use -o to overwrite")
		}
		taskfileChecks, err := os.Create(tfchecks)
		if err != nil {
			return err
		}
		defer taskfileChecks.Close()
		// change the delimiters to [[ and ]]
		// this is because the template uses {{ and }} for the template
		taskfileChecksTemplate := template.Must(template.New("taskfilechecks").
			Delims("[[", "]]").
			Parse(string(tpl.TaskfileChecksTemplate)))
		err = taskfileChecksTemplate.Execute(taskfileChecks, e)
		if err != nil {
			return err
		}
		summary, err := getTaskSummary(e.AbsolutePath)
		if err != nil {
			return err
		}
		e.Summary = &summary
		summaryFile, err := os.Create(fmt.Sprintf("%s/TASKS.md", e.AbsolutePath))
		if err != nil {
			return err
		}
		defer summaryFile.Close()
		summaryTemplate := template.Must(template.New("summary").
			Parse(string(tpl.TaskSummaryTemplate())))
		err = summaryTemplate.Execute(summaryFile, e)
		if err != nil {
			return err
		}
		SetKeyAndSaveConfig(keyTaskfile, true)

	}
	if e.GoReleaser {
		rf := fmt.Sprintf("%s/.goreleaser.yml", e.AbsolutePath)
		if exists(rf) && !e.Overwrite {
			return fmt.Errorf("found existing .goreleaser.yml, use -o to overwrite")
		}
		goReleaser, err := os.Create(rf)
		if err != nil {
			return err
		}
		defer goReleaser.Close()
		goReleaserTemplate := template.Must(template.New("goreleaser").
			Delims("[[", "]]").
			Parse(string(tpl.GoReleaserTemplate)))
		err = goReleaserTemplate.Execute(goReleaser, e)
		if err != nil {
			return err
		}
		tfrelease := fmt.Sprintf("%s/Taskfile.release.yml", e.AbsolutePath)
		if exists(tfrelease) && !e.Overwrite {
			return fmt.Errorf("found existing Taskfile.release.yml, use -o to overwrite")
		}
		taskfileRelease, err := os.Create(tfrelease)
		if err != nil {
			return err
		}
		defer taskfileRelease.Close()
		// change the delimiters to [[ and ]]
		// this is because the template uses {{ and }} for the template
		taskfileReleaseTemplate := template.Must(template.New("taskfilerelease").
			Delims("[[", "]]").
			Parse(string(tpl.TaskfileReleaseTemplate)))
		err = taskfileReleaseTemplate.Execute(taskfileRelease, e)
		if err != nil {
			return err
		}
		summary, err := getTaskSummary(e.AbsolutePath)
		if err != nil {
			return err
		}
		e.Summary = &summary
		summaryFile, err := os.Create(fmt.Sprintf("%s/TASKS.md", e.AbsolutePath))
		if err != nil {
			return err
		}
		defer summaryFile.Close()
		summaryTemplate := template.Must(template.New("summary").
			Parse(string(tpl.TaskSummaryTemplate())))
		err = summaryTemplate.Execute(summaryFile, e)
		if err != nil {
			return err
		}
		SetKeyAndSaveConfig(keyReleaser, true)
	}
	if e.DevContainer {
		dc := fmt.Sprintf("%s/.devcontainer/devcontainer.json", e.AbsolutePath)
		if exists(dc) && !e.Overwrite {
			return fmt.Errorf("found existing .devcontainer/devcontainer.json file, use -o to overwrite")
		}
		if err := os.MkdirAll(fmt.Sprintf("%s/.devcontainer", e.AbsolutePath), 0755); err != nil {
			return err
		}
		devContainer, err := os.Create(dc)
		if err != nil {
			return err
		}
		defer devContainer.Close()
		devContainerTemplate := template.Must(template.New("devcontainer").Parse(string(tpl.DevContainerTemplate)))
		err = devContainerTemplate.Execute(devContainer, e)
		if err != nil {
			return err
		}
		dockerfile, err := os.Create(fmt.Sprintf("%s/.devcontainer/Dockerfile", e.AbsolutePath))
		if err != nil {
			return err
		}
		defer dockerfile.Close()
		dockerfileTemplate := template.Must(template.New("dockerfile").Parse(string(tpl.DockerfileTemplate)))
		err = dockerfileTemplate.Execute(dockerfile, e)
		if err != nil {
			return err
		}
		SetKeyAndSaveConfig(keyDevcontainer, true)

	}
	if e.ActionsGo || e.ActionsPages || e.ActionsRelease {
		if err := os.MkdirAll(fmt.Sprintf("%s/.github/workflows", e.AbsolutePath), 0755); err != nil {
			return err
		}
	}
	if e.ActionsGo {
		goyml := fmt.Sprintf("%s/.github/workflows/go.yml", e.AbsolutePath)
		if exists(goyml) && !e.Overwrite {
			return fmt.Errorf("found existing go workflow file, use -o to overwrite")
		}
		actionsGo, err := os.Create(goyml)
		if err != nil {
			return err
		}
		defer actionsGo.Close()
		actionsGoTemplate := template.Must(template.New("actionsgo").
			Delims("[[", "]]").
			Parse(string(tpl.GoActionTemplate)))
		err = actionsGoTemplate.Execute(actionsGo, e)
		if err != nil {
			return err
		}
		SetKeyAndSaveConfig(keyActionsGo, true)
	}
	if e.ActionsPages {
		pagesYml := fmt.Sprintf("%s/.github/workflows/pages.yml", e.AbsolutePath)
		if exists(pagesYml) && !e.Overwrite {
			return fmt.Errorf("found existing pages workflow file, use -o to overwrite")
		}
		actionsPages, err := os.Create(pagesYml)
		if err != nil {
			return err
		}
		defer actionsPages.Close()
		actionsPagesTemplate := template.Must(template.New("actionspages").
			Delims("[[", "]]").
			Parse(string(tpl.ActionsPagesTemplate)))
		err = actionsPagesTemplate.Execute(actionsPages, e)
		if err != nil {
			return err
		}
		SetKeyAndSaveConfig(keyActionsPages, true)
	}
	if e.ActionsRelease {
		releaseYml := fmt.Sprintf("%s/.github/workflows/release.yml", e.AbsolutePath)
		if exists(releaseYml) && !e.Overwrite {
			return fmt.Errorf("found existing release workflow file, use -o to overwrite")
		}
		actionsRelease, err := os.Create(fmt.Sprintf("%s/.github/workflows/release.yml", e.AbsolutePath))
		if err != nil {
			return err
		}
		defer actionsRelease.Close()
		actionsReleaseTemplate := template.Must(template.New("actionsrelease").
			Delims("[[", "]]").
			Parse(string(tpl.ActionsReleaseTemplate)))
		err = actionsReleaseTemplate.Execute(actionsRelease, e)
		if err != nil {
			return err
		}
		SetKeyAndSaveConfig(keyActionsRelease, true)
	}
	return nil
}

func (p *Project) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create cmd/root.go
	if _, err = os.Stat(fmt.Sprintf("%s/cmd", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/cmd", p.AbsolutePath), 0751))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/cmd/root.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.RootTemplate())))
	err = rootTemplate.Execute(rootFile, p)
	if err != nil {
		return err
	}

	gitIgnore, err := os.Create(fmt.Sprintf("%s/.gitignore", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer gitIgnore.Close()
	gitIgnoreTemplate := template.Must(template.New("gitignore").Parse(string(tpl.GitIgnoreTemplate)))
	err = gitIgnoreTemplate.Execute(gitIgnore, p)
	if err != nil {
		return err
	}

	// create license
	return p.createLicenseFile()
}

func (p *Project) createLicenseFile() error {
	data := map[string]interface{}{
		"copyright": copyrightLine(),
	}
	licenseFile, err := os.Create(fmt.Sprintf("%s/LICENSE", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	licenseTemplate := template.Must(template.New("license").Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}

func (c *Command) Create() error {
	cmdFile, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", c.AbsolutePath, c.CmdName))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(template.New("sub").Parse(string(tpl.AddCommandTemplate())))
	err = commandTemplate.Execute(cmdFile, c)
	if err != nil {
		return err
	}
	return nil
}
func (c *Command) Docs() error {
	cmdFile, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", c.AbsolutePath, c.CmdName))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(template.New("sub").Parse(string(tpl.AddDocsTemplate())))
	err = commandTemplate.Execute(cmdFile, c)
	if err != nil {
		return err
	}
	err = fs.WalkDir(tpl.DocFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		os.CopyFS(".", tpl.DocFS)
		return nil
	})
	if err != nil {
		return err
	}
	modFile, err := os.Create(fmt.Sprintf("%s/docs/go.mod", c.AbsolutePath))
	if err != nil {
		return err
	}
	defer modFile.Close()
	docModTemplate := template.Must(template.New("docmod").Parse(string(tpl.GoModTemplate)))
	err = docModTemplate.Execute(modFile, c)
	if err != nil {
		return err
	}

	docsTask, err := os.Create(fmt.Sprintf("%s/docs/Taskfile.yml", c.AbsolutePath))
	if err != nil {
		return err
	}
	defer docsTask.Close()
	docsTaskTemplate := template.Must(template.New("docstask").Parse(string(tpl.DocsTaskfileTemplate)))
	err = docsTaskTemplate.Execute(docsTask, c)
	if err != nil {
		return err
	}
	return err
}
