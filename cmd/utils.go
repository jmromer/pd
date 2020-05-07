/*
Copyright © 2020 Jake Romer <mail@jakeromer.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// Exists returns true if the given file or directory exists, else false.
// TODO: concentrate error checking here
func exists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// TODO: Add doc
func expandPath(path string) string {
	if len(path) == 0 {
		return homeDir()
	}

	// resolve to lexical file path
	path = filepath.Clean(path)
	if path == "." {
		pwd, err := os.Getwd()
		check(err)
		return pwd
	}

	// expand home directory if a tilde is passed
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", homeDir(), 1)
	}

	// expand to absolute path
	abspath, err := filepath.Abs(path)
	check(err)

	return abspath
}

// TODO: Add doc
func homeDir() string {
	home, err := homedir.Dir()
	check(err)
	return home
}

// TODO: Add doc
func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// list contents of `path` using exa
func listFilesExa(path string) (string, error) {
	var out bytes.Buffer
	var output string

	cmd := exec.Command(
		"exa",
		"--all",
		"--color=always",
		"--git",
		"--group-directories-first",
		"--header",
		"--long",
		path,
	)

	cmd.Stdout = &out
	err := cmd.Run()

	if err == nil {
		output = out.String()
	}

	return output, err
}

// list contents of `path` using ls
func listFilesLs(path string) (string, error) {
	var out bytes.Buffer
	var output string

	cmd := exec.Command(
		"ls",
		"--almost-all",
		"--color=always",
		"--group-directories-first",
		"--human-readable",
		"-l",
		path,
	)

	cmd.Stdout = &out
	err := cmd.Run()

	if err == nil {
		output = out.String()
	}

	return output, err
}
