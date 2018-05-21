package editor

// The following package has been inspired by
// https://github.com/kioopi/extedit

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/mkchoi212/fac/conflict"
)

const defaultEditor = "vim"

type Session struct {
	input     Content
	SplitFunc bufio.SplitFunc
}

// Open starts a text-editor with the contents of content.
// It returns edited content after user closes the editor
func Open(c *conflict.Conflict) (output []string, err error) {
	s := &Session{SplitFunc: bufio.ScanLines}
	lines := c.File.Lines[c.Start : c.End+1]

	content := strings.NewReader(strings.Join(lines, ""))
	input, err := contentFromReader(content, s.SplitFunc)
	if err != nil {
		return
	}

	fileName, err := writeTmpFile(input.reader)
	if err != nil {
		return
	}

	cmd := editorCmd(fileName)
	err = cmd.Run()
	if err != nil {
		return
	}

	newContent, err := contentFromFile(fileName, s.SplitFunc)
	if err != nil {
		return
	}

	output = newContent.c
	for i := range output {
		output[i] = output[i] + "\n"
	}
	return
}

// writeTmpFile writes content to a temporary file and returns
// the path to the file
func writeTmpFile(content io.Reader) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	io.Copy(f, content)
	f.Close()
	return f.Name(), nil
}

// editorCmd creates a os/exec.Cmd to open
// filename in an editor ready to be run()
func editorCmd(filename string) *exec.Cmd {
	editorPath := os.Getenv("EDITOR")
	if editorPath == "" {
		editorPath = defaultEditor
	}
	editor := exec.Command(editorPath, filename)

	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	return editor
}
