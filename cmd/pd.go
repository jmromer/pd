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
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	finder "github.com/b4b4r07/go-finder"
	"github.com/b4b4r07/go-finder/source"
	"github.com/logrusorgru/aurora"
)

// Use FZF to select a project directory, printing to stdout the absolute path
// to the project directory.
func SelectProject() {
	fzf, err := finder.New(
		"fzf",
		"--ansi",
		"--bind 'ctrl-b:preview-up'",
		"--bind 'ctrl-f:preview-down'",
		"--cycle",
		"--no-multi",
		"--no-sort",
		"--preview='pd --fzf-preview {+}'",
		"--reverse",
		"--tiebreak=index",
	)
	check(err)

	if !exists(historyFile) {
		RefreshLog(true)
	}

	fzf.Read(historyFileSource())
	selection, err := fzf.Run()
	check(err)

	if len(selection) > 0 {
		abspath := projectLabelToAbsPath(selection[0])
		fmt.Println(abspath)
		addLogEntry(abspath)
		RefreshLog(false)
	} else {
		fmt.Println(workingDir())
	}
}

// FzfPreview triggers a preview (file listing) of the directory associated with
// the given project label.
//
// List the files in the directory associated with the given project label.
// Used in the FZF preview window. Attempts to do this using exa, then tree if
// exa is unavailable or fails. Falls back to ls if both fail.
//
// Examples:
// pd --fzf-preview my-project Documents/projects
// pd --fzf-preview my-other-project
func FzfPreview(label string) {
	path := projectLabelToAbsPath(label)
	abbreviated := strings.Replace(path, homeDir(), "~", 1)
	list, err := listFilesExa(path, abbreviated)

	switch {
	case err != nil:
		list, err = listFilesTree(path)
		fallthrough
	case err != nil:
		list, err = listFilesLs(path, abbreviated)
		fallthrough
	case err == nil && len(list) > 0:
		fmt.Println(list)
	case len(list) == 0:
		fmt.Println("Empty")
	case !exists(path):
		fmt.Println("Directory does not exist.")
	default:
		fmt.Println("Could not list contents.")
	}
}

// RefreshLog finds all version-controlled projects $HOME and refresh the
// history. Refreshing the history removes any directories that no longer exist,
// and re-aggregates and re-ranks entries.
func RefreshLog(searchForProjects bool) {
	var projects []string

	if searchForProjects {
		// scan home directory for all project paths
		projects = collectUserProjects()
	} else {
		projects = []string{}
	}

	// retrieve current log entries
	logEntries := currentlyLoggedProjects()

	// keep only those found projects not currently in the log
	entries := collectEntries(projects, logEntries)

	refreshProjectListing(entries)
}

// ChangeDirectory resolves the given path to an absolute path
// Logs an entry to the history log and refreshes the project listing
func ChangeDirectory(target string) {
	projectPath := findDirectory(target)
	fmt.Println(projectPath)
	addLogEntry(projectPath)
	RefreshLog(false)
}

// Append a log entry for the the given absolute path to the pd history file
func addLogEntry(abspath string) {
	f, err := os.OpenFile(
		expandPath(historyFile),
		os.O_APPEND|os.O_WRONLY,
		0600,
	)
	check(err)
	defer f.Close()
	writeLogEntry(buildLogEntry(abspath), f)
}

// Resolve the given path to an absolute path to a directory.
// If given a path to a file, return the containing directory.
func findDirectory(path string) string {
	target := expandPath(path)

	// use the containing directory if `path` is a file
	if stat, err := os.Stat(target); err == nil && !stat.IsDir() {
		target = filepath.Dir(target)
	}

	return target
}

// Scan the current user's home directory for projects
//
// "Projects" being directories under version control or that include
// .projectile file.
//
// For a given file path, skip to the next iteration if:
//
//   1. The file isn't a directory
//   2. The directory begins with a `.`
//   3. The directory is in the given skip list
//
// Otherwise, collect the directory and skip to the next iteration without
// recursing into the directory (nested projects won't be logged unless manually
// visited).
func collectUserProjects() []string {
	if debug {
		fmt.Println("Finding project directories...")
	}
	projects := []string{}
	skipDirs := map[string]bool{
		// TODO: Add config file knob for this
		os.ExpandEnv("$HOME/Library"): true,
	}

	filepath.Walk(
		homeDir(),
		func(path string, info os.FileInfo, err error) error {
			switch {
			case err != nil:
				return nil
			case !info.IsDir():
				// if the given file isn't a directory, we can skip it
				return nil
			case strings.HasPrefix(info.Name(), "."):
				// if the given directory is a dotfile directory
				return filepath.SkipDir
			case skipDirs[path]:
				// if the given directory is in the skip list
				return filepath.SkipDir
			case isProject(path):
				// if the given directory is a project,
				// log its path and don't recurse into it.
				projects = append(projects, path)
				return filepath.SkipDir
			default:
				return nil
			}
		},
	)

	return projects
}

