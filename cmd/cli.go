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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
  pd -

Given no arguments, open FZF to allow fuzzy-selecting a directory to cd into.
`

// WIP:
// implement skip patterns
// eliminate lossy path / label parsing
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
		dirStackFlag := regexp.MustCompile("\\A[-+][0-9]*\\z")

		switch {
		case len(target) == 0:
			SelectProject()
		case target == "--help":
			fmt.Println(help)
		case strings.HasPrefix(target, "--fzf-preview"):
			label := strings.Replace(target, "--fzf-preview", "", 1)
			FzfPreview(strings.Trim(label, " "))
		case target == "--pd-refresh":
			RefreshLog(true)
		case dirStackFlag.MatchString(target):
			fmt.Println(target)
		default:
			ChangeDirectory(target)
		}
	},
	DisableFlagParsing: true,
	SilenceErrors:      true,
	SilenceUsage:       true,
}

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
