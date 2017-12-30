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
)

func printLines(v *gocui.View, lines []string) {
	v.Clear()
	for _, line := range lines {
		fmt.Fprintln(v, line)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func parseInput(g *gocui.Gui, v *gocui.View) error {
	in := strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)

	switch {
	case in == "a":
		conflicts[cur].Resolve(g, v, Local)
	case in == "d":
		conflicts[cur].Resolve(g, v, Incoming)
	case in == "h":
		conflicts[cur].Select(g, true)
	case in == "q":
		globalQuit(g)
	default:
		printPrompt(g, Colorize("[wasd] >>", Red))
		return nil
	}
	printPrompt(g, Colorize("[wasd] >>", Green))
	return nil
}

func keyBindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, parseInput)
	return err
}

func main() {
	var err error
	conflicts, err = Find()
	if err != nil {
		log.Panicln("No conflicts found")
	}
	conflictCount = len(conflicts)

	printSummary()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	g.Cursor = true

	if err := keyBindings(g); err != nil {
		log.Panicln(err)
	}

	conflicts[0].Select(g, false)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
