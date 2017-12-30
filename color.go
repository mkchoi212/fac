package main

import "fmt"

const (
	Gray   = 0
	Red    = 1
	Green  = 2
	Yellow = 3
	Blue   = 4
	Purple = 5
)

const (
	Regular   = 1
	Light     = 2
	Highlight = 3
	Underline = 4
)

func Colorize(s string, c int) string {
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", c, Regular, s)
}

func ColorizeLight(s string, c int) string {
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", c, Light, s)
}
