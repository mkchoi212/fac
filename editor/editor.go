package editor

// The following package has been inspired by
// https://github.com/kioopi/extedit

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/mkchoi212/fac/conflict"
)

const defaultEditor = "vim"

var execCommand = exec.Command

// Open starts a text-editor with lines from `Content`
// It returns the manually edited lines from the text-editor when the user closes the editor
func Open(c *conflict.Conflict) (output []string, err error) {
	lines := c.File.Lines[c.Start : c.End+1]

	content := strings.NewReader(strings.Join(lines, ""))
	input, err := contentFromReader(content)

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

	newContent, err := contentFromFile(fileName)
	if err != nil {
		return
	}

	output = newContent.c
	return
}

// writeTmpFile writes content to a temporary file and returns the path to the file
// It returns an error if the temporary file cannot be created
func writeTmpFile(content io.Reader) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	io.Copy(f, content)
	f.Close()
	return f.Name(), nil
}

// editorCmd returns a os/exec.Cmd to open the provided file
func editorCmd(filename string) *exec.Cmd {
	editorEnv := os.Getenv("EDITOR")
	if editorEnv == "" {
		editorEnv = defaultEditor
	}

	editorVars := strings.Split(editorEnv, " ")

	path := editorVars[0]
	args := []string{filename}

	if len(editorVars) > 1 {
		args = append(editorVars[1:], args...)
	}

	editor := execCommand(path, args...)

	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	return editor
}
