package conflict

import (
	"errors"
	"path"
	"strings"
)

// Supported git conflict styles
const (
	text = iota
	start
	diff3
	separator
	end
)

// identifyStyle identifies the conflict marker style of provided text
func identifyStyle(line string) (style int) {
	line = strings.TrimSpace(line)

	if strings.Contains(line, "<<<<<<<") {
		style = start
	} else if strings.Contains(line, ">>>>>>>") {
		style = end
	} else if line == "||||||| merged common ancestors" {
		style = diff3
	} else if line == "=======" {
		style = separator
	} else {
		style = text
	}
	return
}

// GroupConflictMarkers groups the provided lines of a file into individual `Conflict` structs
// It returns an error if the line contains invalid number of conflicts.
// This may be caused if the user manually edits the file.
func GroupConflictMarkers(lines []string) (conflicts []Conflict, err error) {
	conf := Conflict{Start: -1}
	expected := start

	for idx, line := range lines {
		style := identifyStyle(line)
		if style != expected && style != diff3 && style != text {
			return nil, errors.New("Invalid number of markers")
		}

		switch style {
		case text:
			continue
		case start:
			conf.Start = idx
			expected = separator
		case separator:
			conf.Middle = idx
			expected = end
		case diff3:
			conf.Diff3 = []int{idx}
			expected = separator
		case end:
			conf.End = idx
			expected = start
		}

		if style == end {
			if !(conf.Valid()) {
				return nil, errors.New("Invalid number of remaining conflict markers")
			}
			conflicts = append(conflicts, conf)
			conf = Conflict{Start: -1}
		}
	}

	return conflicts, nil
}

func ExtractConflictsIn(f File) (conflicts []Conflict, err error) {
	conflicts, err = GroupConflictMarkers(f.Lines)
	if err != nil {
		return
	}

	for i := range conflicts {
		c := &conflicts[i]
		c.File = &f

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
	topPath, err := topLevelPath(cwd)
	if err != nil {
		return
	}

	targetFiles, err := conflictedFiles(topPath)
	if err != nil {
		return
	}

	for _, fname := range targetFiles {
		absPath := path.Join(topPath, fname)
		file := File{Name: fname, AbsolutePath: absPath}

		if err = file.Read(); err != nil {
			return
		}

		conflicts, err := ExtractConflictsIn(file)
		if err != nil {
			return nil, err
		}

		file.Conflicts = conflicts
		files = append(files, file)
	}
	return
}
