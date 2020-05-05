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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var historyFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pd",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
}

// select - fuzzy-select a project / directory
// add - add project / directory to list
// sync - re-rank projects
// collect - find all VC projects in home directory
// directory - abs path from log entry label

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "~/.pdrc", "config file")
	rootCmd.PersistentFlags().StringVar(
		&historyFile, "history", "~/.pd_history", "history file")
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
