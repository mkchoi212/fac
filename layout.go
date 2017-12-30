package main

import (
	"github.com/jroimartin/gocui"
)

const (
	Current = "current"
	Foreign = "foreign"
	Panel   = "panel"
	Prompt  = "prompt"
	Input   = "input prompt"
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
		prompt := Colorize("[wasd] >>", Green)
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
