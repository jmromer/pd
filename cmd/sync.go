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

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync called")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
