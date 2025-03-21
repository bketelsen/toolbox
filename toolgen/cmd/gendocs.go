/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/cobra/doc"
)

// docsCmd represents the docs command
var gendocsCmd = &cobra.Command{
	Use:    "gendocs",
	Hidden: true,
	Short:  "Generates documentation for the toolgen project",
	Run: func(cmd *cobra.Command, args []string) {
		bp := cmd.GlobalConfig().GetString("basepath")
		cmd.Logger.Info("Base path for documentation", "basepath", bp)
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return bp + "/docs/toolgen/cli/" + strings.ToLower(base) + "/"
		}
		filePrepender := func(filename string) string {
			now := time.Now().Format(time.RFC3339)
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/docs/toolgen/cli/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}
		err := os.MkdirAll("./docs/content/docs/toolgen/cli/", 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = doc.GenMarkdownTreeCustom(rootCmd, "./docs/content/docs/toolgen/cli/", filePrepender, linkHandler)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gendocsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	gendocsCmd.Flags().StringP("basepath", "b", "/toolbox", "Base path for the documentation (default is /toolbox)")

}

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`
