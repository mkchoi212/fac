package conflict

import (
	"testing"
)

var commands = []struct {
	command string
	ok      bool
}{
	{"time", true},
	{"foobar", false},
}

func TestRun(t *testing.T) {
}
