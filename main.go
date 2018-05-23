package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
	"github.com/mkchoi212/fac/editor"
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

func parseInput(g *gocui.Gui, v *gocui.View) error {
	in := strings.TrimSuffix(v.Buffer(), "\n")
	v.Clear()
	v.SetCursor(0, 0)

	if err := Evaluate(g, v, conflicts[cur], in); err != nil {
		if err == ErrUnknownCmd {
			consecutiveError++
			if consecutiveError > 3 {
				Select(conflicts[cur], g, true)
			}
		} else {
			return err
		}
	} else {
		consecutiveError = 0
	}

	PrintPrompt(g)
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

func findConflicts() (files []conflict.File, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	if files, err = conflict.Find(cwd); err != nil {
		return
	}

	for i := range files {
		file := &files[i]
		for j := range file.Conflicts {
			conflicts = append(conflicts, &file.Conflicts[j])
		}
	}

	return
}

func runUI() error {
	for {
		if err := Start(); err != nil {
			if err == ErrOpenEditor {
				newLines, err := editor.Open(conflicts[cur])
				if err != nil {
					return err
				}
				if err = conflicts[cur].Update(newLines); err != nil {
					consecutiveError++
				}
			} else if err == gocui.ErrQuit {
				break
			}
		}
	}

	return nil
}

func die(err error) {
	fmt.Println(color.Red(color.Regular, "fac: %s", strings.TrimSuffix(err.Error(), "\n")))
	os.Exit(1)
}

func main() {
	// Find and parse conflicts
	files, err := findConflicts()
	if err != nil {
		die(err)
	}

	if len(conflicts) == 0 {
		fmt.Println(color.Green(color.Regular, "No conflicts detected 🎉"))
		os.Exit(0)
	}

	if err = runUI(); err != nil {
		die(err)
	}

	for _, file := range files {
		if err = file.WriteChanges(); err != nil {
			die(err)
		}
	}

	printSummary(conflicts)
}
