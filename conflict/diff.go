package conflict

// Diff is incomplete (TODO)
func (c *Conflict) Diff() []string {
	lines, _ := DiffLines("")
	return lines
}
