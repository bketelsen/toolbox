# yaml-language-server: $schema=https://goreleaser.com/static/schema-pro.json
# vim: set ts=2 sw=2 tw=0 fo=cnqoj

version: 2
project_name: [[ .AppName ]] 
before:
  hooks:
    # You may remove this if you don't use go modules.
    - go mod tidy
after:
  hooks:
    - go run main.go docs
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
    goarch:
      - amd64
      - arm64
    binary: [[ .AppName ]] 
    id: [[ .AppName ]] 
    # Default: '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser'.
    # Templates: allowed.
    ldflags:
      - '-s -w -X [[ .PkgName ]]/cmd.version={{.Version}} -X [[ .PkgName ]]/cmd.commit={{.Commit}}'

archives:
  - formats: [ 'tar.gz' ]
    # this name template makes the OS and Arch compatible with the results of `uname`.
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    # use zip for windows archives
    format_overrides:
      - goos: windows
        formats: [ 'zip' ]

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"

release:
  footer: >-

    ---

    Released by [GoReleaser](https://github.com/goreleaser/goreleaser).
