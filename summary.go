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
	z - toggle diff

	j - scroll down
	k - scroll up

	h | ? - print help
	q | Ctrl+c - quit
	`
	fmt.Fprintf(v, Blue(Regular, instruction))
}

func printSummary() {
	resolvedCnt := 0
	var line string

	for _, c := range conflicts {
		if c.Choice != 0 {
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

func conflictsIn(fname string) (list []Conflict) {
	for _, c := range conflicts {
		if c.AbsolutePath == fname && c.Choice != 0 {
			list = append(list, c)
		}
	}
	return
}

func WriteChanges(fname string) error {
	targetConflicts := conflictsIn(fname)

	var replacementLines []string

	for _, c := range targetConflicts {
		if c.Choice == Local {
			replacementLines = append([]string{}, c.CurrentLines...)
		} else {
			replacementLines = append([]string{}, c.ForeignLines...)
		}

		i := 0
		for ; i < len(replacementLines); i++ {
			allFileLines[fname][c.Start+i-1] = replacementLines[i]
		}
		for ; c.End-c.Start >= i; i++ {
			allFileLines[fname][c.Start+i-1] = ""
		}
	}

	for _, l := range allFileLines[fname] {
		if l != "" {
			fmt.Printf(l)
		}
	}
	return nil
}
