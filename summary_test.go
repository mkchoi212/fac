package main

import (
	"strings"
	"testing"

	"github.com/mkchoi212/fac/conflict"
)

func TestFinalizeChanges(t *testing.T) {
	c := conflict.Conflict{}
	c.Choice = Local
	c.Start = 4
	c.End = 10
	c.LocalPureLines = []string{
		"$ go get github.com/mkchoi212/fac\n",
	}

	dummyLines := []string{
		"## 👷 Installation\n",
		"Execute:\n",
		"```bash\n",
		"<<<<<<< Updated upstream:assets/README.md\n",
		"$ go get github.com/mkchoi212/fac\n",
		"||||||| merged common ancestors\n",
		"$ go get github.com/parliament/fac\n",
		"=======\n",
		"$ go get github.com/parliament/facc\n",
		">>>>>>> Stashed changes:README.md\n",
		"```\n",
	}

	output := strings.Join(FinalizeChanges([]conflict.Conflict{c}, dummyLines), "")
	expected := strings.Join([]string{
		"## 👷 Installation\n",
		"Execute:\n",
		"```bash\n",
		"$ go get github.com/mkchoi212/fac\n",
		"```\n",
	}, "")

	if output != expected {
		t.Errorf("FinalizeChanges was incorrect: got \n%s, want \n%s", output, expected)
	}
}
