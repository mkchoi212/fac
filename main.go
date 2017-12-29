// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	curIdx        = 0
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
	v.Clear()
	v.SetCursor(0, 0)
	_ = v.Buffer()

	switch {
	default:
		printPrompt(g, Red("[a | d] >>"))
	}
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
	conflicts[0].Diff()

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

	g.Update(func(g *gocui.Gui) error {
		conflicts[0].Select(g)
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
