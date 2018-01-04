package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/conflict"
	"github.com/mkchoi212/fac/style"
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
	fmt.Fprintf(v, style.Blue(instruction))
}

func printSummary() {
	resolvedCnt := 0
	var line string

	for _, c := range conflict.All {
		if c.Choice != 0 {
			line = style.Green("✔ %s: %d", c.FileName, c.Start)
			resolvedCnt++
		} else {
			line = style.Red("✘ %s: %d", c.FileName, c.Start)
		}
		fmt.Println(line)
	}

	var buf bytes.Buffer
	if resolvedCnt != conflict.Count {
		buf.WriteString("\nResolved ")
		buf.WriteString(style.RedLight("%d ", resolvedCnt))
		buf.WriteString("conflict(s) out of ")
		buf.WriteString(style.RedLight("%d", conflict.Count))
	} else {
		buf.WriteString(style.Green("\nFixed All Conflicts 🎉"))
	}
	fmt.Println(buf.String())
}

func writeChanges(absPath string) (err error) {
	f, err := os.Create(absPath)
	if err != nil {
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range conflict.FileLines[absPath] {
		if _, err = w.WriteString(line); err != nil {
			return
		}
	}
	err = w.Flush()
	return
}

func FinalizeChanges(absPath string) (err error) {
	targetConflicts := conflict.In(absPath)

	var replacementLines []string

	for _, c := range targetConflicts {
		if c.Choice == Local {
			replacementLines = append([]string{}, c.LocalLines...)
		} else {
			replacementLines = append([]string{}, c.IncomingLines...)
		}

		i := 0
		for ; i < len(replacementLines); i++ {
			conflict.FileLines[absPath][c.Start+i-1] = replacementLines[i]
		}
		for ; c.End-c.Start >= i; i++ {
			conflict.FileLines[absPath][c.Start+i-1] = ""
		}
	}

	if err = writeChanges(absPath); err != nil {
		return
	}
	return
}
