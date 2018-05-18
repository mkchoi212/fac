package conflict

type test struct {
	path          string
	diffCheck     []string
	markers       []int
	highlightable bool
	parsable      bool
}

var tests = []test{
	{
		path:          "testdata/assets/README.md",
		diffCheck:     readmeDiffCheck,
		markers:       []int{20, 22, 24, 26, 32, 35, 38, 41},
		highlightable: true,
		parsable:      true,
	},
	{
		path:          "testdata/CircularCrownSelector.swift",
		diffCheck:     ccDiffCheck,
		markers:       []int{14, 22, 30, 38},
		highlightable: true,
		parsable:      true,
	},
	{
		path:          "testdata/lorem_ipsum",
		diffCheck:     loremDiffCheck,
		markers:       []int{3, 6, 9},
		highlightable: false,
		parsable:      true,
	},
	{
		path:          "testdata/lorem_ipsum",
		diffCheck:     loremDiffCheck,
		markers:       []int{3, 9},
		highlightable: false,
		parsable:      false,
	},
}

var loremDiffCheck = []string{
	"lorem_ipsum:3: leftover conflict marker",
	"lorem_ipsum:6: leftover conflict marker",
	"lorem_ipsum:9: leftover conflict marker",
}

var ccDiffCheck = []string{
	"CircularCrownSelector.swift:14: leftover conflict marker",
	"CircularCrownSelector.swift:22: leftover conflict marker",
	"CircularCrownSelector.swift:30: leftover conflict marker",
	"CircularCrownSelector.swift:31: trailing whitespace.",
	"+    let red = CGFloat(arc4random())",
	"CircularCrownSelector.swift:32: trailing whitespace.",
	"+    let green = CGFloat(arc4random())",
	"CircularCrownSelector.swift:34: trailing whitespace.",
	"+    let blue = CGFloat(arc4random())",
	"CircularCrownSelector.swift:36: trailing whitespace.",
	"+  }",
	"CircularCrownSelector.swift:38: leftover conflict marker",
}

var readmeDiffCheck = []string{
	"assets/README.md:20: leftover conflict marker",
	"assets/README.md:20: trailing whitespace.",
	"+<<<<<<< Updated upstream:assets/README.md",
	"assets/README.md:22: leftover conflict marker",
	"assets/README.md:22: trailing whitespace.",
	"+||||||| merged common ancestors",
	"assets/README.md:23: trailing whitespace.",
	"+$ go get github.com/parliament/fac",
	"assets/README.md:24: leftover conflict marker",
	"assets/README.md:24: trailing whitespace.",
	"+=======",
	"assets/README.md:25: trailing whitespace.",
	"+$ go get github.com/parliament/facc",
	"assets/README.md:26: leftover conflict marker",
	"assets/README.md:26: trailing whitespace.",
	"+>>>>>>> Stashed changes:README.md",
	"assets/README.md:32: leftover conflict marker",
	"assets/README.md:32: trailing whitespace.",
	"+<<<<<<< Updated upstream:assets/README.md",
	"assets/README.md:35: leftover conflict marker",
	"assets/README.md:38: leftover conflict marker",
	"assets/README.md:38: trailing whitespace.",
	"+=======",
	"assets/README.md:40: trailing whitespace.",
	"+rotten_brew install facc",
	"assets/README.md:41: leftover conflict marker",
	"assets/README.md:41: trailing whitespace.",
	"+>>>>>>> Stashed changes:README.md",
	"assets/README.md:46: trailing whitespace.",
	"+> **Please note facc does NOT support diff3 merge conflict outputs yet!**",
	"}",
}
