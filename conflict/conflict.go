package conflict

import (
	"sort"
	"strings"

	"github.com/mkchoi212/fac/color"
)

// Conflict represents a single conflict that may have occurred
type Conflict struct {
	Choice       int
	FileName     string
	AbsolutePath string
	Start        int
	Middle       int
	End          int
	Diff3        []int

	LocalLines           []string
	LocalPureLines       []string
	IncomingLines        []string
	ColoredLocalLines    []string
	ColoredIncomingLines []string

	CurrentName string
	ForeignName string

	TopPeek     int
	BottomPeek  int
	DisplayDiff bool
}

// Supported git conflict styles
const (
	text = iota
	start
	diff3
	separator
	end
)

// IdentifyStyle identifies the conflict marker style of provided text
func IdentifyStyle(line string) (style int) {
	line = strings.TrimSpace(line)

	if strings.Contains(line, "<<<<<<<") {
		style = start
	} else if strings.Contains(line, ">>>>>>>") {
		style = end
	} else if line == "||||||| merged common ancestors" {
		style = diff3
	} else if line == "=======" {
		style = separator
	} else {
		style = text
	}
	return
}

// Valid checks if the parsed conflict has corresponding begin, separator,
// and middle line numbers
func (c *Conflict) Valid() bool {
	return c.Middle != 0 && c.End != 0
}

func (c *Conflict) Equal(c2 *Conflict) bool {
	return c.AbsolutePath == c2.AbsolutePath && c.Start == c2.Start
}

func (c *Conflict) ToggleDiff() {
	c.DisplayDiff = !(c.DisplayDiff)
}

// Extract extracts lines where conflicts exist
// and corresponding branch names
func (c *Conflict) Extract(lines []string) error {
	c.LocalLines = lines[c.Start : c.Middle-1]
	if len(c.Diff3) != 0 {
		sort.Ints(c.Diff3)
		diff3Barrier := c.Diff3[0]
		c.LocalPureLines = lines[c.Start : diff3Barrier-1]
	} else {
		c.LocalPureLines = c.LocalLines
	}
	c.IncomingLines = lines[c.Middle : c.End-1]
	c.CurrentName = strings.Split(lines[c.Start-1], " ")[1]
	c.ForeignName = strings.Split(lines[c.End-1], " ")[1]
	return nil
}

func (c *Conflict) PaddingLines() (topPadding, bottomPadding []string) {
	lines := FileLines[c.AbsolutePath]
	start, end := c.Start-1, c.End

	if c.TopPeek >= start {
		c.TopPeek = start
	} else if c.TopPeek < 0 {
		c.TopPeek = 0
	}

	for _, l := range lines[start-c.TopPeek : start] {
		topPadding = append(topPadding, color.Black(color.Regular, l))
	}

	if c.BottomPeek >= len(lines)-c.End {
		c.BottomPeek = len(lines) - c.End
	} else if c.BottomPeek < 0 {
		c.BottomPeek = 0
	}

	for _, l := range lines[end : end+c.BottomPeek] {
		bottomPadding = append(bottomPadding, color.Black(color.Regular, l))
	}
	return
}

// In finds `Conflict`s that are from the provided file name
func In(fname string, conflicts []Conflict) (res []Conflict) {
	for _, c := range conflicts {
		if c.AbsolutePath == fname && c.Choice != 0 {
			res = append(res, c)
		}
	}
	return
}
