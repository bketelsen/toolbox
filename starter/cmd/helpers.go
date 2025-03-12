// Copyright © 2015 Steve Francia <spf@spf13.com>.
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
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cobra "github.com/bketelsen/toolbox/cobra"
)

var srcPaths []string

func init() {
	// Initialize srcPaths.
	envGoPath := os.Getenv("GOPATH")
	goPaths := filepath.SplitList(envGoPath)
	if len(goPaths) == 0 {
		// Adapted from https://github.com/Masterminds/glide/pull/798/files.
		// As of Go 1.8 the GOPATH is no longer required to be set. Instead there
		// is a default value. If there is no GOPATH check for the default value.
		// Note, checking the GOPATH first to avoid invoking the go toolchain if
		// possible.

		goExecutable := os.Getenv("COBRA_GO_EXECUTABLE")
		if len(goExecutable) <= 0 {
			goExecutable = "go"
		}

		out, err := exec.Command(goExecutable, "env", "GOPATH").Output()
		cobra.CheckErr(err)

		toolchainGoPath := strings.TrimSpace(string(out))
		goPaths = filepath.SplitList(toolchainGoPath)
		if len(goPaths) == 0 {
			cobra.CheckErr("$GOPATH is not set")
		}
	}
	srcPaths = make([]string, 0, len(goPaths))
	for _, goPath := range goPaths {
		srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		cobra.CheckErr(err)
	}
	return true
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

type TaskSummary struct {
	Tasks []struct {
		Name     string `json:"name,omitempty"`
		Desc     string `json:"desc,omitempty"`
		Summary  string `json:"summary,omitempty"`
		Aliases  []any  `json:"aliases,omitempty"`
		UpToDate bool   `json:"up_to_date,omitempty"`
		Location struct {
			Line     int    `json:"line,omitempty"`
			Column   int    `json:"column,omitempty"`
			Taskfile string `json:"taskfile,omitempty"`
		} `json:"location,omitempty"`
	} `json:"tasks,omitempty"`
	Location string `json:"location,omitempty"`
}

// getTaskSummary returns a TaskSummary struct from the Taskfile.yml
func getTaskSummary(path string) (TaskSummary, error) {
	var taskSummary TaskSummary
	taskfile := filepath.Join(path, "Taskfile.yml")
	if !exists(taskfile) {
		cobra.CheckErr("Taskfile.yml not found")
	}
	out, err := exec.Command("task", "--list-all", "--json").Output()
	if err != nil {
		return taskSummary, err
	}

	err = json.Unmarshal(out, &taskSummary)

	return taskSummary, err
}
