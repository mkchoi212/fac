package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/binding"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
	"github.com/mkchoi212/fac/editor"
)

var (
	viewOrientation   = Vertical
	conflicts         = []*conflict.Conflict{}
	keyBinding        = binding.Binding{}
	cur               = 0
	consecutiveErrCnt = 0
	g                 *gocui.Gui
)

// globalQuit is invoked when the user quits the contact and or
// when all conflicts have been resolved
func globalQuit(g *gocui.Gui, err error) {
	g.Update(func(g *gocui.Gui) error {
		return err
	})
}

// quit is invoked when the user presses "Ctrl+C"
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// findConflicts looks at the current directory and returns an
// array of `File`s that contain merge conflicts
// It returns an error if it fails to parse the conflicts
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
			// Set context lines to show
			c := &file.Conflicts[j]
			if err = c.SetContextLines(keyBinding[binding.DefaultContextLines]); err != nil {
				return
			}
			conflicts = append(conflicts, c)
		}
	}

	return
}

// runUI initializes, configures, and starts a fresh instance of gocui
func runUI() (err error) {
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}

	defer g.Close()
	g.SetManagerFunc(layout)
	g.Cursor = true

	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return
	}

	if keyBinding[binding.ContinuousEvaluation] == "false" {
		if err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, ParseInput); err != nil {
			return
		}
	}

	Select(g, conflicts[cur], false)

	if err = g.MainLoop(); err != nil {
		return
	}

	return
}

// mainLoop manages how the main instances of gocui are created and destroyed
func mainLoop() error {
	for {
		if err := runUI(); err != nil {
			// Instantiates a fresh instance of gocui
			// when user opens an editor because screen is dirty
			if err == ErrOpenEditor {
				newLines, err := editor.Open(conflicts[cur])
				if err != nil {
					return err
				}
				if err = conflicts[cur].Update(newLines); err != nil {
					consecutiveErrCnt++
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
	ParseFlags()

	var err error

	keyBinding, err = binding.LoadSettings()
	if err != nil {
		die(err)
	}

	files, err := findConflicts()
	if err != nil {
		die(err)
	}

	if len(conflicts) == 0 {
		fmt.Println(color.Green(color.Regular, "No conflicts detected 🎉"))
		os.Exit(0)
	}

	if err = mainLoop(); err != nil {
		die(err)
	}

	for _, file := range files {
		if err = file.WriteChanges(); err != nil {
			die(err)
		}
	}

	PrintSummary(conflicts)
}
