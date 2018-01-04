package color

import (
	"fmt"
)

const (
	FgBlack = iota
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgPurple
)

const (
	Regular = iota + 1
	Light
	Highlight
	Underline
)

func colorString(color int, style int, format string, a ...interface{}) string {
	var str string
	if len(a) == 0 {
		str = format
	} else {
		str = fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", color, style, str)
}

// Black is a convenient helper function to return a string with black foreground
func Black(style int, format string, a ...interface{}) string {
	return colorString(FgBlack, style, format, a...)
}

// Red is a convenient helper function to return a string with red foreground
func Red(style int, format string, a ...interface{}) string {
	return colorString(FgRed, style, format, a...)
}

// Green is a convenient helper function to return a string with green foreground
func Green(style int, format string, a ...interface{}) string {
	return colorString(FgGreen, style, format, a...)
}

// Yellow is a convenient helper function to return a string with yellow foreground
func Yellow(style int, format string, a ...interface{}) string {
	return colorString(FgYellow, style, format, a...)
}

// Blue is a convenient helper function to return a string with blue foreground
func Blue(style int, format string, a ...interface{}) string {
	return colorString(FgBlue, style, format, a...)
}

// Purple is a convenient helper function to return a string with purple foreground
func Purple(style int, format string, a ...interface{}) string {
	return colorString(FgPurple, style, format, a...)
}
