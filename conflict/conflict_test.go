package conflict

import (
	"reflect"
	"testing"

	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/testhelper"
)

var dummyFile = struct {
	lines  []string
	start  int
	middle int
	end    int
	diff3  []int
}{
	lines: []string{
		"<<<<<<< Updated upstream:assets/README.md",
		"brew tap mkchoi212/fac https://github.com/mkchoi212/fac.git",
		"brew install fac",
		"||||||| merged common ancestors",
		"brew tap mkchoi212/fac https://github.com/parliament/fac.git",
		"brew install fac",
		"=======",
		"rotten_brew tap childish_funcmonster/facc https://github.com/parliament/facc.git",
		"rotten_brew install facc",
		">>>>>>> Stashed changes:README.md",
		"hello",
		"world",
	},
	start:  0,
	diff3:  []int{3},
	middle: 5,
	end:    9,
}

func TestValidConflict(t *testing.T) {
	invalid := Conflict{
		Start: 2,
		End:   8,
	}
	testhelper.Assert(t, !invalid.Valid(), "%s should be invalid", invalid)

	valid := Conflict{
		Start:  2,
		Middle: 4,
		End:    8,
	}
	testhelper.Assert(t, valid.Valid(), "%s should be valid", invalid)
}

func TestExtract(t *testing.T) {
	c := Conflict{
		Start:  dummyFile.start,
		Middle: dummyFile.middle,
		End:    dummyFile.end,
		Diff3:  dummyFile.diff3,
	}

	err := c.Extract(dummyFile.lines)
	testhelper.Ok(t, err)

	expectedCurrentName := "Updated"
	testhelper.Equals(t, expectedCurrentName, c.CurrentName)

	expectedForeignName := "Stashed"
	testhelper.Equals(t, expectedForeignName, c.ForeignName)

	expectedIncoming := dummyFile.lines[dummyFile.middle+1 : dummyFile.end]
	testhelper.Equals(t, expectedIncoming, c.IncomingLines)

	expectedPureIncoming := dummyFile.lines[dummyFile.start+1 : dummyFile.diff3[0]]
	testhelper.Equals(t, expectedPureIncoming, c.LocalPureLines)
}

func TestUpdate(t *testing.T) {
	for _, test := range tests {
		// Read "manually written" lines
		f := File{AbsolutePath: test.path}
		err := f.Read()
		testhelper.Ok(t, err)

		// Ignore files with more than one conflicts
		// to simulate manual edit
		if test.numConflicts > 1 {
			continue
		}

		// Update empty `Conflict`
		c := Conflict{File: &f}
		err = c.Update(f.Lines)

		if test.shouldPass {
			testhelper.Ok(t, err)
		} else {
			testhelper.Assert(t, err != nil, "%s should not have passed", f)
		}
	}
}

func TestEqual(t *testing.T) {
	foo := File{AbsolutePath: "/bin/foo"}
	bar := File{AbsolutePath: "/bin/bar"}

	c1, c2 := Conflict{Start: 45, File: &foo}, Conflict{Start: 45, File: &foo}
	c3 := Conflict{Start: 45, File: &bar}

	testhelper.Assert(t, c1.Equal(&c2), "%s and %s should be equal", c1, c2)
	testhelper.Assert(t, !(c1.Equal(&c3)), "%s and %s should not be equal", c1, c2)
}

func TestSetContextLines(t *testing.T) {
	var testCases = []struct {
		in       string
		expected int
		err      bool
	}{
		{"2", 2, false},
		{"0", 0, false},
		{"-2", 0, true},
		{"foo", 0, true},
	}

	c := Conflict{}

	for _, tt := range testCases {
		if err := c.SetContextLines(tt.in); err != nil {
			testhelper.Assert(t, tt.err, "Didn't expect error to be returned, input: %s", tt.in)
		}

		testhelper.Assert(t, c.TopPeek == tt.expected, "TopPeek, expected: %s, returned: %s", tt.expected, c.TopPeek)
		testhelper.Assert(t, c.BottomPeek == tt.expected, "BottomPeek, expected:  %s, returned: %s", tt.expected, c.BottomPeek)
	}
}

func TestPaddingLines(t *testing.T) {
	f := File{Lines: dummyFile.lines}
	c := Conflict{
		Start: dummyFile.start,
		End:   dummyFile.end,
		File:  &f,
	}

	top, bottom := c.PaddingLines()
	testhelper.Assert(t, len(top) == 0 && len(bottom) == 0, "top and bottom peak should initially be 0")

	c.TopPeek--
	c.BottomPeek++
	top, bottom = c.PaddingLines()
	expectedBottom := color.Black(color.Regular, f.Lines[dummyFile.end])

	testhelper.Equals(t, len(top), 0)
	testhelper.Equals(t, len(bottom), 1)
	testhelper.Equals(t, bottom[0], expectedBottom)
}

func TestHighlight(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path, Name: test.path}
		err := f.Read()
		testhelper.Ok(t, err)

		conflicts, err := ExtractConflicts(f)
		if test.shouldPass {
			testhelper.Ok(t, err)
		} else {
			testhelper.Assert(t, err != nil, "%s is not parsable", f)
		}

		for _, c := range conflicts {
			_ = c.HighlightSyntax()

			if test.highlightable {
				testhelper.Assert(t, !reflect.DeepEqual(c.IncomingLines, c.ColoredIncomingLines), "%s should be highlighted", f)
			} else {
				testhelper.Equals(t, c.IncomingLines, c.ColoredIncomingLines)
			}
		}
	}
}
