package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
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

	relMid, relEnd := c.Middle-c.Start, c.End-c.Start
	c.CurrentLines = allLines[1:relMid]
	c.ForeignLines = allLines[relMid+1 : relEnd]
	c.CurrentName = strings.Split(allLines[0], " ")[1]
	c.ForeignName = strings.Split(allLines[c.End-c.Start], " ")[1]
	return nil
}

func parse(diff []byte, dict map[string][]int) error {
	parts := bytes.Split(diff, []byte(":"))

	if len(parts) < 3 {
		return errors.New("Could not parse line")
	}

	fname, lineData := string(parts[0]), parts[1]

	if lineNum, err := strconv.Atoi(string(lineData)); err == nil {
		lines := append(dict[fname], lineNum)
		dict[fname] = lines
	}
	return nil
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
		_ = parse(line, diffDict)
	}

	conflicts := []Conflict{}

	for fname := range diffDict {
		conf := Conflict{}
		sort.Ints(diffDict[fname])
		lines := diffDict[fname]
		conf.Start, conf.Middle, conf.End = lines[0], lines[1], lines[2]
		conf.FileName = fname
		conf.AbsolutePath = path.Join(cmd.Dir, fname)
		conf.Resolved = false

		conflicts = append(conflicts, conf)
	}

	for i := range conflicts {
		if err := conflicts[i].ExtractLines(); err != nil {
			log.Panicln(err)
		}
	}

	return conflicts, nil
}
