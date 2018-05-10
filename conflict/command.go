package conflict

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

// RunCommand runs the given command with arguments and returns the output
// Refer to https://stackoverflow.com/questions/10385551/get-exit-code-go
func run(name string, dir string, args ...string) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// Success and exitCode should be 0
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return
}

// MarkerLocations returns line numbers of all git conflict markers
func MarkerLocations(cwd string) ([]string, error) {
	stdout, stderr, _ := run("git", cwd, "--no-pager", "diff", "--check")

	if len(stderr) != 0 {
		return nil, errors.New(stderr)
	} else if len(stdout) == 0 {
		return nil, NewErrNoConflict("No conflicts detected 🎉")
	}

	return strings.Split(stdout, "\n"), nil
}

// TopLevelPath finds the top level path of the current git repository
func TopLevelPath(cwd string) (string, error) {
	stdout, stderr, _ := run("git", cwd, "rev-parse", "--show-toplevel")

	if len(stderr) != 0 {
		return "", errors.New(stderr)
	} else if len(stdout) == 0 {
		return "", errors.New(stderr)
	}

	return string(strings.Split(stdout, "\n")[0]), nil
}

// DiffLines is incomplete (TODO)
func DiffLines(cwd string) ([]string, error) {
	stdout, _, _ := run("git", cwd, "--no-pager", "diff", "--color")
	return []string{stdout}, nil
}
