// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
)

var (
	cur              = 0
	consecutiveError = 0
)

func printLines(v *gocui.View, lines []string) {
	for _, line := range lines {
		fmt.Fprint(v, line)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func parseInput(g *gocui.Gui, v *gocui.View) error {
	evalCmd := func(in rune, g *gocui.Gui) {
		switch in {
		case 'j':
			Scroll(g, &conflict.All[cur], Up)
		case 'k':
			Scroll(g, &conflict.All[cur], Down)
		case 'w':
			conflict.All[cur].TopPeek++
			Select(&conflict.All[cur], g, false)
		case 's':
			conflict.All[cur].BottomPeek++
			Select(&conflict.All[cur], g, false)
		case 'a':
			Resolve(&conflict.All[cur], g, v, Local)
		case 'd':
			Resolve(&conflict.All[cur], g, v, Incoming)
		case 'n':
			MoveToItem(Down, g, v)
		case 'p':
			MoveToItem(Up, g, v)
		case 'v':
			ViewOrientation = ^ViewOrientation
			layout(g)
		case 'h', '?':
			Select(&conflict.All[cur], g, true)
		case 'q':
			globalQuit(g)
		// case 'z':
		//	conflict.All[cur].ToggleDiff()
		//	Select(&conflict.All[cur], g, false)
		default:
			PrintPrompt(g, color.Red(color.Regular, "[wasd] >>"))
			consecutiveError++
			return
		}
		PrintPrompt(g, color.Green(color.Regular, "[wasd] >>"))

		if consecutiveError == 2 {
			consecutiveError = 0
			Select(&conflict.All[cur], g, true)
		}
	}

	in := strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)

	if len(in) > 1 {
		for _, r := range [...]rune{'a', 'd', 'h', 'z'} {
			if strings.ContainsRune(in, r) {
				PrintPrompt(g, color.Red(color.Regular, "[wasd] >>"))
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
	if err := conflict.Find(); err != nil {
		switch err.(type) {
		case *conflict.ErrNoConflict:
			fmt.Println(color.Green(color.Regular, err.Error()))
		default:
			fmt.Print(color.Red(color.Regular, err.Error()))
		}
		return
	}

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

	Select(&conflict.All[0], g, false)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	g.Close()

	for fname := range conflict.FileLines {
		if err := FinalizeChanges(fname); err != nil {
			fmt.Println(color.Red(color.Underline, "%s\n", err))
		}
	}
	printSummary()
}
