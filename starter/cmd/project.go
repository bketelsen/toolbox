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
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
}

type Command struct {
	CmdName   string
	CmdParent string
	*Project
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
		fmt.Printf("path=%q, isDir=%v\n", path, d.IsDir())
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
	return err
}
