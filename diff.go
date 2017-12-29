package main

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func (c *Conflict) Diff() {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(strings.Join(c.CurrentLines, "\n"), strings.Join(c.ForeignLines, "\n"), false)
	fmt.Println(dmp.DiffPrettyText(diffs))
}
