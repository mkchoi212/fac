package conflict

import "os"

func (c *Conflict) Diff() []string {
	cwd, _ := os.Getwd()
	stdout, _, _ := RunCommand("git", cwd, "--no-pager", "diff", "--color")
	return []string{stdout}
}
