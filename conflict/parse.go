package conflict

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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

func ReadFile(absPath string) error {
	input, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	r := bufio.NewReader(input)

	for {
		data, err := r.ReadBytes('\n')
		if err == nil || err == io.EOF {
			FileLines[absPath] = append(FileLines[absPath], string(data))
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}
	return nil
}

func (c *Conflict) ExtractLines() error {
	lines := FileLines[c.AbsolutePath]
	c.LocalLines = lines[c.Start : c.Middle-1]
	c.IncomingLines = lines[c.Middle : c.End-1]
	c.CurrentName = strings.Split(lines[c.Start-1], " ")[1]
	c.ForeignName = strings.Split(lines[c.End-1], " ")[1]
	return nil
}

func parseRawOutput(diff string, dict map[string][]int) error {
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

func New(absPath string, lines []int) ([]Conflict, error) {
	// Check for diff3 output before parsing
	for _, line := range FileLines[absPath] {
		if strings.Contains(line, "||||||| merged common ancestors") {
			return nil, errors.New(`fac does not support diff3 styled outputs yet 😞
Run below command to change to a compatible conflict style

>> git config --global merge.conflictstyle merge`)
		}
	}

	if len(lines)%3 != 0 {
		return nil, errors.New("Invalid number of remaining conflict markers")
	}

	parsedConflicts := []Conflict{}
	for i := 0; i < len(lines); i += 3 {
		conf := Conflict{}
		conf.Start = lines[i]
		conf.Middle = lines[i+1]
		conf.End = lines[i+2]
		conf.AbsolutePath = absPath
		_, conf.FileName = path.Split(absPath)
		parsedConflicts = append(parsedConflicts, conf)
	}

	return parsedConflicts, nil
}

func Find() (err error) {
	cwd, _ := os.Getwd()

	stdout, stderr, _ := RunCommand("git", cwd, "rev-parse", "--show-toplevel")
	if len(stderr) != 0 {
		return errors.New(stderr)
	} else if len(stdout) == 0 {
		return errors.New("no git top-level path")
	}
	topLevelPath := string(strings.Split(stdout, "\n")[0])

	stdout, stderr, _ = RunCommand("git", cwd, "--no-pager", "diff", "--check")

	if len(stderr) != 0 {
		return errors.New(stderr)
	} else if len(stdout) == 0 {
		return NewErrNoConflict("No conflicts detected 🎉")
	}

	stdoutLines := strings.Split(stdout, "\n")
	diffMap := make(map[string][]int)
	FileLines = make(map[string][]string)

	for _, line := range stdoutLines {
		if len(line) == 0 {
			continue
		}

		if err = parseRawOutput(line, diffMap); err != nil {
			return
		}
	}

	for fname := range diffMap {
		absPath := path.Join(topLevelPath, fname)
		if err = ReadFile(absPath); err != nil {
			return
		}
		if conflicts, err := New(absPath, diffMap[fname]); err == nil {
			All = append(All, conflicts...)
		} else {
			return err
		}
	}

	Count = len(All)
	if Count == 0 {
		return NewErrNoConflict("No conflicts detected 🎉")
	}

	for i := range All {
		if err = All[i].ExtractLines(); err != nil {
			return
		}
		if err = All[i].HighlightSyntax(); err != nil {
			return
		}
	}

	return nil
}
