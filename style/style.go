package style

import "github.com/fatih/color"

var (
	Green      = color.New(color.FgGreen, color.Bold).SprintfFunc()
	GreenLight = color.New(color.FgGreen).SprintfFunc()

	Red      = color.New(color.FgRed, color.Bold).SprintfFunc()
	RedLight = color.New(color.FgRed).SprintfFunc()

	Blue = color.New(color.FgBlue, color.Bold).SprintfFunc()
	Grey = color.New(color.FgBlack, color.Faint).SprintfFunc()
)
