package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// Conflict represents a single conflict that may have occured
type Conflict struct {
	Resolved     bool
	FileName     string
	AbsolutePath string
	Start        int
	Middle       int
	End          int

	CurrentLines []string
	ForeignLines []string

	CurrentName string
	ForeignName string
}

// Resolve resolves the conflict
func (c *Conflict) Resolve() {
	c.Resolved = true
}

func (c *Conflict) ExtractLines() error {
	input, err := os.Open(c.AbsolutePath)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	line := 0
	r := bufio.NewReader(input)
	allLines := []string{}

	for {
		data, err := r.ReadBytes('\n')
		line++
		if line > c.End && line < c.Start {
			continue
		}

		if err == nil || err == io.EOF {
			if len(data) > 0 && data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
			if len(data) > 0 && data[len(data)-1] == '\r' {
				data = data[:len(data)-1]
			}
			allLines = append(allLines, string(data))
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}

	c.CurrentLines = allLines[c.Start : c.Middle-1]
	c.ForeignLines = allLines[c.Middle : c.End-1]
	c.CurrentName = strings.Split(allLines[0], " ")[1]
	c.ForeignName = strings.Split(allLines[c.End-1], " ")[1]
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
	fmt.Println(lines)
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
		return nil, err
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
			log.Panic(err)
		}
	}

	fmt.Println(conflicts)

	for i := range conflicts {
		if err := conflicts[i].ExtractLines(); err != nil {
			log.Panicln(err)
		}
	}

	return conflicts, nil
}
