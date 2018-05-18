package conflict

import (
	"reflect"
	"testing"
)

var conflictFile = struct {
	readPath  string
	writePath string
	markers   []int
	lineNum   int
}{
	readPath:  "testdata/CircularCrownSelector.swift",
	writePath: "testdata/output.swift",
	markers:   []int{145, 153, 161, 169},
	lineNum:   171,
}

func TestRead(t *testing.T) {
	f := File{AbsolutePath: conflictFile.readPath}
	if err := f.Read(); err != nil {
		t.Error("Read failed: could not read file")
	}

	if len(f.Lines) != conflictFile.lineNum {
		t.Errorf("Read failed: got %d lines, wanted %d lines", len(f.Lines), conflictFile.lineNum)
	}
}

func TestWriteChanges(t *testing.T) {
	f := File{AbsolutePath: conflictFile.readPath}
	if err := f.Read(); err != nil {
		t.Error("WriteChanges/Read failed")
	}

	conflicts, err := parseConflictsIn(f, conflictFile.markers)
	if err != nil {
		t.Error("WriteChanges/parseConflicts failed")
	}

	f.Conflicts = conflicts
	targetConflict := &f.Conflicts[0]
	targetConflict.Choice = Local

	f.AbsolutePath = conflictFile.writePath
	if err := f.WriteChanges(); err != nil {
		t.Errorf("WriteChages failed: %s", err.Error())
	}

	expected := f.Lines[142:153]
	f.Lines = nil
	if err := f.Read(); err != nil {
		t.Error("WriteChanges/Read failed")
	}

	output := f.Lines[142:]
	if reflect.DeepEqual(output, expected) {
	}
}
