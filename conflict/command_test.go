package conflict

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/mkchoi212/fac/testhelper"
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

	if args == "--no-pager,diff,--name-only,--diff-filter=U" {
		fmt.Fprintf(os.Stdout, "lorem_ipsum\nassets/README.md\n")
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

		if test.ok {
			testhelper.Assert(t, exitCode == 0, "expected no errors but got %s", stderr)
		} else {
			testhelper.Assert(t, exitCode != 0, "expected errors but got %s", stdout)
		}
	}
}

func TestConflictedFiles(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []string{mockValidDirectory, mockInvalidDirectory}
	for _, test := range tests {
		out, err := conflictedFiles(test)

		if test == mockValidDirectory {
			testhelper.Ok(t, err)
		} else if test == mockInvalidDirectory {
			testhelper.Assert(t, err != nil, "expected errors but got %s", out)
		}
	}
}

func TestTopLevelPath(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []string{mockValidDirectory, mockInvalidDirectory}
	for _, test := range tests {
		out, err := topLevelPath(test)

		if test == mockValidDirectory {
			testhelper.Ok(t, err)
		} else if test == mockInvalidDirectory {
			testhelper.Assert(t, err != nil, "expected errors but got %s", out)
		}
	}
}