func currentlyLoggedProjects() map[string]LogEntry {
	entries := make(map[string]LogEntry)

	if !exists(historyFile) {
		return entries
	}

	fp, err := os.Open(historyFile)
	check(err)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		currCount, _ := strconv.Atoi(line[0])
		abspath := line[1]
		label := line[2]

		entry, isAlreadyCounted := entries[abspath]
		if isAlreadyCounted {
			currCount += entry.Count
		}

		entries[abspath] = LogEntry{
			Count:        currCount,
			AbsolutePath: abspath,
			Label:        label,
		}
	}
	return entries
}

func collectEntries(foundPaths []string, currEntries map[string]LogEntry) (entries []LogEntry) {
	// Keep all current entries
	for _, entry := range currEntries {
		entries = append(entries, entry)
	}
	// Keep new entries
	for _, abspath := range foundPaths {
		_, isAlreadyLogged := currEntries[abspath]
		if !isAlreadyLogged {
			entries = append(entries, buildLogEntry(abspath))
		}
	}
	return
}

// Refresh the pd history file
// Re-aggregate and re-rank entries, remove directories that no longer exist.
func refreshProjectListing(entries []LogEntry) {
	if debug {
		fmt.Println("Refreshing project listing...")
	}

	// aggregate log entries, sorting by count in desc order
	sort.Sort(ByCount(entries))

	// write sorted entries to log
	f, err := os.Create(expandPath(historyFile))
	check(err)
	defer f.Close()
	for _, entry := range entries {
		if exists(entry.AbsolutePath) {
			writeLogEntry(entry, f)
		}
	}

	if debug {
		fmt.Println("Completed at", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	}
}

// Each entry in the pd history file consists of a count, an abs path, and the
// colored project label to be used in the FZF interface.
type LogEntry struct {
	Count        int
	AbsolutePath string
	Label        string
}

type ByCount []LogEntry

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[j].Count < a[i].Count }

// Given an absolute path, parse out a project label and return a new LogEntry.
func buildLogEntry(abspath string) LogEntry {
	homeDir := fmt.Sprintf("%s/", homeDir())
	path := strings.Replace(abspath, homeDir, "", -1)
	components := strings.Split(path, "/")

	last := len(components) - 1
	location := strings.Join(components[0:last], "/")
	label := fmt.Sprintf(
		"%s %s",
		aurora.Blue(components[last]).String(),
		aurora.Gray(12-1, location).String(),
	)

	return LogEntry{Count: 1, AbsolutePath: abspath, Label: label}
}

// Stream the history file's contents
func historyFileSource() source.Source {
	return func(out io.WriteCloser) error {
		fp, err := os.Open(historyFile)
		check(err)
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			entry := strings.Split(scanner.Text(), ",")
			abspath := entry[1]
			label := entry[2]
			if exists(abspath) {
				fmt.Fprintln(out, label)
			}
		}
		return scanner.Err()
	}
}

// Given a project label, re-construct the absolute path that was used to
// generate it.
//
// TODO: Find a better way to do this, since it's a lossy process and leaves
// some un-fixable corner cases. Ideally we want to display the project label
// when fuzzy-selecting but return the associated absolute path without parsing
// it out.
func projectLabelToAbsPath(label string) string {
	comps := strings.Split(label, " ")
	proj := comps[0]
	abspath := homeDir()

	if len(comps) > 1 {
		path := strings.Join(comps[1:], " ")
		if strings.HasPrefix(path, "/") {
			abspath = path
		} else {
			abspath = filepath.Join(abspath, path)
		}
	}

	return filepath.Join(abspath, proj)
}

// Write the given LogEntry to the given file handle.
// History entry format is CSV with
// `count`, `absolute path`, and `project label`
func writeLogEntry(entry LogEntry, file *os.File) {
	line := fmt.Sprintf(
		"%d,%s,%s\n",
		entry.Count,
		entry.AbsolutePath,
		entry.Label,
	)
	_, err := file.WriteString(line)
	check(err)
}
