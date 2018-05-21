package main

import (
	"testing"

	"github.com/jroimartin/gocui"
	"github.com/mkchoi212/fac/color"
)

func TestPrintPrompt(t *testing.T) {
	g := gocui.Gui{}
	makePrompt(&g)
	PrintPrompt(&g, color.Red)

	_, err := g.View(Prompt)
	if err != nil {
		t.Fatalf("PrintPrompt failed: %s", err.Error())
	}
}
