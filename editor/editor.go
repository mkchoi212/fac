package editor

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
	result    Content
	SplitFunc bufio.SplitFunc
}

// Open starts a text-editor with the contents of content.
// It returns edited content after user closes the editor
func Open(conf *conflict.Conflict) (output []string, err error) {
	s := NewSession()

	lines := append([]string{}, conf.LocalLines...)
	lines = append(lines, "=======\n")
	lines = append(lines, conf.IncomingLines...)

	content := strings.NewReader(strings.Join(lines, ""))
	input, err := contentFromReader(content, s.SplitFunc)

	if err != nil {
		return
	}

	fileName, err := writeTmpFile(input)
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

	return
}

func NewSession() *Session {
	return &Session{SplitFunc: bufio.ScanLines}
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
