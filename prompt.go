package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mkchoi212/fac/key"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
	"github.com/mkchoi212/fac/conflict"
)

// ErrUnknownCmd is returned when user inputs an invalid character
var ErrUnknownCmd = errors.New("This person doesn't know whats going on")

// ErrOpenEditor is returned when the user wants to open an editor
// Note that the current instance of gocui must be destroyed before opening an editor
var ErrOpenEditor = errors.New("Screen is tainted after opening vim")

// PrintPrompt prints the promptString on the bottom left corner of the screen
// Note that the prompt is composed of two separate views,
// one that displays just the promptString, and another that takes input from the user
func PrintPrompt(g *gocui.Gui) {
	promptString := binding.Summary()

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(Prompt)
		if err != nil {
			return err
		}
		v.Clear()
		v.MoveCursor(0, 0, true)

		if consecutiveErrCnt == 0 {
			fmt.Fprintf(v, color.Green(color.Regular, promptString))
		} else {
			fmt.Fprintf(v, color.Red(color.Regular, promptString))
		}
		return nil
	})
}

// Evaluate evalutes the user's input character by character
// It returns `ErrUnknownCmd` if the string contains an invalid command
// It also returns `ErrNeedRefresh` if user uses `e` command to open vim
func Evaluate(g *gocui.Gui, v *gocui.View, conf *conflict.Conflict, input string) (err error) {
	for _, c := range input {
		switch string(c) {
		case binding[key.ScrollUp]:
			Scroll(g, conflicts[cur], Up)
		case binding[key.ScrollDown]:
			Scroll(g, conflicts[cur], Down)
		case binding[key.ShowLinesUp]:
			conflicts[cur].TopPeek++
			Select(g, conflicts[cur], false)
		case binding[key.ShowLinesDown]:
			conflicts[cur].BottomPeek++
			Select(g, conflicts[cur], false)
		case binding[key.SelectLocal]:
			Resolve(g, v, conflicts[cur], conflict.Local)
		case binding[key.SelectIncoming]:
			Resolve(g, v, conflicts[cur], conflict.Incoming)
		case binding[key.NextConflict]:
			Move(g, v, Down)
		case binding[key.PreviousConflict]:
			Move(g, v, Up)
		case binding[key.ToggleViewOrientation]:
			viewOrientation = ^viewOrientation
			layout(g)
		case binding[key.EditCode]:
			return ErrOpenEditor
		case binding[key.ShowHelp], "?":
			Select(g, conflicts[cur], true)
		case binding[key.QuitApplication]:
			globalQuit(g)
		default:
			return ErrUnknownCmd
		}
	}
	return
}

// ParseInput is invoked when the user presses "Enter"
// It `evaluate`s the user's query and reflects the state on the UI
func ParseInput(g *gocui.Gui, v *gocui.View) error {
	in := strings.TrimSuffix(v.Buffer(), "\n")

	if err := Evaluate(g, v, conflicts[cur], in); err != nil {
		if err == ErrUnknownCmd {
			consecutiveErrCnt++
			if consecutiveErrCnt > 3 {
				Select(g, conflicts[cur], true)
			}
		} else {
			return err
		}
	} else {
		consecutiveErrCnt = 0
	}

	PrintPrompt(g)
	return nil
}

// PromptEditor handles user's interaction with the prompt
// Note that user's `ContinuousEvaluation` setting value changes its behavior
func PromptEditor(v *gocui.View, k gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		if key.ContinuousEvaluation == "true" {
			v.Clear()
			v.EditWrite(ch)
			ParseInput(g, v)
			v.SetCursor(0, 0)
		} else {
			v.EditWrite(ch)
		}
		return
	}

	switch k {
	case gocui.KeyEnter:
		ParseInput(g, v)
		v.Clear()
		v.SetCursor(0, 0)
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
	case gocui.KeyDelete:
		v.EditDelete(false)
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case gocui.KeyArrowDown:
		v.SetCursor(len(v.Buffer())-1, 0)
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
