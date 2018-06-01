package conflict

import (
	"strings"
	"testing"

	"github.com/mkchoi212/fac/testhelper"
)

func TestRead(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}
		err := f.Read()

		testhelper.Ok(t, err)
		testhelper.Equals(t, len(f.Lines), test.numLines)
	}
}

func TestWriteChanges(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}

		// Read file
		err := f.Read()
		testhelper.Ok(t, err)

		if !test.shouldPass {
			continue
		}

		// Exract conflicts and resolve them
		conflicts, err := ExtractConflicts(f)
		testhelper.Ok(t, err)
		for i := range test.resolveDecision {
			conflicts[i].Choice = test.resolveDecision[i]
		}
		f.Conflicts = conflicts

		// Write changes to file
		f.AbsolutePath = "testdata/.test_output"
		err = f.WriteChanges()
		testhelper.Ok(t, err)

		// Read changes and compare
		f.AbsolutePath = "testdata/.test_output"
		f.Lines = nil
		err = f.Read()
		testhelper.Ok(t, err)

		for i := range f.Lines {
			f.Lines[i] = strings.TrimSuffix(f.Lines[i], "\n")
		}
		testhelper.Equals(t, f.Lines, test.resolved)
	}
}
