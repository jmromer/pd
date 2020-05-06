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

// TODO;
// evacuation strategy (from original)
// internalize preview logic
// documentation
// function name capitalization conventions
// skip patterns
// profile & optimize
// customize fzf command
// skip -3 args?

var cfgFile string
var historyFile string
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "pd",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		SelectProject()
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}
		AddLogEntry(args[0])
	},
}

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		path := strings.Trim(strings.Join(args, " "), " ")
		if strings.HasPrefix(path, "-") {
			fmt.Println(path)
			return
		} else {
			fmt.Println(ChangeDirectory(path))
			SyncProjectListing()
		}
	},
}

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		CollectProjects()
	},
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		SelectProject()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		SyncProjectListing()
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	check(err)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "~/.pdrc", "config file")
	rootCmd.PersistentFlags().StringVarP(
		&historyFile, "file", "f", "~/.pd_history", "history file")
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(selectCmd)
	rootCmd.AddCommand(syncCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Use config file from the flag.
	viper.SetConfigFile(ExpandPath(cfgFile))
	// read in environment variables that match
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	viper.ReadInConfig()
	historyFile = ExpandPath(historyFile)
}
