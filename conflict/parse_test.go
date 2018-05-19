package conflict

import (
	"os/exec"
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

func TestFind(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	files, err := Find(".")
	if err != nil {
		t.Errorf("Find failed: %s", err.Error())
	}

	if len(files) != 3 {
		t.Errorf("Find failed: got %d files, want 3", len(files))
	}

	for _, f := range files {
		for _, test := range tests {
			if f.AbsolutePath == test.path && test.parsable {
				if len(f.Conflicts) != test.numConflicts {
					t.Errorf("Find failed: got %d conflicts in %s, want %d",
						len(f.Conflicts), test.path, test.numConflicts)
				}
			}
		}
	}
}
