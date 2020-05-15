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
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// Expand the given path to an absolute path
// Expand placeholders like `.` and `~` to their canonical form.
// Return the home directory if given the empty string.
func expandPath(path string) string {
	path = filepath.Clean(path)
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", homeDir(), 1)
	}
	abspath, err := filepath.Abs(path)
	check(err)
	return abspath
}

// Return the current user's home directory.
func homeDir() string {
	home, err := homedir.Dir()
	check(err)
	return home
}

// Return the current user's current working directory.
func workingDir() string {
	pwd, err := os.Getwd()
	check(err)
	return pwd
}

// Return true if the given path is a project.
// (i.e., a directory under version control or a Projectile project)
func isProject(path string) bool {
	return isVersionControlled(path) || isProjectile(path)
}

// Return true if the given path is under version control.
// Currently only supports Git.
func isVersionControlled(path string) bool {
	return exists(filepath.Join(path, ".git/"))
}

// Return true if the given path is a projectile project.
func isProjectile(path string) bool {
	return exists(filepath.Join(path, ".projectile"))
}

// if there's an error, print it out and exit with a failure exit code.
func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// List the contents of `path` using exa.
func listFilesExa(path string, abbreviated string) (string, error) {
	fmt.Println(abbreviated)
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
	return capturedOutput(cmd)
}

// List the contents of `path` using ls.
func listFilesLs(path string, abbreviated string) (string, error) {
	fmt.Println(abbreviated)
	cmd := exec.Command(
		"ls",
		"--almost-all",
		"--color=always",
		"--group-directories-first",
		"--human-readable",
		"-l",
		path,
	)
	return capturedOutput(cmd)
}

// List the contents of `path` using tree.
func listFilesTree(path string) (string, error) {
	cmd := exec.Command(
		"tree",
		"-C",
		"-L",
		"2",
		path,
	)
	return capturedOutput(cmd)
}

// Capture the output of command `cmd`, returning a string.
func capturedOutput(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	var output string

	cmd.Stdout = &out
	err := cmd.Run()

	if err == nil {
		output = out.String()
	}

	return output, err
}
