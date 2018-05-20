package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
)

func globalQuit(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("")
		return quit(g, v)
	})
}

// PrintPrompt prints the promptString on the bottom left corner of the screen
// Note that the prompt is composed of two seperate views
// One that displays just the promptString, and another that takes input from the user
func PrintPrompt(g *gocui.Gui, colorize func(style int, format string, a ...interface{}) string) {
	promptString := "[w,a,s,d,?] >>"

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Prompt)
		if err != nil {
			return err
		}
		v.Clear()
		v.MoveCursor(0, 0, true)
		fmt.Fprintf(v, colorize(color.Regular, promptString))
		return nil
	})
}
