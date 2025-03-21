# https://taskfile.dev
# Generated File, changes may be lost
# Add `Taskfile.custom.yml` in this directory with your additions

version: '3'

vars:
  VERSION: 0.0.1

includes:
  docs:
    taskfile: Taskfile.docs.yml
    optional: true
  checks:
    taskfile: Taskfile.checks.yml
    optional: true
  release:
    taskfile: Taskfile.release.yml
    optional: true
  custom:
    taskfile: Taskfile.custom.yml
    optional: true

tasks:
  build:
    desc: Build the application
    summary: |
      Build the application with ldflags to set the version with a -dev suffix.

      Output: '[[ .AppName ]]' in project root.
    cmds:
      - go build -o [[ .AppName ]] -ldflags '-s -w -X [[ .PkgName ]]/cmd.version={{.VERSION}}-dev' main.go
    silent: true

  tools:
    desc: Install required tools
    cmds:
      - go install github.com/bketelsen/toolbox/toolgen@latest

  direnv:
    desc: Add direnv hook to your bashrc
    cmds:
      - direnv hook bash >> ~/.bashrc
    silent: true
[[ if .Config.GetBool "docs" ]]
  generate:
    desc: Generate CLI documentation
    deps: [tools]
    cmds:
      - go run main.go gendocs -b "[[ .Config.GetString "basepath" ]]"
    silent: true
[[ end ]]
