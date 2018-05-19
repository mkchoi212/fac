package conflict

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

var mockValidDirectory = "/"
var mockInvalidDirectory = "/hello/world"

var commands = []struct {
	command string
	ok      bool
}{
	{"time", true},
	{"ls", true},
	{"foobar", false},
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Getenv(argsEnvFlag)

	// Purposely fail test
	if args == "false" {
		fmt.Fprintf(os.Stderr, "Mock exec: Command tragically failed")
		os.Exit(1)
	}

	// TopLevelPath arguements
	if args == "rev-parse,--show-toplevel" {
		fmt.Fprintf(os.Stdout, "testdata")
		os.Exit(0)
	}

	// MarkerLocation arguements
	if args == "--no-pager,diff,--check" {
		allCheckOutput := append([]string{}, loremDiffCheck...)
		allCheckOutput = append(allCheckOutput, ccDiffCheck...)
		allCheckOutput = append(allCheckOutput, readmeDiffCheck...)
		fmt.Fprintf(os.Stdout, strings.Join(allCheckOutput, "\n"))
		os.Exit(0)
	}

	fmt.Fprintf(os.Stdout, "Mock exec: Command succeeded")
	os.Exit(0)
}

// Allows us to mock exec.Command, thanks to
// https://npf.io/2015/06/testing-exec-command/
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(cmd.Env, "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestRun(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	for _, test := range commands {
		stdout, stderr, exitCode := run(test.command, ".", strconv.FormatBool(test.ok))

		if test.ok && exitCode != 0 {
			t.Errorf("run failed: got %s with exit code %d, expected no errors", stderr, exitCode)
		} else if !(test.ok) && exitCode == 0 {
			t.Errorf("run failed: got %s with exit code %d, expected errors", stdout, exitCode)
		}
	}
}

func TestTopLevelPath(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []string{mockValidDirectory, mockInvalidDirectory}
	for _, test := range tests {
		out, err := TopLevelPath(test)

		if test == mockValidDirectory && err != nil {
			t.Errorf("TopLevelPath failed: got %s, expected no errors", err.Error())
		} else if test == mockInvalidDirectory && err == nil {
			t.Errorf("TopLevelPath failed: got %s, expected errors", out)
		}
	}
}

func TestMarkerLocations(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []string{mockValidDirectory, mockInvalidDirectory}
	for _, test := range tests {
		out, err := MarkerLocations(test)

		if test == mockValidDirectory && err != nil {
			t.Errorf("TopLevelPath failed: got %s, expected no errors", err.Error())
		} else if test == mockInvalidDirectory && err == nil {
			t.Errorf("TopLevelPath failed: got %s, expected errors", out)
		}
	}
}
