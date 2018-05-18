package conflict

import (
	"errors"
	"path"
	"strconv"
	"strings"
)

func parseGitMarkerInfo(diff string) (fname string, lineNum int, err error) {
	parts := strings.Split(diff, ":")

	if len(parts) < 3 || !strings.Contains(diff, "leftover conflict marker") {
		err = errors.New("Line does not contain marker location info")
		return
	}

	fname, lineData := string(parts[0]), parts[1]
	if lineNum, err = strconv.Atoi(string(lineData)); err != nil {
		return
	}
	return
}

func parseConflictsIn(f File, markerLocations []int) (conflicts []Conflict, err error) {
	var conf Conflict

	for _, lineNum := range markerLocations {
		line := f.Lines[lineNum-1]
		index := lineNum - 1

		switch IdentifyStyle(line) {
		case start:
			conf = Conflict{}
			conf.Start = index
		case separator:
			conf.Middle = index
		case diff3:
			conf.Diff3 = append(conf.Diff3, index)
		case end:
			conf.End = index
			conflicts = append(conflicts, conf)
		default:
			continue
		}
	}

	for i := range conflicts {
		c := &conflicts[i]
		c.File = &f

		if !(c.Valid()) {
			return nil, errors.New("Invalid number of remaining conflict markers")
		}

		if err = c.Extract(f.Lines); err != nil {
			return
		}

		if err = c.HighlightSyntax(); err != nil {
			return
		}
	}

	return
}

// Find runs `git --no-pager diff --check` in order to detect git conflicts
// It returns an array of `File`s where each `File` contains conflicts within itself
// If the parsing fails, it returns an error
func Find(cwd string) (files []File, err error) {
	topPath, err := TopLevelPath(cwd)
	if err != nil {
		return
	}

	markerLocations, err := MarkerLocations(topPath)
	if err != nil {
		return
	}

	markerLocMap := make(map[string][]int)
	for _, line := range markerLocations {
		if len(line) == 0 {
			continue
		}

		fname, line, ok := parseGitMarkerInfo(line)
		if ok != nil {
			continue
		}
		markerLocMap[fname] = append(markerLocMap[fname], line)
	}

	for fname := range markerLocMap {
		absPath := path.Join(topPath, fname)
		file := File{Name: fname, AbsolutePath: absPath}

		if err = file.Read(); err != nil {
			return
		}

		conflicts, err := parseConflictsIn(file, markerLocMap[fname])
		if err != nil {
			return nil, err
		}

		file.Conflicts = conflicts
		files = append(files, file)
	}
	return
}
