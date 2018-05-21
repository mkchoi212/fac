package main

import (
	"fmt"
	"os"
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

func globalQuit(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("")
		return quit(g, v)
	})
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func parseInput(g *gocui.Gui, v *gocui.View) (err error) {
	in := strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)

	if err = Evaluate(g, v, conflicts[cur], in); err != nil {
		if err == ErrUnknownCmd {
			consecutiveError++
			PrintPrompt(g, color.Red)

			if consecutiveError > 3 {
				Select(conflicts[cur], g, true)
				consecutiveError = 0
			}
		} else {
			return
		}
	}

	return nil
}

// Start initializes, configures, and starts a fresh instance of gocui
func Start() (err error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}

	defer g.Close()
	g.SetManagerFunc(layout)
	g.Cursor = true

	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return
	}
	if err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, parseInput); err != nil {
		return
	}

	Select(conflicts[cur], g, false)
	if err = g.MainLoop(); err != nil {
		return
	}

	return
}

func main() {
	// Find and parse conflicts
	cwd, _ := os.Getwd()
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
		file := &files[i]
		for j := range file.Conflicts {
			conflicts = append(conflicts, &file.Conflicts[j])
		}
	}

	// Main GUI loop
	for {
		if err := Start(); err != nil {
			if err == gocui.ErrQuit {
				break
			} else if err == ErrNeedRefresh {
				continue
			}
		}
	}

	for _, file := range files {
		if err = file.WriteChanges(); err != nil {
			fmt.Println(color.Red(color.Underline, "%s\n", err))
		}
	}
	printSummary(conflicts)
}
