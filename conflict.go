package main

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
)

// Conflict represents a single conflict that may have occured
type Conflict struct {
	Resolved     bool
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

	topPeek    int
	bottomPeek int
}

const (
	Local    = 0
	Incoming = 1
)

func (c *Conflict) isEqual(c2 *Conflict) bool {
	return c.FileName == c2.FileName && c.Start == c2.Start
}

func (c *Conflict) Select(g *gocui.Gui, withHelp bool) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Panel)
		if err != nil {
			return err
		}
		v.Clear()

		for idx, conflict := range conflicts {
			var out string
			if conflict.Resolved {
				out = Colorize(fmt.Sprintf("✅  %s:%d", conflict.FileName, conflict.Start), Green)
			} else {
				out = Colorize(fmt.Sprintf("%d.  %s:%d", idx+1, conflict.FileName, conflict.Start), Red)
			}

			if conflict.isEqual(c) {
				fmt.Fprintf(v, "%s <-\n", out)
			} else {
				fmt.Fprintf(v, "%s\n", out)
			}
		}

		if withHelp {
			printHelp(v)
		}
		return nil
	})

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Current)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.WriteString(c.CurrentName)
		buf.WriteString(" (Current Change) ")
		v.Title = buf.String()

		printLines(v, c.ColoredCurrentLines)

		v, err = g.View(Foreign)
		if err != nil {
			return err
		}
		buf.Reset()
		buf.WriteString(c.ForeignName)
		buf.WriteString(" (Incoming Change) ")
		v.Title = buf.String()

		printLines(v, c.ColoredForeignLines)
		return nil
	})
	return nil
}

func (c *Conflict) Resolve(g *gocui.Gui, v *gocui.View, version int) error {
	g.Update(func(g *gocui.Gui) error {
		c.Resolved = true
		nextConflict(g, v)
		return nil
	})
	return nil
}

func nextConflict(g *gocui.Gui, v *gocui.View) error {
	originalCur := cur

	for originalCur != cur {
		cur++
		if cur >= conflictCount {
			cur = 0
		}
	}

	if originalCur == cur {
		globalQuit(g)
	}

	conflicts[cur].Select(g, false)
	return nil
}
