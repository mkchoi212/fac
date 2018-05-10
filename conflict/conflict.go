package conflict

import (
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

	LocalLines           []string
	IncomingLines        []string
	ColoredLocalLines    []string
	ColoredIncomingLines []string

	CurrentName string
	ForeignName string

	TopPeek     int
	BottomPeek  int
	DisplayDiff bool
}

// ErrNoConflict is used to indicate that there
// are no errors present in the git repo
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
	return c.AbsolutePath == c2.AbsolutePath && c.Start == c2.Start
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

// In finds `Conflict`s that are from the provided file name
func In(conflicts []Conflict, fname string) (res []Conflict) {
	for _, c := range conflicts {
		if c.AbsolutePath == fname && c.Choice != 0 {
			res = append(res, c)
		}
	}
	return
}
