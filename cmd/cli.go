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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO:
// eliminate lossy path / label parsing
// implement skip patterns
// expiration strategy
// profile & optimize

var cfgFile string
var historyFile string
var debug bool

var rootCmd = &cobra.Command{
	Use:   "pd",
	Short: "A project / directory manager and FZF-powered fuzzy-selector.",
	Run: func(cmd *cobra.Command, args []string) {
		target := strings.Trim(strings.Join(args, " "), " ")
		if len(target) == 0 {
			selectProject()
			return
		}

		if target == "--help" {
			fmt.Println(help)
			return
		}

		// Given a project label, generate a preview (file listing) of its contents.
		// Examples:
		// pd --pd-preview=my-project Documents/projects
		// pd --pd-preview=my-other-project
		if strings.HasPrefix(target, "--pd-preview=") {
			label := strings.Replace(target, "--pd-preview=", "", 1)
			listProjectFiles(label)
			return
		}

		// Find all version-controlled projects $HOME and refresh the history.
		// 	Refreshing the history removes any directories that no longer exist,
		// 	and re-aggregates and re-ranks entries.
		if target == "--pd-refresh" {
			collectProjects()
			refreshProjectListing()
			return
		}

		if strings.HasPrefix(target, "-") ||
			strings.HasPrefix(target, "+") {
			fmt.Println(target)
			return
		}

		dir := findDirectory(target)
		addLogEntry(dir)
		fmt.Println(dir)

		refreshProjectListing()
	},
	DisableFlagParsing: true,
	SilenceErrors:      true,
	SilenceUsage:       true,
}

var help = `
p/d

Usage:

  pd [directory name]

Intended to be used in tandem with cd as follows:

  cd $(pd ~/Documents)

Given a file path, print its absolute form (resolving symlinks) and save to
history. If a path to a non-directory is given, use its containing
directory instead.

Examples:

  pd ~/Documents/projects/my-project
  pd ~/my-other-project
  pd ./projects/my-project/some-file.txt

Given a position on the directory stack, no-op. Print that back out to leave the
behavior of cd unchanged.

Examples:
  pd -2
  pd +1

Given no arguments, open FZF to allow fuzzy-selecting a directory to cd into.
`

func init() {
	cobra.OnInitialize(initConfig)
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	check(err)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath(homeDir())
	viper.SetConfigName(".pdrc")
	viper.ReadInConfig()

	historyFile = expandPath("~/.pd_history")
}
