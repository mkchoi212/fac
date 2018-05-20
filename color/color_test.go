package color

import (
	"os"
	"testing"
)

// Setup tests
var tests = []struct {
	color    func(style int, format string, a ...interface{}) string
	input    []string
	expected string
}{
	{Black, []string{"%s", "foobar"}, "\033[30;1mfoobar[0m"},
	{Red, []string{"%s %s", "foobar", "hey"}, "\033[31;1mfoobar hey[0m"},
	{Green, []string{"%s", ""}, "\033[32;1m[0m"},
	{Blue, []string{"foobar"}, "\033[34;1mfoobar[0m"},
}

func TestColors(t *testing.T) {
	// Redirect stdout
	oldStdout := os.Stdout
	_, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("ColorString failed to redirect stdout because %s", err.Error())
	}
	os.Stdout = w

	// Restore old stdout
	defer func() {
		w.Close()
		os.Stdout = oldStdout
	}()

	// Check output
	var out string
	for _, test := range tests {
		if len(test.input) == 1 {
			// With formatter
			out = test.color(Regular, test.input[0])
		} else {
			// Without formatter
			s := make([]interface{}, len(test.input)-1)
			for i, v := range test.input[1:] {
				s[i] = v
			}
			out = test.color(Regular, test.input[0], s...)
		}

		if test.expected != out {
			t.Errorf("Color failed: got %s, want %s", out, test.expected)
		}
	}
}
