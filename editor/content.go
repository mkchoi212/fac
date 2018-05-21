package editor

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Content represents the lines read from a io.Reader
// `c` holds the actual lines and
// `reader` contains the io.Reader representation of the lines
type Content struct {
	c      []string
	reader io.Reader
}

func (c Content) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

func (c Content) String() string {
	return strings.Join(c.c, "")
}

// contentFromReader creates a new `Content` object
// filled with the lines from the provided Reader
// It returns an error if anything other than io.EOF is raised
func contentFromReader(content io.Reader) (c Content, err error) {
	reader := bufio.NewReader(content)

	for {
		line, err := reader.ReadString('\n')
		c.c = append(c.c, line)

		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return
	}

	c.reader = strings.NewReader(c.String())
	return
}

// contentFromFile reads the content from the file
// It returns an error if the file does not exist
func contentFromFile(filename string) (Content, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Content{}, err
	}
	defer file.Close()

	return contentFromReader(file)
}
