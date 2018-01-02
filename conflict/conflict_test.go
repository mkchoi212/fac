package conflict

import (
	"testing"
)

func TestEqual(t *testing.T) {
	c1, c2, c3 := Conflict{}, Conflict{}, Conflict{}

	c1.FileName, c2.FileName, c3.FileName = "foobar", "foobar", "foobar"
	c1.Start, c2.Start, c3.Start = 45, 45, 45

	c1.AbsolutePath, c2.AbsolutePath = "/path/foobar", "/path/foobar"
	c3.AbsolutePath = "/other/path/foobar"

	if c1.Equal(&c2) != true {
		t.Errorf("%v and %v should be equal", c1, c2)
	}

	if c1.Equal(&c3) != false {
		t.Errorf("%v and %v should NOT be equal", c1, c3)
	}
}

func TestToggleDiff(t *testing.T) {
	c1 := Conflict{}
	c1.ToggleDiff()
	if c1.DisplayDiff != true {
		t.Errorf("%v should be toggled ON", c1)
	}

	c1.ToggleDiff()
	if c1.DisplayDiff != false {
		t.Errorf("%v should be toggled OFF", c1)
	}
}
