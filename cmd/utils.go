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
	"fmt"
	"os"
	"strings"

	. "github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
)

type LogEntry struct {
	projectName  string
	location     string
	absolutePath string
}

// ParseLogEntry parses the given absolute path `abspath`, returning an LogEntry
// struct with `projectName`, `location`, and `absolutePath`
func ParseLogEntry(abspath string) LogEntry {
	homeDir := fmt.Sprintf("%s/", os.ExpandEnv("$HOME"))
	path := strings.Replace(abspath, homeDir, "", -1)
	components := strings.Split(path, "/")

	last := len(components) - 1
	location := strings.Join(components[0:last], "/")
	projectName := components[last]

	return LogEntry{
		projectName:  Blue(projectName).String(),
		location:     Gray(12-1, location).String(),
		absolutePath: abspath,
	}
}

// Exists returns true if the given file or directory exists, else false.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// TODO: Add doc
func PanicIfErr(e error) {
	if e != nil {
		panic(e)
	}
}

// TODO: Add doc
func HomeDir() string {
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home
}
