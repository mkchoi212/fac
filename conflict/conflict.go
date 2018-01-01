package conflict

import (
	"strings"

	"github.com/mkchoi212/fac/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Conflict represents a single conflict that may have occured
type Conflict struct {
	Choice       int
	FileName     string
	AbsolutePath string
	Start        int
	Middle       int
	End          int

	CurrentLines        []string
	ForeignLines        []string
	ColoredCurrentLines []string
	ColoredForeignLines []string

	CurrentName string
	ForeignName string

	TopPeek     int
	BottomPeek  int
	DisplayDiff bool
}

var All = []Conflict{}
var Count int

type ErrNoConflict struct {
	message string
}

func NewErrNoConflict(message string) *ErrNoConflict {
	return &ErrNoConflict{
		message: message,
	}
}

func (e *ErrNoConflict) Error() string {
	return e.message
}

func (c *Conflict) Equal(c2 *Conflict) bool {
	return c.FileName == c2.FileName && c.Start == c2.Start
}

func (c *Conflict) ToggleDiff() {
	c.DisplayDiff = !(c.DisplayDiff)
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

func In(fname string) (list []Conflict) {
	for _, c := range All {
		if c.AbsolutePath == fname && c.Choice != 0 {
			list = append(list, c)
		}
	}
	return
}

func (c *Conflict) Diff() []string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(strings.Join(c.CurrentLines, ""), strings.Join(c.ForeignLines, ""), false)
	return []string{dmp.DiffPrettyText(diffs)}
}
