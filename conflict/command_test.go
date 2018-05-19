package conflict

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var mockValidDirectory = "/"
var mockInvalidDirectory = "/hello/world"

var commands = []struct {
	command string
	path    string
	ok      bool
}{
	{"time", mockValidDirectory, true},
	{"ls", mockValidDirectory, true},
	{"less", mockInvalidDirectory, false},
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	cwd := os.Getenv(cwdEnvFlag)

	if cwd == mockInvalidDirectory {
		fmt.Fprintf(os.Stderr, "Mock exec: Command tragically failed")
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Mock exec: Command succeeded")
		os.Exit(0)
	}
}

// Allows us to mock exec.Command, thanks to
// https://npf.io/2015/06/testing-exec-command/
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestRun(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	for _, test := range commands {
		stdout, stderr, exitCode := run(test.command, test.path)

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
