// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	cur           = 0
	conflictCount = 0
	conflicts     = []Conflict{}

	consecutiveError = 0
)

func printLines(v *gocui.View, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(v, line)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func parseInput(g *gocui.Gui, v *gocui.View) error {
	evalCmd := func(in rune, g *gocui.Gui) {
		switch {
		case in == 'j':
			scroll(g, &conflicts[cur], Up)
		case in == 'k':
			scroll(g, &conflicts[cur], Down)
		case in == 'w':
			conflicts[cur].topPeek++
			conflicts[cur].Select(g, false)
		case in == 's':
			conflicts[cur].bottomPeek++
			conflicts[cur].Select(g, false)
		case in == 'a':
			conflicts[cur].Resolve(g, v, Local)
		case in == 'd':
			conflicts[cur].Resolve(g, v, Incoming)
		case in == 'h' || in == '?':
			conflicts[cur].Select(g, true)
		case in == 'q':
			globalQuit(g)
		case in == 'z':
			conflicts[cur].toggleDiff()
			conflicts[cur].Select(g, false)
		default:
			printPrompt(g, Red(Regular, "[wasd] >>"))
			consecutiveError++
		}
		if consecutiveError == 2 {
			consecutiveError = 0
			conflicts[cur].Select(g, false)
		}
	}

	in := strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)

	if len(in) > 1 {
		for _, r := range [...]rune{'a', 'd', 'h', 'z'} {
			if strings.ContainsRune(in, r) {
				printPrompt(g, Red(Regular, "[wasd] >>"))
				return nil
			}
		}
	}

	for _, c := range in {
		evalCmd(c, g)
	}
	return nil
}

func main() {
	var err error
	conflicts, err = Find()
	if err != nil {
		log.Panicln("No conflicts found")
	}
	conflictCount = len(conflicts)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	g.SetManagerFunc(layout)
	g.Cursor = true

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, parseInput); err != nil {
		log.Panic(err)
	}

	conflicts[0].Select(g, false)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	g.Close()
	printSummary()
}
