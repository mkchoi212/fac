package main

import "fmt"

const (
	Red    = 1
	Green  = 2
	Yellow = 3
	Blue   = 4
	Purple = 5
)

func Colorize(s string, c int) string {
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", c, 1, s)
}
