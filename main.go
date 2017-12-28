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
	current       = 0
	conflictCount = 0
	conflicts     = []Conflict{}
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	height := maxY - 2
	branchViewWidth := (maxX / 5) * 2

	if _, err := g.SetView("current", 0, 0, branchViewWidth, height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if _, err := g.SetView("foreign", branchViewWidth, 0, branchViewWidth*2, height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if v, err := g.SetView("panel", branchViewWidth*2, 0, maxX-2, height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Conflicts"
	}

	return nil
}

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
				out = fmt.Sprintf("✅  \033[3%d;%dm%s:%d \033[0m", 2, 1, conflict.FileName, conflict.Start)
			} else {
				out = fmt.Sprintf("%d. \033[3%d;%dm%s:%d \033[0m", idx+1, 1, 1, conflict.FileName, conflict.Start)
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
		v.Title = conf.CurrentName
		printLines(v, conf.ColoredCurrentLines)

		v, err = g.View("foreign")
		if err != nil {
			return err
		}
		v.Title = conf.ForeignName
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

func keyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextConflict); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, resolveConflict); err != nil {
		return err
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
	defer g.Close()

	g.SetManagerFunc(layout)

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
