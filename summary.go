package main

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
)

func printHelp(v *gocui.View) {
	instruction := `
	w - show more lines up
	s - show more lines down
	a - use local version
	d - use incoming version

	j - scroll down
	k - scroll up

	v - [v]iew orientation
	n - [n]ext conflict
	p - [p]revious conflict

	h | ? - [h]elp
	q | Ctrl+c - [q]uit
	`
	fmt.Fprintf(v, color.Blue(color.Regular, instruction))
}

func printSummary() {
	resolvedCnt := 0
	var line string

	for _, c := range conflicts {
		if c.Choice != 0 {
			line = color.Green(color.Regular, "✔ %s: %d", c.File.Name, c.Start)
			resolvedCnt++
		} else {
			line = color.Red(color.Regular, "✘ %s: %d", c.File.Name, c.Start)
		}
		fmt.Println(line)
	}

	var buf bytes.Buffer
	if resolvedCnt != len(conflicts) {
		buf.WriteString("\nResolved ")
		buf.WriteString(color.Red(color.Light, "%d ", resolvedCnt))
		buf.WriteString("conflict(s) out of ")
		buf.WriteString(color.Red(color.Light, "%d", len(conflicts)))
	} else {
		buf.WriteString(color.Green(color.Regular, "\nFixed All Conflicts 🎉"))
	}
	fmt.Println(buf.String())
}
