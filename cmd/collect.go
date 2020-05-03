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

	"github.com/spf13/cobra"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run:   CollectProjects,
}

func init() {
	rootCmd.AddCommand(collectCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// TODO: Add doc
func CollectProjects(cmd *cobra.Command, args []string) {
	// TODO: make skip dirs configurable, abbreviated or absolute path
	// strings.HasPrefix(path, "~/")
	// path.IsAbs(*destination)
	skipDirs := map[string]bool{
		os.ExpandEnv("$HOME/Library"): true,
	}

	homeDir := HomeDir()
	projectListPath := fmt.Sprintf("%s/.pd_history", homeDir)

	file, err := os.Create(projectListPath)
	PanicIfErr(err)
	defer file.Close()

	walkErr := filepath.Walk(homeDir,
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

	PanicIfErr(walkErr)
}

// TODO: Add doc
func writeLogEntry(path string, file *os.File) {
	entry := ParseLogEntry(path)

	line := fmt.Sprintf(
		"1,%s,%s,%s\n",
		entry.projectName,
		entry.location,
		entry.absolutePath,
	)

	_, writeErr := file.WriteString(line)
	PanicIfErr(writeErr)

	syncErr := file.Sync()
	PanicIfErr(syncErr)
}
