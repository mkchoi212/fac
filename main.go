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
	conflicts        = []*conflict.Conflict{}
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
			Scroll(g, conflicts[cur], Up)
		case 'k':
			Scroll(g, conflicts[cur], Down)
		case 'w':
			conflicts[cur].TopPeek++
			Select(conflicts[cur], g, false)
		case 's':
			conflicts[cur].BottomPeek++
			Select(conflicts[cur], g, false)
		case 'a':
			Resolve(conflicts[cur], g, v, conflict.Local)
		case 'd':
			Resolve(conflicts[cur], g, v, conflict.Incoming)
		case 'n':
			MoveToItem(Down, g, v)
		case 'p':
			MoveToItem(Up, g, v)
		case 'v':
			ViewOrientation = ^ViewOrientation
			layout(g)
		case 'h', '?':
			Select(conflicts[cur], g, true)
		case 'q':
			globalQuit(g)
		default:
			PrintPrompt(g, color.Red(color.Regular, promptString))
			consecutiveError++
			if consecutiveError == 2 {
				consecutiveError = 0
				Select(conflicts[cur], g, true)
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
	// Find and parse conflicts
	//cwd, _ := os.Getwd()
	cwd := "../dummy_repo"
	files, err := conflict.Find(cwd)
	if err != nil {
		fmt.Println(color.Red(color.Regular, err.Error()))
		return
	}
	if len(files) == 0 {
		fmt.Println(color.Green(color.Regular, "No conflicts detected 🎉"))
		return
	}

	for i := range files {
		file := files[i]
		for j := range file.Conflicts {
			conflicts = append(conflicts, &file.Conflicts[j])
		}
	}

	// Setup CUI Environment
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

	// Main UI loop
	Select(conflicts[0], g, false)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	g.Close()

	for _, file := range files {
		if err = file.WriteChanges(); err != nil {
			fmt.Println(color.Red(color.Underline, "%s\n", err))
		}
	}
	printSummary()
}
