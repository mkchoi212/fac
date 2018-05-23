package main

import (
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

	Up   = 3
	Down = 4

	Horizontal = 5
	Vertical   = -6
)

var (
	ViewOrientation = Vertical
	inputHeight     = 2
)

func layout(g *gocui.Gui) (err error) {
	if err = makePanels(g); err != nil {
		return
	}

	if err = makeOverviewPanel(g); err != nil {
		return
	}

	if err = makePrompt(g); err != nil {
		return
	}

	return
}

func makePanels(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewHeight := maxY - inputHeight
	branchViewWidth := (maxX / 5) * 2
	isOdd := maxY%2 == 1

	var x0, x1, y0, y1 int
	var x2, x3, y2, y3 int

	if ViewOrientation == Horizontal {
		x0, x1 = 0, branchViewWidth
		y0, y1 = 0, viewHeight
		x2, x3 = branchViewWidth, branchViewWidth*2
		y2, y3 = 0, viewHeight
	} else {
		branchViewWidth = branchViewWidth * 2
		viewHeight = (maxY - inputHeight) / 2

		x0, x1 = 0, branchViewWidth
		y0, y1 = 0, viewHeight
		x2, x3 = 0, branchViewWidth
		y2, y3 = viewHeight, viewHeight*2
		if isOdd {
			y3++
		}
	}

	if v, err := g.SetView(Current, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
	}

	if v, err := g.SetView(Foreign, x2, y2, x3, y3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
	}

	return nil
}

func makeOverviewPanel(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewHeight := maxY - inputHeight
	branchViewWidth := (maxX / 5) * 2

	if v, err := g.SetView(Panel, branchViewWidth*2, 0, maxX-2, viewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Conflicts"
	}
	return nil
}

func makePrompt(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	inputHeight := 2
	viewHeight := maxY - inputHeight

	// Prompt view
	if v, err := g.SetView(Prompt, 0, viewHeight, 21, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		PrintPrompt(g)
	}

	// User input view
	if v, err := g.SetView(Input, 17, viewHeight, maxX, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Editable = true
		v.Wrap = false
		if _, err := g.SetCurrentView(Input); err != nil {
			return err
		}
	}
	return nil
}

func Select(c *conflict.Conflict, g *gocui.Gui, showHelp bool) error {
	// Update side panel
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Panel)
		if err != nil {
			return err
		}
		v.Clear()

		for idx, conflict := range conflicts {
			var out string
			if conflict.Choice != 0 {
				out = color.Green(color.Regular, "✔ %s:%d", conflict.File.Name, conflict.Start)
			} else {
				out = color.Red(color.Regular, "%d. %s:%d", idx+1, conflict.File.Name, conflict.Start)
			}

			if conflict.Equal(c) {
				fmt.Fprintf(v, "-> %s\n", out)
			} else {
				fmt.Fprintf(v, "%s\n", out)
			}
		}

		if showHelp {
			printHelp(v)
		}
		return nil
	})

	// Update code view
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Current)
		if err != nil {
			return err
		}
		v.Title = fmt.Sprintf("%s %s", c.CurrentName, "(Local Version)")

		top, bottom := c.PaddingLines()
		v.Clear()
		printLines(v, top)
		printLines(v, c.ColoredLocalLines)
		printLines(v, bottom)
		if c.Choice == conflict.Local {
			v.FgColor = gocui.ColorGreen
		}

		v, err = g.View(Foreign)
		if err != nil {
			return err
		}
		v.Title = fmt.Sprintf("%s %s", c.ForeignName, "(Incoming Version)")

		top, bottom = c.PaddingLines()
		v.Clear()
		printLines(v, top)
		printLines(v, c.ColoredIncomingLines)
		printLines(v, bottom)
		if c.Choice == conflict.Incoming {
			v.FgColor = gocui.ColorGreen
		}

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

		if cur >= len(conflicts) {
			cur = 0
		} else if cur < 0 {
			cur = len(conflicts) - 1
		}

		if conflicts[cur].Choice == 0 || originalCur == cur {
			break
		}
	}

	if originalCur == cur && conflicts[cur].Choice != 0 {
		globalQuit(g)
	}

	Select(conflicts[cur], g, false)
	return nil
}

func Scroll(g *gocui.Gui, c *conflict.Conflict, direction int) {
	if direction == Up {
		c.TopPeek--
	} else if direction == Down {
		c.TopPeek++
	} else {
		return
	}

	Select(c, g, false)
}
