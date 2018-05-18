package conflict

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// File represents a single file that contains git merge conflicts
type File struct {
	AbsolutePath string
	Name         string
	Lines        []string
	Conflicts    []Conflict
}

// readFile reads all lines of a given file
func (f *File) Read() (err error) {
	input, err := os.Open(f.AbsolutePath)
	if err != nil {
		return
	}
	defer input.Close()

	r := bufio.NewReader(input)

	for {
		data, err := r.ReadBytes('\n')
		if err == nil || err == io.EOF {
			// gocui currently doesn't support printing \r
			line := strings.Replace(string(data), "\r", "", -1)
			f.Lines = append(f.Lines, line)
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}

	return
}

// WriteChanges writes all the resolved conflicts in a given file
// to the file system.
func (f File) WriteChanges() (err error) {
	var replacementLines []string

	for _, c := range f.Conflicts {
		if c.Choice == Local {
			replacementLines = append([]string{}, c.LocalPureLines...)
		} else if c.Choice == Incoming {
			replacementLines = append([]string{}, c.IncomingLines...)
		} else {
			continue
		}

		i := 0
		for ; i < len(replacementLines); i++ {
			f.Lines[c.Start+i] = replacementLines[i]
		}
		for ; i <= c.End-c.Start; i++ {
			f.Lines[c.Start+i] = ""
		}
	}

	if err = write(f.AbsolutePath, f.Lines); err != nil {
		return
	}
	return
}

func write(absPath string, lines []string) (err error) {
	f, err := os.Create(absPath)
	if err != nil {
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err = w.WriteString(line); err != nil {
			return
		}
	}
	err = w.Flush()
	return
}
