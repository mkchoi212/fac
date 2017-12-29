package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func printHelp(v *gocui.View) {
	instruction := `
	w - show more lines up
	s - show more lines down
	a - select left screen
	d - select right screen

	h - print help
	Ctrl+c - quit application
	`
	fmt.Fprintf(v, Colorize(instruction, Purple))
}
