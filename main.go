// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	current       = 0
	conflictCount = 0
	conflicts     = []Conflict{}
)

func printLines(v *gocui.View, lines []string) {
	v.Clear()
	for _, line := range lines {
		fmt.Fprintln(v, line)
	}
}

func selectConflict(i int, g *gocui.Gui) error {

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("panel")
		if err != nil {
			return err
		}
		v.Clear()

		for idx, conflict := range conflicts {
			var out string
			if conflict.Resolved {
				out = Green(fmt.Sprintf("✅  %s:%d", conflict.FileName, conflict.Start))
			} else {
				out = Red(fmt.Sprintf("%d.  %s:%d", idx+1, conflict.FileName, conflict.Start))
			}

			if idx == i {
				fmt.Fprintf(v, "%s <-\n", out)
			} else {
				fmt.Fprintf(v, "%s\n", out)
			}
		}
		return nil
	})

	g.Update(func(g *gocui.Gui) error {
		conf := conflicts[i]

		v, err := g.View("current")
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.WriteString(conf.CurrentName)
		buf.WriteString(" (Current Change) ")
		v.Title = buf.String()

		printLines(v, conf.ColoredCurrentLines)

		v, err = g.View("foreign")
		if err != nil {
			return err
		}
		buf.Reset()
		buf.WriteString(conf.ForeignName)
		buf.WriteString(" (Incoming Change) ")
		v.Title = buf.String()

		printLines(v, conf.ColoredForeignLines)
		return nil
	})

	return nil
}

func nextConflict(g *gocui.Gui, v *gocui.View) error {
	current = current + 1
	if current >= conflictCount {
		current = 0
	}

	selectConflict(current, g)
	return nil
}

func resolveConflict(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		conflicts[current].Resolve()
		selectConflict(current, g)
		return nil
	})
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func parseInput(g *gocui.Gui, v *gocui.View) error {
	g.Update(updatePrompt(v))
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
		selectConflict(0, g)
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
