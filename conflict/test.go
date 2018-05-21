package conflict

type test struct {
	path            string
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
		markers:         []int{3, 6, 9},
		resolved:        loremResolved,
		resolveDecision: []int{Local},
		numConflicts:    1,
		numLines:        11,
		highlightable:   false,
		shouldPass:      true,
	},
	{
		path:            "testdata/invalid_lorem_ipsum",
		markers:         nil,
		resolved:        nil,
		resolveDecision: nil,
		numConflicts:    0,
		numLines:        11,
		highlightable:   false,
		shouldPass:      false,
	},
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
