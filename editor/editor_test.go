package editor

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/mkchoi212/fac/conflict"
	"github.com/mkchoi212/fac/testhelper"
)

func TestEditorCmd(t *testing.T) {
	editor := editorCmd(".")
	testhelper.Assert(t, editor != nil, "editor should not be nil")
}

func TestWriteTmpFile(t *testing.T) {
	r := strings.NewReader("Hello, Reader!")
	name, err := writeTmpFile(r)

	testhelper.Ok(t, err)
	testhelper.Assert(t, name != "", "tmp file name should not be empty")
}

func TestOpen(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	f := conflict.File{AbsolutePath: "testdata/lorem_ipsum"}
	err := f.Read()
	testhelper.Ok(t, err)

	c := conflict.Conflict{File: &f, Start: 0, End: 5}
	output, err := Open(&c)

	testhelper.Ok(t, err)
	testhelper.Equals(t, f.Lines, output)
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	os.Exit(0)
}

// Allows us to mock exec.Command, thanks to
// https://npf.io/2015/06/testing-exec-command/
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(cmd.Env, "GO_WANT_HELPER_PROCESS=1")
	return cmd
}
