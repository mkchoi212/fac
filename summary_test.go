package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
)

var dummyFile = conflict.File{Name: "Foobar"}
var dummyConflicts = []*conflict.Conflict{
	{Choice: conflict.Local, File: &dummyFile},
	{Choice: conflict.Local, File: &dummyFile},
	{Choice: 0, File: &dummyFile},
}

var expected_two_resolved = []string{
	"[32;1m✔ Foobar: 0[0m",
	"[32;1m✔ Foobar: 0[0m",
	"[31;1m✘ Foobar: 0[0m",
	"",
	"Resolved [31;2m2 [0mconflict(s) out of [31;2m3[0m",
}

var expected_all_resolved = []string{
	"[32;1m✔ Foobar: 0[0m",
	"[32;1m✔ Foobar: 0[0m",
	"[32;1m✔ Foobar: 0[0m",
	"[32;1m",
	"Fixed All Conflicts 🎉[0m",
	"",
}

func TestPrintHelp(t *testing.T) {
	var b bytes.Buffer
	printHelp(&b)

	out := b.String()
	if out != color.Blue(color.Regular, instruction) {
		t.Errorf("PrintHelp failed: wanted blue colored %s..., got %s", instruction[:45], out)
	}
}

func TestSummary(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printSummary(dummyConflicts)
	dummyConflicts[2].Choice = conflict.Local
	printSummary(dummyConflicts)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout
	output := strings.Split(string(out), "\n")
	expected := append(expected_two_resolved, expected_all_resolved...)

	if len(expected) != len(output) {
		t.Errorf("Summary failed: got \n%v, wanted \n%v", output, expected)
	}

	for i := range expected {
		expectedLine := []byte(expected[i])
		outputLine := []byte(output[i])
		// Remove ESC bytes
		outputLine = bytes.Trim(outputLine, string([]byte{27}))

		if len(expectedLine) == 0 {
			continue
		}

		if !(reflect.DeepEqual(expectedLine, outputLine)) {
			t.Errorf("Summary failed: got %s, wanted %s", outputLine, expectedLine)
		}
	}
}
