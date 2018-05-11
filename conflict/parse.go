package conflict

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var FileLines map[string][]string

func (c *Conflict) HighlightSyntax() error {
	var lexer chroma.Lexer

	if lexer = lexers.Match(c.FileName); lexer == nil {
		for _, block := range [][]string{c.LocalLines, c.IncomingLines} {
			fmt.Print(strings.Join(block, "\n"))
			if trial := lexers.Analyse(strings.Join(block, "")); trial != nil {
				lexer = trial
				break
			}
		}
	}

	if lexer == nil {
		lexer = lexers.Fallback
		c.ColoredLocalLines = c.LocalLines
		c.ColoredIncomingLines = c.IncomingLines
		return nil
	}

	style := styles.Get("emacs")
	formatter := formatters.Get("terminal")

	var it chroma.Iterator
	var err error
	buf := new(bytes.Buffer)

tokenizer:
	for i, block := range [][]string{c.LocalLines, c.IncomingLines} {
		for _, line := range block {
			if it, err = lexer.Tokenise(nil, line); err != nil {
				break tokenizer
			}
			if err = formatter.Format(buf, style, it); err != nil {
				break tokenizer
			}

			if i == 0 {
				c.ColoredLocalLines = append(c.ColoredLocalLines, buf.String())
			} else {
				c.ColoredIncomingLines = append(c.ColoredIncomingLines, buf.String())
			}
			buf.Reset()
		}
	}
	return err
}

// ReadFile reads all lines of a given file
func ReadFile(absPath string) ([]string, error) {
	input, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	r := bufio.NewReader(input)
	lines := []string{}

	for {
		data, err := r.ReadBytes('\n')
		if err == nil || err == io.EOF {
			// gocui currently doesn't support printing \r
			line := strings.Replace(string(data), "\r", "", -1)
			lines = append(lines, line)
		}

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}
	return lines, nil
}

func parseGitMarkerInfo(diff string, dict map[string][]int) error {
	parts := strings.Split(diff, ":")

	if len(parts) < 3 || !strings.Contains(diff, "marker") {
		return nil
	}

	fname, lineData := string(parts[0]), parts[1]

	if lineNum, err := strconv.Atoi(string(lineData)); err == nil {
		lines := append(dict[fname], lineNum)
		dict[fname] = lines
	}
	return nil
}

func newConflicts(absPath string, fname string, markerLocations []int) ([]Conflict, error) {
	parsedConflicts := []Conflict{}
	lines := FileLines[absPath]

	var conf Conflict

	for _, lineNum := range markerLocations {
		line := lines[lineNum-1]

		switch IdentifyStyle(line) {
		case start:
			conf = Conflict{}
			conf.AbsolutePath = absPath
			conf.FileName = fname
			conf.Start = lineNum
		case separator:
			conf.Middle = lineNum
		case diff3:
			conf.Diff3 = append(conf.Diff3, lineNum)
		case end:
			conf.End = lineNum
			parsedConflicts = append(parsedConflicts, conf)
		default:
			continue
		}
	}

	// Verify all markers are properly parsed/paired
	for _, c := range parsedConflicts {
		if !(c.Valid()) {
			return nil, errors.New("Invalid number of remaining conflict markers")
		}
	}
	return parsedConflicts, nil
}

// Find runs `git --no-pager diff --check` in order to detect git conflicts
// If there are no conflicts, it returns a `ErrNoConflict`
// If there are conflicts, it parses the corresponding files
func Find(cwd string) ([]Conflict, error) {
	allConflicts := []Conflict{}

	topPath, ok := TopLevelPath(cwd)
	if ok != nil {
		return nil, ok
	}

	markerLocations, ok := MarkerLocations(topPath)
	if ok != nil {
		return nil, ok
	}

	markerLocMap := make(map[string][]int)
	FileLines = make(map[string][]string)

	for _, line := range markerLocations {
		if len(line) == 0 {
			continue
		}

		if err := parseGitMarkerInfo(line, markerLocMap); err != nil {
			return nil, err
		}
	}

	for fname := range markerLocMap {
		absPath := path.Join(topPath, fname)

		lines, err := ReadFile(absPath)
		if err != nil {
			return nil, err
		}
		FileLines[absPath] = append(FileLines[absPath], lines...)

		if conflicts, err := newConflicts(absPath, fname, markerLocMap[fname]); err == nil {
			allConflicts = append(allConflicts, conflicts...)
		} else {
			return nil, err
		}
	}

	for i := range allConflicts {
		fileLines := FileLines[allConflicts[i].AbsolutePath]
		if err := allConflicts[i].Extract(fileLines); err != nil {
			return nil, err
		}
		if err := allConflicts[i].HighlightSyntax(); err != nil {
			return nil, err
		}
	}

	return allConflicts, nil
}
