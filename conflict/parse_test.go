package conflict

import (
	"os/exec"
	"testing"

	"github.com/mkchoi212/fac/testhelper"
)

func TestIdentifyStyle(t *testing.T) {
	styles := []int{}

	for _, line := range dummyFile.lines {
		style := identifyStyle(line)
		if style != text {
			styles = append(styles, style)
		}
	}

	expected := []int{start, diff3, separator, end}
	testhelper.Equals(t, expected, styles)
}

func TestParseConflictMarkers(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}
		err := f.Read()
		testhelper.Ok(t, err)

		conflicts, err := GroupConflictMarkers(f.Lines)
		if test.shouldPass {
			testhelper.Ok(t, err)
		} else {
			testhelper.Assert(t, err != nil, "%s should fail", f.AbsolutePath)
		}
		testhelper.Equals(t, len(conflicts), test.numConflicts)
	}
}

func TestParseConflictsIn(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}
		if err := f.Read(); err != nil && test.highlightable {
			t.Error("ParseConflicts/Read failed")
		}

		_, err := ExtractConflicts(f)
		if test.shouldPass {
			testhelper.Assert(t, err == nil, "function should have succeeded")
		} else {
			testhelper.Assert(t, err != nil, "function should have failed")
		}
	}
}

func TestFind(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	files, err := Find("testdata")
	testhelper.Ok(t, err)
	testhelper.Equals(t, len(files), 2)

	for _, f := range files {
		for _, test := range tests {
			if f.AbsolutePath == test.path && test.shouldPass {
				testhelper.Equals(t, len(f.Conflicts), test.numConflicts)
			}
		}
	}
}
