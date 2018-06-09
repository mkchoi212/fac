package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mkchoi212/fac/binding"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
)

// PrintHelp prints the current key binding rules in the side panel
func PrintHelp(v io.Writer, binding *binding.Binding) {
	fmt.Fprintf(v, color.Blue(color.Regular, binding.Help()))
}

// PrintSummary prints the summary of the fac session after the user
// either quits the program or has resolved all conflicts
func PrintSummary(conflicts []*conflict.Conflict) {
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
