/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}