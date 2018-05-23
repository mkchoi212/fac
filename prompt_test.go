package main

import (
	"testing"

	"github.com/jroimartin/gocui"
)

func TestPrintPrompt(t *testing.T) {
	g := gocui.Gui{}
	makePrompt(&g)
	PrintPrompt(&g)

	_, err := g.View(Prompt)
	if err != nil {
		t.Fatalf("PrintPrompt failed: %s", err.Error())
	}
}
