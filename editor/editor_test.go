package editor

import (
	"strings"
	"testing"

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
