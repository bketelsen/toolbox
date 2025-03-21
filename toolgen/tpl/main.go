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

package tpl

import (
	"embed"
)

//go:embed all:site/*
var SiteFS embed.FS

//go:embed all:content/*
var ContentFS embed.FS

//go:embed go.mod.tpl
var GoModTemplate []byte

//go:embed hugo.mod.tpl
var HugoModTemplate []byte

//go:embed hugo.yaml.tpl
var HugoYamlTemplate []byte

//go:embed go.yml.tpl
var GoActionTemplate []byte

//go:embed Taskfile.yml.tpl
var TaskfileTemplate []byte

//go:embed Taskfile.checks.yml.tpl
var TaskfileChecksTemplate []byte

//go:embed Taskfile.release.yml.tpl
var TaskfileReleaseTemplate []byte

//go:embed Taskfile.docs.yml.tpl
var DocsTaskfileTemplate []byte

//go:embed .goreleaser.yaml.tpl
var GoReleaserTemplate []byte

//go:embed pages.yml.tpl
var ActionsPagesTemplate []byte

//go:embed release.yml.tpl
var ActionsReleaseTemplate []byte

//go:embed devcontainer.json.tpl
var DevContainerTemplate []byte

//go:embed Dockerfile.tpl
var DockerfileTemplate []byte

//go:embed gitignore.tpl
var GitIgnoreTemplate []byte

//go:embed install.sh.tpl
var InstallScriptTemplate []byte

//go:embed main.go.tpl
var MainTemplate []byte

//go:embed root.go.tpl
var RootTemplate []byte

//go:embed command.go.tpl
var AddCommandTemplate []byte

//go:embed docgen.go.tpl
var AddDocsTemplate []byte

func TaskSummaryTemplate() []byte {
	return []byte(`
## Available Tasks
{{ range .Summary.Tasks }}
### {{ .Name }}

{{ if .Desc }}Description: {{  .Desc }}
{{- end }}
{{ if .Summary }}Summary: {{ .Summary }}
{{- end }}
Run this task:
` + "```" + `
task {{ .Name }}
` + "```" + `
{{ end }}
`)
}
