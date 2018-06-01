package color

import (
	"fmt"
)

// Defines the color of the text
const (
	FgBlack = iota
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgPurple
)

// Defines the style of the text
const (
	Regular = iota + 1
	Light
	Highlight
	Underline
)

// Apply applies a `color` and a `style` to the provided `format`
// It does this by using ANSI color codes
func Apply(color int, style int, format string, a ...interface{}) string {
	var str string
	if len(a) == 0 {
		str = format
	} else {
		str = fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", color, style, str)
}

// Red returns a string with red foreground
func Red(style int, format string, a ...interface{}) string {
	return Apply(FgRed, style, format, a...)
}

// Yellow returns a string with red foreground
func Yellow(style int, format string, a ...interface{}) string {
	return Apply(FgYellow, style, format, a...)
}

// Green returns a string with green foreground
func Green(style int, format string, a ...interface{}) string {
	return Apply(FgGreen, style, format, a...)
}

// Blue returns a string with blue foreground
func Blue(style int, format string, a ...interface{}) string {
	return Apply(FgBlue, style, format, a...)
}

// Black returns a string with black foreground
func Black(style int, format string, a ...interface{}) string {
	return Apply(FgBlack, style, format, a...)
}
