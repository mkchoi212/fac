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

	j - scroll down
	k - scroll up

	h | ? - print help
	q | Ctrl+c - quit
	`
	fmt.Fprintf(v, Black(Regular, instruction))
}

func printSummary() {
	resolvedCnt := 0
	var line string

	for _, c := range conflicts {
		if c.Resolved {
			line = Green(Regular, "✔ %s: %d", c.FileName, c.Start)
			resolvedCnt++
		} else {
			line = Red(Regular, "✘ %s: %d", c.FileName, c.Start)
		}
		fmt.Println(line)
	}

	var buf bytes.Buffer
	if resolvedCnt != len(conflicts) {
		buf.WriteString("\nResolved ")
		buf.WriteString(Red(Light, "%d ", resolvedCnt))
		buf.WriteString("conflict(s) out of ")
		buf.WriteString(Red(Light, "%d", len(conflicts)))
	} else {
		buf.WriteString(Green(Regular, "\nFixed All Conflicts 🎉"))
	}
	fmt.Println(buf.String())
}
