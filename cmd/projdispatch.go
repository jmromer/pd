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

// TODO: Add doc
// append to log file
func AddLogEntry(abspath string) {
	f, err := os.OpenFile(
		ExpandPath(historyFile),
		os.O_APPEND|os.O_WRONLY,
		0600,
	)
	check(err)
	defer f.Close()
	writeLogEntry(buildLogEntry(abspath), f)
}

// TODO: Add doc
func ChangeDirectory(path string) string {
	if len(path) == 0 || path == "~" {
		return HomeDir()
	}

	// resolve symlinks in case path contains one
	target, err := filepath.EvalSymlinks(ExpandPath(path))
	check(err)

	// use the containing directory if `path` is a file
	if stat, err := os.Stat(target); err == nil && !stat.IsDir() {
		target = filepath.Dir(target)
	}

	AddLogEntry(target)
	return target
}

// TODO: Add doc
func CollectProjects() {
	fmt.Println("Finding project directories...")
	skipDirs := map[string]bool{
		os.ExpandEnv("$HOME/Library"): true,
	}

	file, err := os.Create(historyFile)
	check(err)
	defer file.Close()

	filepath.Walk(
		HomeDir(),
		func(path string, info os.FileInfo, err error) error {
			return collectProjectDir(path, skipDirs, info, file, err)
		},
	)
	file.Sync()
}

// TODO: Add doc
func SelectProject() {
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
	check(err)

	if !Exists(historyFile) {
		CollectProjects()
	}

	fzf.Read(historyFileSource())
	selection, err := fzf.Run()
	check(err)

	if len(selection) > 0 {
		fmt.Println(projectLabelToAbsPath(selection[0]))
	}
}

// TODO: Add doc
// read history entries
// re-rank them, aggregating multiple entries
func SyncProjectListing() {
	if verbose {
		fmt.Println("Syncing project listing...")
	}

	entryLabels := map[string]string{}
	entryCounts := map[string]int{}

	fp, err := os.Open(historyFile)
	check(err)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		entry := strings.Split(scanner.Text(), ",")
		count := 0
		count, _ = strconv.Atoi(entry[0])
		abspath := entry[1]
		label := entry[2]

		if Exists(abspath) {
			entryLabels[abspath] = label
			entryCounts[abspath] += count
		}
	}

	// aggregate log entries, sorting by count in desc order
	i := 0
	entries := make([]LogEntry, len(entryCounts))
	for path, ct := range entryCounts {
		entries[i] = LogEntry{Count: ct, AbsolutePath: path, Label: entryLabels[path]}
		i += 1
	}
	sort.Sort(ByCount(entries))

	// write sorted entries to log
	f, err := os.Create(ExpandPath(historyFile))
	check(err)
	defer f.Close()
	for _, entry := range entries {
		writeLogEntry(entry, f)
	}

	if verbose {
		fmt.Println("Synced at", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	}
}

type LogEntry struct {
	Count        int
	AbsolutePath string
	Label        string
}

type ByCount []LogEntry

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[j].Count < a[i].Count }

// TODO: add doc
func buildLogEntry(abspath string) LogEntry {
	homeDir := fmt.Sprintf("%s/", HomeDir())
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

// TODO: add doc
func collectProjectDir(path string, skipDirs map[string]bool, info os.FileInfo, file *os.File, err error) error {
	// If the given file isn't a directory, we can skip it
	if err != nil || !info.IsDir() {
		return nil
	}

	// if the given directory is skippable (or a .dotfile directory)
	// skip it, don't recurse into it
	if strings.HasPrefix(info.Name(), ".") || skipDirs[path] {
		return filepath.SkipDir
	}

	// if the given directory contains a Git repo,
	// log its path and don't recurse into it.
	isGitProj := Exists(filepath.Join(path, ".git"))
	if err == nil && isGitProj {
		writeLogEntry(buildLogEntry(path), file)
		return filepath.SkipDir
	}

	return nil
}

// stream history file contents
func historyFileSource() source.Source {
	return func(out io.WriteCloser) error {
		fp, err := os.Open(historyFile)
		check(err)
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			entry := strings.Split(scanner.Text(), ",")
			fmt.Fprintln(out, entry[2])
		}
		return scanner.Err()
	}
}

// TODO: Add doc
func projectLabelToAbsPath(selection string) string {
	comps := strings.Fields(selection)
	proj := comps[0]
	abspath := HomeDir()

	if len(comps) > 1 {
		path := comps[1]
		if strings.HasPrefix(path, "/") {
			abspath = path
		} else {
			abspath = strings.Join([]string{abspath, path}, "/")
		}
	}

	return strings.Join([]string{abspath, proj}, "/")
}

// TODO: Add doc
// history entry format: `count, abspath, history log entry`
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
