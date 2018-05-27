package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
	"github.com/mkchoi212/fac/key"
)

func printHelp(v io.Writer, binding *key.Binding) {
	fmt.Fprintf(v, color.Blue(color.Regular, binding.Help()))
}

func printSummary(conflicts []*conflict.Conflict) {
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
