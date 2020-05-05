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
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	finder "github.com/b4b4r07/go-finder"
	"github.com/b4b4r07/go-finder/source"
	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "TODO: Add description",
	Long:  `TODO: Add description`,
	Run:   SelectProject,
}

func init() {
	rootCmd.AddCommand(selectCmd)
}

func SelectProject(cmd *cobra.Command, args []string) {
	fzf, err := finder.New(
		"fzf",
		"--ansi",
		"--bind 'ctrl-b:preview-up'",
		"--bind 'ctrl-f:preview-down'",
		"--cycle",
		"--no-multi",
		"--no-sort",
		"--preview='cd-dispatch-fzf-preview {+}'",
		"--reverse",
		"--tiebreak=index",
	)
	PanicIfErr(err)

	histExists, _ := Exists(historyFile)
	if !histExists {
		FindAndSaveProjects()
	}

	fzf.Read(HistoryFileSource())
	selection, err := fzf.Run()
	PanicIfErr(err)

	if len(selection) > 0 {
		fmt.Println(selection)
	}
}

// HistoryFileSource streams history file contents
func HistoryFileSource() source.Source {
	return func(out io.WriteCloser) error {
		fp, err := os.Open(historyFile)
		PanicIfErr(err)
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			entry := strings.Split(scanner.Text(), ",")
			fmt.Fprintln(out, entry[2])
		}
		return scanner.Err()
	}
}
