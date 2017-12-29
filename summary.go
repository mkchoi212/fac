package main

import (
	"bytes"
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

func printSummary() {
	resolvedCnt := 0
	var line string

	for _, c := range conflicts {
		if c.Resolved {
			line = Colorize(fmt.Sprintf("✔ %s: %d", c.FileName, c.Start), Green)
			resolvedCnt++
		} else {
			line = Colorize(fmt.Sprintf("✘ %s: %d", c.FileName, c.Start), Red)
		}
		fmt.Println(line)
	}

	var buf bytes.Buffer
	if resolvedCnt != len(conflicts) {
		buf.WriteString("\nResolved ")
		buf.WriteString(ColorizeLight(fmt.Sprintf("%d ", resolvedCnt), Red))
		buf.WriteString("conflict(s) out of ")
		buf.WriteString(ColorizeLight(fmt.Sprintf("%d", len(conflicts)), Red))
	} else {
		buf.WriteString(Colorize(fmt.Sprintf("\nFixed All Conflicts 🎉"), Green))
	}
	fmt.Println(buf.String())
}
