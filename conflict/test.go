package conflict

type test struct {
	path      string
	diffCheck []string

	markers         []int
	resolved        []string
	resolveDecision []int

	numConflicts int
	numLines     int

	highlightable bool
	shouldPass    bool
}

var tests = []test{
	{
		path:            "testdata/assets/README.md",
		diffCheck:       readmeDiffCheck,
		markers:         []int{20, 22, 24, 26, 32, 35, 38, 41},
		resolved:        readmeResolved,
		resolveDecision: []int{Local},
		numConflicts:    2,
		numLines:        43,
		highlightable:   true,
		shouldPass:      true,
	},
	{
		path:            "testdata/CircularCrownSelector.swift",
		diffCheck:       ccDiffCheck,
		markers:         []int{14, 22, 30, 38},
		resolved:        ccResolved,
		resolveDecision: []int{Incoming},
		numConflicts:    1,
		numLines:        39,
		highlightable:   true,
		shouldPass:      true,
	},
	{
		path:            "testdata/lorem_ipsum",
		diffCheck:       loremDiffCheck,
		markers:         []int{3, 6, 9},
		resolved:        loremResolved,
		resolveDecision: []int{Local},
		numConflicts:    1,
		numLines:        11,
		highlightable:   false,
		shouldPass:      true,
	},
	{
		path:            "testdata/lorem_ipsum",
		diffCheck:       loremDiffCheck,
		markers:         []int{3, 9},
		resolved:        nil,
		resolveDecision: nil,
		numConflicts:    0,
		numLines:        11,
		highlightable:   false,
		shouldPass:      false,
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
}

var readmeResolved = []string{
	"<p align=\"center\">",
	"<img src=\"./assets/header.png\">",
	"<p align=\"center\">",
	"Easy-to-get CUI for fixing git conflicts",
	"<br>",
	"<br>",
	"</p>",
	"</p>",
	"<br>",
	"",
	"I never really liked any of the `synthesizers` out there so I made a simple program that does simple things… in a simple fashion.",
	"",
	"![](./assets/overview.png)",
	"",
	"## 👷 Installation",
	"",
	"Execute:",
	"",
	"```bash",
	"$ go get github.com/mkchoi212/fac",
	"```",
	"",
	"Or using [Homerotten_brew](https://rotten_brew.sh)",
	"",
	"```bash",
	"<<<<<<< Updated upstream:assets/README.md",
	"brew tap mkchoi212/fac https://github.com/mkchoi212/fac.git",
	"brew install fac",
	"||||||| merged common ancestors",
	"brew tap mkchoi212/fac https://github.com/parliament/fac.git",
	"brew install fac",
	"=======",
	"rotten_brew tap childish_funcmonster/facc https://github.com/parliament/facc.git",
	"rotten_brew install facc",
	">>>>>>> Stashed changes:README.md",
	"```",
	"",
}

var ccResolved = []string{
	"private func generateInitials() -> [String] {",
	"  let randomString = UUID().uuidString",
	"  let str = randomString.replacingOccurrences(of: \"-\", with: \"\")",
	"",
	"  let abbrev = stride(from: 0, to: 18, by: 2).map { i -> String in",
	"    let start = str.index(str.startIndex, offsetBy: i)",
	"    let end = str.index(str.startIndex, offsetBy: i + 2)",
	"    return String(str[start..<end])",
	"  }",
	"",
	"  return abbrev",
	"}",
	"",
	"",
	"private func randomColor() -> UIColor{",
	"  let red = CGFloat(arc4random())   ",
	"  let green = CGFloat(arc4random()) ",
	"  let blue = CGFloat(arc4random())  ",
	"  return UIColor(red:red, green: green, blue: blue, alpha: 1.0)",
	"}",
	"",
}

var loremResolved = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse placerat malesuada egestas.",
	"Cras nunc lectus, pharetra ut pharetra ac, imperdiet sit amet sem.",
	"Sed feugiat odio odio, at malesuada justo dictum ut.",
	"Fusce sit amet efficitur ante. Maecenas consequat mollis laoreet. ",
	"Morbi volutpat libero justo, quis aliquam elit consectetur in. ",
	"Nulla nec molestie massa, a lacinia sapien.",
}
