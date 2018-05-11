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
	all              = []conflict.Conflict{}
	numConflicts     = 0
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
			Scroll(g, &all[cur], Up)
		case 'k':
			Scroll(g, &all[cur], Down)
		case 'w':
			all[cur].TopPeek++
			Select(&all[cur], g, false)
		case 's':
			all[cur].BottomPeek++
			Select(&all[cur], g, false)
		case 'a':
			Resolve(&all[cur], g, v, Local)
		case 'd':
			Resolve(&all[cur], g, v, Incoming)
		case 'n':
			MoveToItem(Down, g, v)
		case 'p':
			MoveToItem(Up, g, v)
		case 'v':
			ViewOrientation = ^ViewOrientation
			layout(g)
		case 'h', '?':
			Select(&all[cur], g, true)
		case 'q':
			globalQuit(g)
		default:
			PrintPrompt(g, color.Red(color.Regular, promptString))
			consecutiveError++
			if consecutiveError == 2 {
				consecutiveError = 0
				Select(&all[cur], g, true)
			}
			return
		}
		consecutiveError = 0
		PrintPrompt(g, color.Green(color.Regular, promptString))
	}

	in := strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)

	if len(in) > 1 {
		for _, r := range [...]rune{'a', 'd', 'h', 'z'} {
			if strings.ContainsRune(in, r) {
				PrintPrompt(g, color.Red(color.Regular, promptString))
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
	//cwd, _ := os.Getwd()
	cwd := "./test"
	conflicts, err := conflict.Find(cwd)
	if err != nil {
		fmt.Println(color.Red(color.Regular, err.Error()))
		return
	}
	if len(conflicts) == 0 {
		fmt.Println(color.Green(color.Regular, "No conflicts detected 🎉"))
		return
	}

	all = conflicts
	numConflicts = len(conflicts)

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
		log.Panicln(err)
	}

	Select(&all[0], g, false)

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
