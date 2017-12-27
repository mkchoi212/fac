package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// Conflict represents a single conflict that may have occured
type Conflict struct {
	Resolved  bool
	FileName  string
	StartLine int
}

// Resolve resolves the conflict
func (c *Conflict) Resolve() {
	c.Resolved = true
}

var FakeData = [...]Conflict{Conflict{false, "index.js", 100},
	Conflict{false, "style.css", 32},
	Conflict{false, "hey.js", 69},
}

func parse(diff []byte) (Conflict, error) {
	parts := bytes.Split(diff, []byte(":"))
	conflict := Conflict{}

	if len(parts) != 3 {
		return conflict, errors.New("Could not parse line")
	}

	for i, d := range parts {
		if i == 0 {
			conflict.FileName = string(d)
		} else if i == 1 {
			if lineNum, err := strconv.Atoi(string(d)); err != nil {
				conflict.StartLine = lineNum
			}
		}
	}
	conflict.Resolved = false
	return conflict, nil
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

	diffs := bytes.Split(cmdOut, []byte("\n"))
	conflicts := []Conflict{}

	for _, line := range diffs {
		if conflict, err := parse(line); err == nil {
			conflicts = append(conflicts, conflict)
		}
	}

	return conflicts, nil
}
