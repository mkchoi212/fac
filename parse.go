package main

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

var allFileLines map[string][]string

func (c *Conflict) Highlight() error {
	var lexer chroma.Lexer

	if lexer = lexers.Match(c.FileName); lexer == nil {
		for _, block := range [][]string{c.CurrentLines, c.ForeignLines} {
			fmt.Print(strings.Join(block, "\n"))
			if trial := lexers.Analyse(strings.Join(block, "")); trial != nil {
				lexer = trial
				break
			}
		}
	}

	if lexer == nil {
		lexer = lexers.Fallback
		c.ColoredCurrentLines = c.CurrentLines
		c.ColoredForeignLines = c.ForeignLines
		return nil
	}

	style := styles.Get("emacs")
	formatter := formatters.Get("terminal")

	var it chroma.Iterator
	var err error
	buf := new(bytes.Buffer)

tokenizer:
	for i, block := range [][]string{c.CurrentLines, c.ForeignLines} {
		for _, line := range block {
			if it, err = lexer.Tokenise(nil, line); err != nil {
				break tokenizer
			}
			if err = formatter.Format(buf, style, it); err != nil {
				break tokenizer
			}

			if i == 0 {
				c.ColoredCurrentLines = append(c.ColoredCurrentLines, buf.String())
			} else {
				c.ColoredForeignLines = append(c.ColoredForeignLines, buf.String())
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
			allFileLines[absPath] = append(allFileLines[absPath], string(data))
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
	lines, ok := allFileLines[c.AbsolutePath]
	if !ok {
		if err := ReadFile(c.AbsolutePath); err != nil {
			log.Panic(err)
		}
	}

	lines, _ = allFileLines[c.AbsolutePath]
	c.CurrentLines = lines[c.Start : c.Middle-1]
	c.ForeignLines = lines[c.Middle : c.End-1]
	c.CurrentName = strings.Split(lines[c.Start-1], " ")[1]
	c.ForeignName = strings.Split(lines[c.End-1], " ")[1]
	return nil
}

func parseRawOutput(diff string, dict map[string][]int) error {
	parts := strings.Split(diff, ":")

	if len(parts) < 3 || !strings.Contains(diff, "marker") {
		return errors.New("Could not parse line")
	}

	fname, lineData := string(parts[0]), parts[1]

	if lineNum, err := strconv.Atoi(string(lineData)); err == nil {
		lines := append(dict[fname], lineNum)
		dict[fname] = lines
	}
	return nil
}

func groupConflictOutput(fname string, cwd string, lines []int) ([]Conflict, error) {
	if len(lines)%3 != 0 {
		return nil, errors.New("Invalid number of remaining conflict markers")
	}

	conflicts := []Conflict{}
	for i := 0; i < len(lines); i += 3 {
		conf := Conflict{}
		conf.Start = lines[i]
		conf.Middle = lines[i+1]
		conf.End = lines[i+2]
		conf.FileName = fname
		conf.AbsolutePath = path.Join(cwd, fname)
		conflicts = append(conflicts, conf)
	}

	return conflicts, nil
}

func Find() ([]Conflict, error) {
	dummyPath := "/Users/mikechoi/src/CSCE-313/"

	cmdName := "git"
	cmdArgs := []string{"--no-pager", "diff", "--check"}
	var (
		cmdOut []byte
		err    error
	)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = dummyPath

	cmdOut, err = cmd.Output()
	if !(strings.Contains(err.Error(), "exit status 2")) {
		return nil, NewErrNoConflict("No conflicts detected 🎉")
	}

	output := bytes.Split(cmdOut, []byte("\n"))
	diffDict := make(map[string][]int)

	for _, line := range output {
		_ = parseRawOutput(string(line), diffDict)
	}

	conflicts := []Conflict{}

	for fname := range diffDict {
		if parsedConflicts, err := groupConflictOutput(fname, cmd.Dir, diffDict[fname]); err == nil {
			for _, c := range parsedConflicts {
				conflicts = append(conflicts, c)
			}
		} else {
			return nil, err
		}
	}

	allFileLines = make(map[string][]string)
	for i := range conflicts {
		if err := conflicts[i].ExtractLines(); err != nil {
			return nil, err
		}
		if err := conflicts[i].Highlight(); err != nil {
			return nil, err
		}
	}

	return conflicts, nil
}
