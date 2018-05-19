package conflict

import (
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}
		if err := f.Read(); err != nil {
			t.Error("Read failed: could not read file")
		}

		if len(f.Lines) != test.numLines {
			t.Errorf("Read failed: got %d lines, wanted %d lines", len(f.Lines), test.numLines)
		}
	}
}

func TestWriteChanges(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}
		if err := f.Read(); err != nil {
			t.Fatalf("WriteChanges/Read failed")
		}

		conflicts, err := parseConflictsIn(f, test.markers)
		if err != nil {
			if !test.shouldPass {
				continue
			}
			t.Fatal("WriteChanges/parseConflicts failed")
		}

		for i := range test.resolveDecision {
			conflicts[i].Choice = test.resolveDecision[i]
		}
		f.Conflicts = conflicts

		f.AbsolutePath = "testdata/.test_output"
		if err := f.WriteChanges(); err != nil {
			t.Errorf("WriteChages failed: %s", err.Error())
		}

		f.Lines = nil
		if err := f.Read(); err != nil {
			t.Error("WriteChanges/Read failed")
		}

		for i := range f.Lines {
			f.Lines[i] = strings.TrimSuffix(f.Lines[i], "\n")
		}

		if !(reflect.DeepEqual(f.Lines, test.resolved)) {
			t.Errorf("WriteChanges failed: got %v, want %v", f.Lines, test.resolved)
		}
	}
}
