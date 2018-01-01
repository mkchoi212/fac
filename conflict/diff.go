package conflict

func (c *Conflict) Diff() []string {
	dummyPath := "/Users/mikechoi/src/CSCE-313/"
	stdout, _, _ := RunCommand("git", dummyPath, "--no-pager", "diff", "--color")
	return []string{stdout}
}
