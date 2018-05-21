package editor

import (
	"bufio"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mkchoi212/fac/testhelper"
)

func TestUseContentAsReader(t *testing.T) {
	c, err := contentFromReader(strings.NewReader("foo\nbar"), bufio.ScanLines)
	testhelper.Ok(t, err)

	byteCnt, err := ioutil.ReadAll(c)
	testhelper.Ok(t, err)
	testhelper.Equals(t, string(byteCnt), "foo\nbar")
}

func TestCreateContentFromFile(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	testhelper.Ok(t, err)
	_, err = f.Write([]byte("foo\nbar"))
	testhelper.Ok(t, err)

	c, err := contentFromFile(f.Name(), bufio.ScanLines)

	testhelper.Equals(t, c.c[0], "foo")
	testhelper.Equals(t, c.c[1], "bar")
}
