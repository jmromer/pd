/*
Copyright © 2023 Jake Romer <jmromer@tensorconclave.com>

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
	"math"
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
		"--exact",
		"--no-multi",
		"--no-sort",
		`--preview="pd --fzf-preview "{+}""`,
		"--reverse",
		"--tiebreak=index",
	)
	check(err)

	if !exists(historyFile) {
		RefreshLog(true)
	}

	// projects: maps abspaths to LogEntries
	projects := currentlyLoggedProjects()

	// bail if no projects logged
	if len(projects) == 0 {
		fmt.Println(workingDir())
		return
	}

	listingEntries, listingIndex := searchListing(projects)
	fzf.Read(listingEntries)

	selection, err := fzf.Run()
	check(err)

	// bail if selection is canceled
	if len(selection) == 0 {
		fmt.Println(workingDir())
		return
	}

	// the selected label is stripped of ansi color codes
	// use listingIndex to retrieve the associated abspath
	abspath := listingIndex[selection[0]]
	fmt.Println(abspath)
	addLogEntry(abspath)

	RefreshLog(false)
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
	abspath := projectLabelToAbsPath(label)
	abbreviated := strings.Replace(abspath, homeDir(), "~", 1)
	list, err := listFilesExa(abspath, abbreviated)

	switch {
	case err != nil:
		list, err = listFilesTree(abspath)
		fallthrough
	case err != nil:
		list, err = listFilesLs(abspath, abbreviated)
		fallthrough
	case err == nil && len(list) > 0:
		fmt.Println(list)
	case len(list) == 0:
		fmt.Println("Empty")
	case !exists(abspath):
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
	sort.Sort(ByName(entries))
	sort.Sort(ByCount(entries))

	writeLogEntries(entries)
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
// Skip if the given path is the home directory, since we always
// display the home directory first in the listing.
func addLogEntry(abspath string) {
	if abspath == homeDir() {
		return
	}
	f, err := os.OpenFile(
		expandPath(historyFile),
		os.O_APPEND|os.O_WRONLY,
		0600,
	)
	check(err)
	defer f.Close()
	entry := buildLogEntry(abspath)
	entry.WriteLogLine(f)
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
		projectName := line[2]
		projectPath := line[3]

		entry, isAlreadyCounted := entries[abspath]
		if isAlreadyCounted {
			currCount += entry.Count
		}

		entries[abspath] = LogEntry{
			Count:   currCount,
			AbsPath: abspath,
			Name:    projectName,
			Path:    projectPath,
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
func writeLogEntries(entries []LogEntry) {
	if debug {
		fmt.Println("Refreshing project listing...")
	}

	// create history file, overwriting if need be
	f, err := os.Create(expandPath(historyFile))
	check(err)
	defer f.Close()

	for _, entry := range entries {
		entry.WriteLogLine(f)
	}

	if debug {
		fmt.Println("Completed at", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	}
}

// Each entry in the pd history file consists of a count, an abs path, and the
// colored project label to be used in the FZF interface.
type LogEntry struct {
	Count   int
	AbsPath string
	Name    string
	Path    string
}

// Project label, formatted (ansi-color)
func (e LogEntry) LabelFormatted() string {
	name := aurora.Blue(e.Name).String()
	path := aurora.Gray(12-1, e.Path).String()
	elts := []string{name, path}
	return strings.Join(elts, " ")
}

// Project label, unformatted
func (e LogEntry) Label() string {
	return strings.Join([]string{e.Name, e.Path}, " ")
}

// History entry format is CSV with
// `count`, `absolute path`, and `project label`
func (e LogEntry) LogLine() string {
	return fmt.Sprintf("%d,%s,%s,%s\n", e.Count, e.AbsPath, e.Name, e.Path)
}

// Write the given LogEntry to the given file handle.
func (e LogEntry) WriteLogLine(file *os.File) {
	if exists(e.AbsPath) {
		_, err := file.WriteString(e.LogLine())
		check(err)
	}
}

// Sorting interfaces
type ByCount []LogEntry

func (a ByCount) Len() int      { return len(a) }
func (a ByCount) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool {
	if a[j].Count == a[i].Count {
		return strings.ToLower(a[i].Name) < strings.ToLower(a[j].Name)
	}
	return a[j].Count < a[i].Count
}

type ByName []LogEntry

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	if a[i].Name == a[j].Name {
		return a[i].Path < a[j].Path
	}
	return strings.ToLower(a[i].Name) < strings.ToLower(a[j].Name)
}

// Given an absolute path, parse out a project label and return a new LogEntry.
func buildLogEntry(abspath string) LogEntry {
	path := strings.Replace(abspath, homeDir(), "~", -1)
	components := strings.Split(path, "/")

	last := len(components) - 1
	location := strings.Join(components[0:last], "/")

	return LogEntry{
		Count:   1,
		AbsPath: abspath,
		Name:    components[last],
		Path:    location,
	}
}

// Build a LogEntry for the home directory
// This is used to represent the home directory in the FZF interface.
// Assign it the maximum possible count so that it always appears first.
func buildHomeLogEntry() LogEntry {
	return LogEntry{
		Count:   math.MaxInt32,
		AbsPath: homeDir(),
		Name:    "~",
		Path:    "",
	}
}

// Given a project label, re-construct the absolute path that was used to
// generate it.
func projectLabelToAbsPath(label string) string {
	label = strings.Trim(label, " ")

	if label == "~" {
		return homeDir()
	}

	comps := strings.Split(label, " ~")
	if len(comps) > 1 {
		projName := comps[0]
		pathLabel := comps[1]
		path := fmt.Sprintf("%s/%s", homeDir(), pathLabel)
		return filepath.Join(path, projName)
	}

	comps = strings.Split(label, " /")
	if len(comps) > 1 {
		projName := comps[0]
		pathLabel := comps[1]
		path := fmt.Sprintf("/%s", pathLabel)
		return filepath.Join(path, projName)
	}

	return ""
}

// Build an FZF listing and a listing index
//
// The `listing` is a slice of formatted labels (ansi-colored)
// The `index` maps labels (without color codes) to abs paths.
//
// Return:
// (0) a Source object to be passed to a finder's Read method
// (1) the `index` mapping
func searchListing(projectIndex map[string]LogEntry) (source.Source, map[string]string) {
	logEntries := []LogEntry{}

	logEntries = append(logEntries, buildHomeLogEntry())
	for _, logEntry := range projectIndex {
		logEntries = append(logEntries, logEntry)
	}

	sort.Sort(ByName(logEntries))
	sort.Sort(ByCount(logEntries))

	listing := []string{}
	index := map[string]string{}

	for _, logEntry := range logEntries {
		index[logEntry.Label()] = logEntry.AbsPath
		listing = append(listing, logEntry.LabelFormatted())
	}

	return source.Slice(listing), index
}
