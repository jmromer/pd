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
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
)

type LogEntry struct {
	count        int
	projectName  string
	location     string
	absolutePath string
}

// NewLogEntry parses the given absolute path `abspath`, returning an LogEntry
// struct with `projectName`, `location`, and `absolutePath`
func NewLogEntry(abspath string) LogEntry {
	homeDir := fmt.Sprintf("%s/", HomeDir())
	path := strings.Replace(abspath, homeDir, "", -1)
	components := strings.Split(path, "/")

	last := len(components) - 1
	location := strings.Join(components[0:last], "/")
	projectName := components[last]

	return LogEntry{
		count:        1,
		projectName:  aurora.Blue(projectName).String(),
		location:     aurora.Gray(12-1, location).String(),
		absolutePath: abspath,
	}
}

// TODO: Add doc
func FindAndSaveProjects() {
	fmt.Println("Finding project directories...")
	// TODO: make skip dirs configurable, abbreviated or absolute path
	// strings.HasPrefix(path, "~/")
	// path.IsAbs(*destination)
	skipDirs := map[string]bool{
		os.ExpandEnv("$HOME/Library"): true,
	}

	file, err := os.Create(historyFile)
	PanicIfErr(err)
	defer file.Close()

	homeDir := HomeDir()
	err = filepath.Walk(homeDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() {
				return nil
			}

			if strings.HasPrefix(info.Name(), ".") || skipDirs[path] {
				return filepath.SkipDir
			}

			isGitProj, err := Exists(filepath.Join(path, ".git"))
			if err == nil && isGitProj {
				writeLogEntry(path, file)
				return filepath.SkipDir
			}

			return nil
		})

	PanicIfErr(err)
}

// TODO: Add doc
// history entry format:
// count, abspath, history log entry
func writeLogEntry(path string, file *os.File) {
	entry := NewLogEntry(path)
	line := fmt.Sprintf(
		"%d,%s,%s %s\n",
		entry.count,
		entry.absolutePath,
		entry.projectName,
		entry.location,
	)

	_, err := file.WriteString(line)
	PanicIfErr(err)

	err = file.Sync()
	PanicIfErr(err)
}
