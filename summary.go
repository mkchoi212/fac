package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

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

func writeChanges(absPath string) (err error) {
	f, err := os.Create(absPath)
	if err != nil {
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range allFileLines[absPath] {
		if _, err = w.WriteString(line); err != nil {
			return
		}
	}
	err = w.Flush()
	return
}

func FinalizeChanges(absPath string) (err error) {
	targetConflicts := conflictsIn(absPath)

	var replacementLines []string

	for _, c := range targetConflicts {
		if c.Choice == Local {
			replacementLines = append([]string{}, c.CurrentLines...)
		} else {
			replacementLines = append([]string{}, c.ForeignLines...)
		}

		i := 0
		for ; i < len(replacementLines); i++ {
			allFileLines[absPath][c.Start+i-1] = replacementLines[i]
		}
		for ; c.End-c.Start >= i; i++ {
			allFileLines[absPath][c.Start+i-1] = ""
		}
	}

	if err = writeChanges(absPath); err != nil {
		return
	}
	return
}
