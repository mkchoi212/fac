package main

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func (c *Conflict) diff() []string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(strings.Join(c.CurrentLines, "\n"), strings.Join(c.ForeignLines, "\n"), false)
	return []string{dmp.DiffPrettyText(diffs)}
}
