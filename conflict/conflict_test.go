package conflict

import (
	"reflect"
	"testing"
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
	middle: 6,
	end:    9,
}

func TestIdentifyStyle(t *testing.T) {
	styles := []int{}

	for _, line := range dummyFile.lines {
		style := IdentifyStyle(line)
		if style != text {
			styles = append(styles, style)
		}
	}

	expected := []int{start, diff3, separator, end}
	if !(reflect.DeepEqual(styles, expected)) {
		t.Errorf("IdentifyStyle failed: got %v, want %v", styles, expected)
	}
}

func TestValidConflict(t *testing.T) {
	invalid := Conflict{CurrentName: "foobar"}
	if invalid.Valid() {
		t.Errorf("Valid failed: got %t, want %t", true, false)
	}

	valid := Conflict{
		File:        &File{},
		CurrentName: "foo",
		ForeignName: "bar",
		Middle:      4,
		End:         8,
	}
	if !(valid.Valid()) {
		t.Errorf("Valid failed: got %t, want %t", false, true)
	}
}

func TestExtract(t *testing.T) {
	c := Conflict{
		Start:  dummyFile.start,
		Middle: dummyFile.middle,
		End:    dummyFile.end,
		Diff3:  dummyFile.diff3,
	}

	if err := c.Extract(dummyFile.lines); err != nil {
		t.Errorf("Extract failed: could not parse file lines")
	}

	expectedCurrentName := "Updated"
	if expectedCurrentName != c.CurrentName {
		t.Errorf("Extract failed: got %s want %s", c.CurrentName, expectedCurrentName)
	}

	expectedForeignName := "Stashed"
	if expectedForeignName != c.ForeignName {
		t.Errorf("Extract failed: got %s want %s", c.ForeignName, expectedForeignName)
	}

	expectedIncoming := dummyFile.lines[dummyFile.middle+1 : dummyFile.end]
	if !(reflect.DeepEqual(c.IncomingLines, expectedIncoming)) {
		t.Errorf("Extract failed: got %v want %v", c.IncomingLines, expectedIncoming)
	}

	expectedPureIncoming := dummyFile.lines[dummyFile.start+1 : dummyFile.diff3[0]]
	if !(reflect.DeepEqual(c.LocalPureLines, expectedPureIncoming)) {
		t.Errorf("Extract failed: got %v want %v", c.LocalPureLines, expectedPureIncoming)
	}
}

func TestEqual(t *testing.T) {
	foo := File{AbsolutePath: "/bin/foo"}
	bar := File{AbsolutePath: "/bin/bar"}

	c1, c2 := Conflict{Start: 45, File: &foo}, Conflict{Start: 45, File: &foo}
	c3 := Conflict{Start: 45, File: &bar}

	if c1.Equal(&c2) != true {
		t.Errorf("Equal failed: got %t, want %t", false, true)
	}

	if c1.Equal(&c3) != false {
		t.Errorf("Equal failed: got %t, want %t", true, false)
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
	if len(top) != 0 || len(bottom) != 0 {
		t.Errorf("PaddingLines failed: got %d, %d lines, want 0, 0 lines", len(top), len(bottom))
	}

	c.TopPeek--
	c.BottomPeek++
	top, bottom = c.PaddingLines()
	expectedBottom := f.Lines[dummyFile.end+1]
	if len(top) != 0 || reflect.DeepEqual(bottom, expectedBottom) {
		t.Errorf("PaddingLines failed: got %v, want %v", bottom, expectedBottom)
	}

}

func TestHighlight(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path, Name: test.path}

		if err := f.Read(); err != nil && test.highlightable {
			t.Error("conflict/Read failed")
		}

		conflicts, err := parseConflictsIn(f, test.markers)
		if err != nil {
			if !test.parsable {
				continue
			}
			t.Error("conflict/parseConflicts failed")
		}

		for _, c := range conflicts {
			if err := c.HighlightSyntax(); err != nil {
				t.Errorf("HighlightSyntax failed: %s", err.Error())
			}

			if test.highlightable && reflect.DeepEqual(c.IncomingLines, c.ColoredIncomingLines) {
				t.Errorf("HighlightSyntax failed: %s has not been highlighted", f.Name)
			}
		}
	}
}
