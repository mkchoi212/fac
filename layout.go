package main

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
)

const (
	Current = "current"
	Foreign = "foreign"
	Panel   = "panel"
	Prompt  = "prompt"
	Input   = "input prompt"

	Local    = 1
	Incoming = 2
	Up       = 3
	Down     = 4
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	inputHeight := 2
	viewHeight := maxY - inputHeight
	branchViewWidth := (maxX / 5) * 2

	if _, err := g.SetView(Current, 0, 0, branchViewWidth, viewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView(Foreign, branchViewWidth, 0, branchViewWidth*2, viewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if v, err := g.SetView(Panel, branchViewWidth*2, 0, maxX-2, viewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Conflicts"
	}

	if v, err := g.SetView(Prompt, 0, viewHeight, 14, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		prompt := color.Green(color.Regular, "[wasd] >>")
		v.Write([]byte(prompt))
		v.MoveCursor(11, 0, true)
	}

	if v, err := g.SetView(Input, 10, viewHeight, maxX, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Editable = true
		v.Wrap = false
		v.Editor = gocui.EditorFunc(promptEditor)
		if _, err := g.SetCurrentView(Input); err != nil {
			return err
		}
	}
	return nil
}

func Select(c *conflict.Conflict, g *gocui.Gui, showHelp bool) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Panel)
		if err != nil {
			return err
		}
		v.Clear()

		for idx, conflict := range conflict.All {
			var out string
			if conflict.Choice != 0 {
				out = color.Green(color.Regular, "✅  %s:%d", conflict.FileName, conflict.Start)
			} else {
				out = color.Red(color.Regular, "%d. %s:%d", idx+1, conflict.FileName, conflict.Start)
			}

			if conflict.Equal(c) {
				fmt.Fprintf(v, "%s <-\n", out)
			} else {
				fmt.Fprintf(v, "%s\n", out)
			}
		}

		if showHelp {
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

		top, bottom := c.PaddingLines()
		v.Clear()
		printLines(v, top)
		if c.DisplayDiff {
			printLines(v, c.Diff())
		} else {
			printLines(v, c.ColoredCurrentLines)
		}
		printLines(v, bottom)

		v, err = g.View(Foreign)
		if err != nil {
			return err
		}
		buf.Reset()
		buf.WriteString(c.ForeignName)
		buf.WriteString(" (Incoming Change) ")
		v.Title = buf.String()

		top, bottom = c.PaddingLines()
		v.Clear()
		printLines(v, top)
		printLines(v, c.ColoredForeignLines)
		printLines(v, bottom)
		return nil
	})
	return nil
}

func Resolve(c *conflict.Conflict, g *gocui.Gui, v *gocui.View, version int) error {
	g.Update(func(g *gocui.Gui) error {
		c.Choice = version
		MoveToItem(Down, g, v)
		return nil
	})
	return nil
}

func MoveToItem(dir int, g *gocui.Gui, v *gocui.View) error {
	originalCur := cur

	for {
		if dir == Up {
			cur--
		} else {
			cur++
		}

		if cur >= conflict.Count {
			cur = 0
		} else if cur < 0 {
			cur = conflict.Count - 1
		}

		if conflict.All[cur].Choice == 0 || originalCur == cur {
			break
		}
	}

	if originalCur == cur {
		globalQuit(g)
	}

	Select(&conflict.All[cur], g, false)
	return nil
}

func Scroll(g *gocui.Gui, c *conflict.Conflict, direction int) {
	if direction == Up {
		c.TopPeek--
		c.BottomPeek++
	} else if direction == Down {
		c.TopPeek++
	} else {
		return
	}

	Select(c, g, false)
}
