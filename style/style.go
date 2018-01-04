package style

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Green      = color.New(color.FgGreen, color.Bold).SprintfFunc()
	GreenLight = color.New(color.FgGreen).SprintfFunc()

	Red      = color.New(color.FgRed, color.Bold).SprintfFunc()
	RedLight = color.New(color.FgRed).SprintfFunc()

	Blue = color.New(color.FgBlue, color.Bold).SprintfFunc()
	Grey = color.New(color.FgHiBlack).SprintfFunc()
)

// CUIGrey is grey color format function used within the CUI
func CUIGrey(line string) string {
	return fmt.Sprintf("\033[30;1m%s\033[0m", line)
}
