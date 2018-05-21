package main

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
	"github.com/mkchoi212/fac/editor"
)

// ErrUnknownCmd is returned when user inputs an invalid character
var ErrUnknownCmd = errors.New("This person doesn't know whats going on")

// ErrNeedRefresh is returned when the user opens vim to edit code
// Note that a new instance of gocui must be created after vim is opened
var ErrNeedRefresh = errors.New("Screen is tainted after opening vim")

// PrintPrompt prints the promptString on the bottom left corner of the screen
// Note that the prompt is composed of two seperate views,
// one that displays just the promptString, and another that takes input from the user
func PrintPrompt(g *gocui.Gui, colorize func(style int, format string, a ...interface{}) string) {
	promptString := "[w,a,s,d,e,?] >>"

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Prompt)
		if err != nil {
			return err
		}
		v.Clear()
		v.MoveCursor(0, 0, true)
		fmt.Fprintf(v, colorize(color.Regular, promptString))
		return nil
	})
}

// Evaluate evalutes the user's input character by character
// It returns `ErrUnknownCmd` if the string contains an invalid command
// It also returns `ErrNeedRefresh` if user uses `e` command to open vim
func Evaluate(g *gocui.Gui, v *gocui.View, conf *conflict.Conflict, input string) (err error) {
	for _, c := range input {
		switch c {
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
		case 'e':
			newLines, err := editor.Open(conflicts[cur])
			if err != nil {
				return err
			}
			conflicts[cur].Update(newLines)
			return ErrNeedRefresh
		case 'h', '?':
			Select(conflicts[cur], g, true)
		case 'q':
			globalQuit(g)
		default:
			return ErrUnknownCmd
		}
	}
	return
}
