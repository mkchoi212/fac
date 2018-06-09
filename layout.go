package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
)

// Following constants define the string literal names of 5 views
// that are instantiated via gocui
const (
	Current = "current"
	Foreign = "foreign"
	Panel   = "panel"
	Prompt  = "prompt"
	Input   = "input"
)

// `Up` and `Down` represent scrolling directions
// `Horizontal` and `Vertical` represent current code view orientation
// Notice how both pairs of directionality are `not`s of each other
const (
	Up   = 1
	Down = ^1

	Horizontal = 2
	Vertical   = ^2
)

// Following constants define input panel's dimensions
const (
	inputHeight    = 2
	inputCursorPos = 17
	promptWidth    = 21
)

func printLines(v *gocui.View, lines []string) {
	for _, line := range lines {
		fmt.Fprint(v, line)
	}
}

// layout is used as fac's main gocui.Gui manager
func layout(g *gocui.Gui) (err error) {
	if err = makeCodePanels(g); err != nil {
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

// makeCodePanels draws the two panels representing "local" and "incoming" lines of code
// `viewOrientation` is taken into consideration as the panels can either be
//  `Vertical` or `Horizontal`
func makeCodePanels(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewHeight := maxY - inputHeight
	branchViewWidth := (maxX / 5) * 2
	isOdd := maxY%2 == 1

	var x0, x1, y0, y1 int
	var x2, x3, y2, y3 int

	if viewOrientation == Horizontal {
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

// makeOverviewPanel draws the panel on the right-side of the CUI
// listing all the conflicts that need to be resolved
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

// makePrompt draws two panels on the bottom of the CUI
// A "prompt view" which prompts the user for available keybindings and
// a "user input view" which is an area where the user can type in queries
func makePrompt(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewHeight := maxY - inputHeight

	// Prompt view
	if v, err := g.SetView(Prompt, 0, viewHeight, promptWidth, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		PrintPrompt(g)
	}

	// User input view
	if v, err := g.SetView(Input, inputCursorPos, viewHeight, maxX, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Editable = true
		v.Wrap = false
		v.Editor = gocui.EditorFunc(PromptEditor)
		if _, err := g.SetCurrentView(Input); err != nil {
			return err
		}
	}
	return nil
}

// Select selects conflict `c` as the current conflict displayed on the screen
// When selecting a conflict, it updates the side panel, and the code view
func Select(g *gocui.Gui, c *conflict.Conflict, showHelp bool) error {
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
			PrintHelp(v, &keyBinding)
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

// Resolve resolves the provided conflict and moves to the next conflict
// in the queue
func Resolve(g *gocui.Gui, v *gocui.View, c *conflict.Conflict, version int) error {
	g.Update(func(g *gocui.Gui) error {
		c.Choice = version
		Move(g, v, Down)
		return nil
	})
	return nil
}

// Move goes to the next conflict in the list in the provided `direction`
func Move(g *gocui.Gui, v *gocui.View, direction int) error {
	originalCur := cur

	for {
		if direction == Up {
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

	// Quit application if all items are resolved
	if originalCur == cur && conflicts[cur].Choice != 0 {
		globalQuit(g)
	}

	Select(g, conflicts[cur], false)
	return nil
}

// Scroll scrolls the two code view panels in `direction` by one line
func Scroll(g *gocui.Gui, c *conflict.Conflict, direction int) {
	if direction == Up {
		c.TopPeek--
	} else if direction == Down {
		c.TopPeek++
	} else {
		return
	}

	Select(g, c, false)
}
