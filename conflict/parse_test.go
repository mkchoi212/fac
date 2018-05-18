package conflict

import (
	"reflect"
	"testing"
)

func TestParseGitInfo(t *testing.T) {
	for _, test := range tests {
		output := []int{}
		for _, line := range test.diffCheck {
			_, lineNum, ok := parseGitMarkerInfo(line)
			if ok != nil {
				continue
			}
			output = append(output, lineNum)
		}

		if !(reflect.DeepEqual(output, test.markers)) && test.parsable {
			t.Errorf("parseGitInfo failed: got %v, want %v", output, test.markers)
		}
	}
}

func TestParseConflictsIn(t *testing.T) {
	for _, test := range tests {
		f := File{AbsolutePath: test.path}
		if err := f.Read(); err != nil && test.highlightable {
			t.Error("ParseConflicts/Read failed")
		}

		_, err := parseConflictsIn(f, test.markers)
		if err != nil && test.parsable {
			t.Errorf("parseConflicts failed: %s", err.Error())
		}
	}
}
